package proto

import (
	"context"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/client/appraisal"
	"github.com/WiggidyW/etco-go/client/shopqueue"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoutil"
	"github.com/WiggidyW/etco-go/staticdb"
)

// TODO: DEDUPLICATE / DRY
// TODO: MAKE THIS LESS COMPLICATED, COMPLEX, AND CONFUSING
// crazy channel logic
// fairly fast but at what cost

type PBShopPurchaseQueueClient struct {
	pbGetShopAppraisalClient PBGetShopAppraisalClient[*staticdb.SyncIndexMap]
	pbNewShopAppraisalClient PBNewShopAppraisalClient[*staticdb.SyncIndexMap]
	rShopQueueClient         shopqueue.ShopQueueClient
}

func NewPBShopPurchaseQueueClient(
	pbGetShopAppraisalClient PBGetShopAppraisalClient[*staticdb.SyncIndexMap],
	pbNewShopAppraisalClient PBNewShopAppraisalClient[*staticdb.SyncIndexMap],
	rShopQueueClient shopqueue.ShopQueueClient,
) PBShopPurchaseQueueClient {
	return PBShopPurchaseQueueClient{
		pbGetShopAppraisalClient,
		pbNewShopAppraisalClient,
		rShopQueueClient,
	}
}

func (gspqc PBShopPurchaseQueueClient) Fetch(
	ctx context.Context,
	params PBPurchaseQueueParams,
) (
	entries []*proto.PurchaseQueueEntry,
	err error,
) {
	rShopQueueRep, err := gspqc.rShopQueueClient.Fetch(
		ctx,
		shopqueue.ShopQueueParams{},
	)
	if err != nil {
		return entries, err
	}
	rShopQueue := rShopQueueRep.ParsedShopQueue

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chnSendShopPurchaseQueueEntry, chnRecvShopPurchaseQueueEntry :=
		chanresult.NewChanResult[*proto.PurchaseQueueEntry](
			ctx,
			len(rShopQueue),
			0,
		).Split()

	// for each contract, start a goroutine to fetch the queue entry
	for _, appraisalCode := range rShopQueue {
		go gspqc.transceiveFetchEntry(
			ctx,
			params,
			appraisalCode,
			chnSendShopPurchaseQueueEntry,
		)
	}

	// collect the queue entries
	entries = make(
		[]*proto.PurchaseQueueEntry,
		0,
		len(rShopQueue),
	)
	for i := 0; i < len(rShopQueue); i++ {
		entry, err := chnRecvShopPurchaseQueueEntry.Recv()
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (gspqc PBShopPurchaseQueueClient) transceiveFetchEntry(
	ctx context.Context,
	params PBPurchaseQueueParams,
	appraisalCode string,
	chnSend chanresult.ChanSendResult[*proto.PurchaseQueueEntry],
) error {
	shopPurchaseQueueEntry, err := gspqc.fetchEntry(
		ctx,
		params,
		appraisalCode,
	)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(shopPurchaseQueueEntry)
	}
}

func (gspqc PBShopPurchaseQueueClient) fetchEntry(
	ctx context.Context,
	params PBPurchaseQueueParams,
	appraisalCode string,
) (entry *proto.PurchaseQueueEntry, err error) {
	entry = &proto.PurchaseQueueEntry{Code: appraisalCode}

	// if params.QueueInclude == PQI_NONE {}

	if params.QueueInclude == PQI_CODE_APPRAISAL {
		appraisalWithCode, err := gspqc.fetchCodeAppraisal(
			ctx,
			params.TypeNamingSession,
			appraisalCode,
		)
		if err != nil {
			return nil, err
		}
		entry.CodeAppraisal = appraisalWithCode.Appraisal
		entry.AppraisalCharacterId = appraisalWithCode.CharacterId

	} else if params.QueueInclude == PQI_CODE_APPRAISAL_AND_NEW_APPRAISAL {
		codeAndNewAppraisals, err := gspqc.fetchCodeAndNewAppraisals(
			ctx,
			params.TypeNamingSession,
			appraisalCode,
		)
		if err != nil {
			return nil, err
		}
		entry.CodeAppraisal =
			codeAndNewAppraisals.codeAppraisal.Appraisal
		entry.AppraisalCharacterId =
			codeAndNewAppraisals.codeAppraisal.CharacterId
		entry.NewAppraisal = codeAndNewAppraisals.newAppraisal
	}

	return entry, nil
}

func (gspqc PBShopPurchaseQueueClient) fetchCodeAndNewAppraisals(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
) (
	appraisals shopCodeAndNewAppraisals,
	err error,
) {
	appraisals.codeAppraisal, err = gspqc.fetchCodeAppraisal(
		ctx,
		namingSesssion,
		appraisalCode,
	)
	if err != nil {
		return appraisals, err
	}

	appraisals.newAppraisal, err = gspqc.fetchNewAppraisal(
		ctx,
		namingSesssion,
		protoutil.NewRBasicItems(
			appraisals.codeAppraisal.Appraisal.Items,
		),
		appraisals.codeAppraisal.Appraisal.LocationId,
	)
	if err != nil {
		return appraisals, err
	}

	return appraisals, nil
}

func (gspqc PBShopPurchaseQueueClient) fetchCodeAppraisal(
	ctx context.Context,
	namingSession *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	appraisalCode string,
) (
	appraisal AppraisalWithCharacter[proto.ShopAppraisal],
	err error,
) {
	return gspqc.pbGetShopAppraisalClient.Fetch(
		ctx,
		PBGetAppraisalParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSession,
			AppraisalCode:     appraisalCode,
		},
	)
}

func (gspqc PBShopPurchaseQueueClient) fetchNewAppraisal(
	ctx context.Context,
	namingSesssion *staticdb.TypeNamingSession[*staticdb.SyncIndexMap],
	rItems []appraisal.BasicItem,
	locationId int64,
) (*proto.ShopAppraisal, error) {
	return gspqc.pbNewShopAppraisalClient.Fetch(
		ctx,
		PBNewShopAppraisalParams[*staticdb.SyncIndexMap]{
			TypeNamingSession: namingSesssion,
			Items:             rItems,
			LocationId:        locationId,
			CharacterId:       0,
			IncludeCode:       false,
		},
	)
}
