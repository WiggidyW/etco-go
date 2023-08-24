package buyback

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"time"

	"github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	f "github.com/WiggidyW/eve-trading-co-go/client/authingfwding/fwding"
	bm "github.com/WiggidyW/eve-trading-co-go/client/market/buyback"
	wb "github.com/WiggidyW/eve-trading-co-go/client/remotedb/appraisal/writebuyback"
	"github.com/WiggidyW/eve-trading-co-go/staticdb"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type F_BuybackAppraisalClient = f.FwdingClient[
	FWD_BuybackAppraisalParams,
	BuybackAppraisalParams,
	appraisal.BuybackAppraisal,
	BuybackAppraisalClient,
]

type BuybackAppraisalClient struct {
	writeCharClient wb.SAC_WriteBuybackAppraisalClient
	writeAnonClient wb.WriteBuybackAppraisalClient // no cache invalidation needed for anonymous appraisals
	marketClient    bm.BuybackPriceClient
}

func (bac BuybackAppraisalClient) Fetch(
	ctx context.Context,
	params BuybackAppraisalParams,
) (*appraisal.BuybackAppraisal, error) {
	systemInfoPtr := staticdb.GetBuybackSystemInfo(params.SystemId)
	if systemInfoPtr == nil {
		return NewRejectedAppraisal(params), nil
	}
	systemInfo := *systemInfoPtr

	// fetch the appraisal items in a goroutine
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := util.NewChanResult[appraisal.BuybackParentItem](
		ctx,
	).Split()
	for _, item := range params.Items {
		go bac.fetchOne(ctx, item, systemInfo, chnSend)
	}

	// initialize the appraisal
	appraisal := &appraisal.BuybackAppraisal{
		// Code: "",
		Items: make([]appraisal.BuybackParentItem, 0, len(params.Items)),
		Price: 0.0,
		// Time: time.Time{},
		Version:     staticdb.BUYBACK_VERSION,
		SystemId:    params.SystemId,
		CharacterId: params.CharacterId,
	}

	// create a hash code if we're saving the appraisal
	var hasher hash.Hash64 = nil
	if params.Save {
		hasher = fnv.New64a()
	}

	// collect the results
	for i := 0; i < len(params.Items); i++ {
		if item, err := chnRecv.Recv(); err != nil {
			return nil, err
		} else {
			appraisal.Price += item.PricePerUnit * float64(
				item.Quantity,
			)
			appraisal.Items = append(appraisal.Items, item)
			if params.Save {
				hashItem(hasher, item)
			}
		}
	}

	// if we aren't saving the appraisal, or everythings rejected, finish
	// (no timestamp, no code)
	if !params.Save || appraisal.Price <= 0.0 {
		return appraisal, nil
	}

	// hash the appraisal and add the code
	hashAppraisal(hasher, appraisal)
	appraisal.Code = getCode(hasher)

	// save the appraisal
	var timestamp *time.Time
	var err error
	wbaParams := wb.WriteBuybackAppraisalParams{Appraisal: *appraisal}
	if params.CharacterId != nil {
		timestamp, err = bac.writeCharClient.Fetch(ctx, wbaParams)
	} else {
		timestamp, err = bac.writeAnonClient.Fetch(ctx, wbaParams)
	}
	if err != nil {
		return nil, err
	}

	// add the timestamp
	appraisal.Time = *timestamp

	return appraisal, nil
}

func getCode(hasher hash.Hash64) string {
	// 16 characters: b + 15 hex digits (first digit trimmed)
	return "b" + fmt.Sprintf("%016x", hasher.Sum64())[1:]
}

func hashAppraisal(
	hasher hash.Hash64,
	appraisal *appraisal.BuybackAppraisal,
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

func hashItem(
	hasher hash.Hash64,
	item appraisal.BuybackParentItem,
) {
	// put the integer values into a buffer
	buf := make([]byte, 36)
	binary.BigEndian.PutUint32(buf[0:], uint32(item.TypeId))
	binary.BigEndian.PutUint64(buf[4:], uint64(item.Quantity))
	binary.BigEndian.PutUint64(buf[12:], math.Float64bits(item.PricePerUnit))
	binary.BigEndian.PutUint64(buf[20:], math.Float64bits(item.Fee))
	binary.BigEndian.PutUint64(buf[28:], uint64(len(item.Children)))

	hasher.Write([]byte(item.Description)) // hash strings directly
	hasher.Write(buf)
}

func (bac BuybackAppraisalClient) fetchOne(
	ctx context.Context,
	item appraisal.BasicItem,
	systemInfo staticdb.BuybackSystemInfo,
	chnSend util.ChanSendResult[appraisal.BuybackParentItem],
) error {
	if apprItem, err := bac.marketClient.Fetch(
		ctx,
		bm.BuybackPriceParams{
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
