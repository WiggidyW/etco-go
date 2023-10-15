package appraisal

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"time"

	"github.com/WiggidyW/chanresult"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/client/market"
	rdbc "github.com/WiggidyW/etco-go/client/remotedb"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

type MakeBuybackAppraisalParams struct {
	Items       []BasicItem
	SystemId    int32
	CharacterId *int32
	Save        bool
}

type MakeBuybackAppraisalClient struct {
	writeClient  rdbc.SAC_WriteBuybackAppraisalClient
	marketClient market.BuybackPriceClient
}

func NewMakeBuybackAppraisalClient(
	writeClient rdbc.SAC_WriteBuybackAppraisalClient,
	marketClient market.BuybackPriceClient,
) MakeBuybackAppraisalClient {
	return MakeBuybackAppraisalClient{writeClient, marketClient}
}

func (bac MakeBuybackAppraisalClient) Fetch(
	ctx context.Context,
	params MakeBuybackAppraisalParams,
) (*rdb.BuybackAppraisal, error) {
	systemInfoPtr := staticdb.GetBuybackSystemInfo(params.SystemId)
	if systemInfoPtr == nil {
		return NewRejectedBuybackAppraisal(params), nil
	}
	systemInfo := *systemInfoPtr

	// fetch the appraisal items in a goroutine
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := chanresult.
		NewChanResult[rdb.BuybackParentItem](ctx, 0, 0).Split()
	for _, item := range params.Items {
		go bac.fetchOne(ctx, item, systemInfo, chnSend)
	}

	// initialize the appraisal
	appraisal := &rdb.BuybackAppraisal{
		// Code: "",
		Items:    make([]rdb.BuybackParentItem, 0, len(params.Items)),
		Price:    0.0,
		Fee:      0.0,
		FeePerM3: systemInfo.M3Fee,
		Tax:      systemInfo.Tax,
		// Time: time.Time{},
		Version:     build.VERSION_BUYBACK,
		SystemId:    params.SystemId,
		CharacterId: params.CharacterId,
	}

	// create a hash code if we're saving the appraisal
	var hasher hash.Hash64
	if params.Save {
		hasher = fnv.New64a() // TODO: Use a sync pool of hashers
	}

	// collect the results
	for i := 0; i < len(params.Items); i++ {
		item, err := chnRecv.Recv()
		if err != nil {
			return nil, err
		}

		// add the item's price and fee to the appraisal if sum is positive
		itemSumPricePerUnit := item.PricePerUnit - item.FeePerUnit
		if itemSumPricePerUnit > 0.0 {
			f64ItemQuantity := float64(item.Quantity)
			appraisal.Price += itemSumPricePerUnit * f64ItemQuantity
			appraisal.Fee += item.FeePerUnit * f64ItemQuantity
		}

		// hash it if saving
		appraisal.Items = append(appraisal.Items, item)
		if params.Save {
			bac.hashItem(hasher, item)
		}
	}

	// if we aren't saving the appraisal, or everythings rejected, finish
	// (no timestamp, no code)
	if !params.Save || appraisal.Price <= 0.0 {
		return appraisal, nil
	}

	// hash the appraisal and add the code
	bac.hashAppraisal(hasher, appraisal)
	appraisal.Code = bac.getHashCode(hasher)

	// save the appraisal
	var timestamp *time.Time
	var err error
	wbaParams := rdbc.WriteBuybackAppraisalParams{Appraisal: *appraisal}
	if params.CharacterId != nil {
		timestamp, err = bac.writeClient.Fetch(ctx, wbaParams)
	} else {
		// no cache invalidation needed for anonymous appraisals
		timestamp, err = bac.writeClient.InnerClient().Fetch(
			ctx,
			wbaParams,
		)
	}
	if err != nil {
		return nil, err
	}

	// add the timestamp
	appraisal.Time = *timestamp

	return appraisal, nil
}

func (bac MakeBuybackAppraisalClient) getHashCode(hasher hash.Hash64) string {
	// 16 characters: u + 15 hex digits (first digit trimmed)
	return "u" + fmt.Sprintf("%016x", hasher.Sum64())[1:]
}

func (bac MakeBuybackAppraisalClient) hashAppraisal(
	hasher hash.Hash64,
	appraisal *rdb.BuybackAppraisal,
) {
	// put the integer values into a buffer
	buf := make([]byte, 24)
	binary.BigEndian.PutUint64(buf[0:], math.Float64bits(appraisal.Price))
	binary.BigEndian.PutUint32(buf[8:], uint32(appraisal.SystemId))
	if appraisal.CharacterId != nil {
		binary.BigEndian.PutUint32(
			buf[12:],
			uint32(*appraisal.CharacterId),
		)
	} else {
		binary.BigEndian.PutUint32(buf[12:], 0)
	}
	binary.BigEndian.PutUint64(buf[16:], uint64(len(appraisal.Items)))

	hasher.Write([]byte(appraisal.Version)) // hash strings directly
	hasher.Write(buf)
}

func (bac MakeBuybackAppraisalClient) hashItem(
	hasher hash.Hash64,
	item rdb.BuybackParentItem,
) {
	// put the integer values into a buffer
	buf := make([]byte, 36)
	binary.BigEndian.PutUint32(buf[0:], uint32(item.TypeId))
	binary.BigEndian.PutUint64(buf[4:], uint64(item.Quantity))
	binary.BigEndian.PutUint64(buf[12:], math.Float64bits(item.PricePerUnit))
	binary.BigEndian.PutUint64(buf[20:], math.Float64bits(item.FeePerUnit))
	binary.BigEndian.PutUint64(buf[28:], uint64(len(item.Children)))

	hasher.Write([]byte(item.Description)) // hash strings directly
	hasher.Write(buf)
}

func (bac MakeBuybackAppraisalClient) fetchOne(
	ctx context.Context,
	item BasicItem,
	systemInfo staticdb.BuybackSystemInfo,
	chnSend chanresult.ChanSendResult[rdb.BuybackParentItem],
) error {
	if apprItem, err := bac.marketClient.Fetch(
		ctx,
		market.BuybackPriceParams{
			BuybackSystemInfo: systemInfo,
			TypeId:            item.TypeId,
			Quantity:          item.Quantity,
		},
	); err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(*apprItem)
	}
}
