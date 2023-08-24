package shop

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"

	"github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	af "github.com/WiggidyW/eve-trading-co-go/client/authingfwding/authing"
	sm "github.com/WiggidyW/eve-trading-co-go/client/market/shop"
	"github.com/WiggidyW/eve-trading-co-go/staticdb"
	"github.com/WiggidyW/eve-trading-co-go/util"
)

type AF_ShopAppraisalClient = af.AuthingClient[
	FWD_ShopAppraisalParams,
	ShopAppraisalParams,
	appraisal.ShopAppraisal,
	ShopAppraisalClient,
]

type ShopAppraisalClient struct {
	marketClient sm.ShopPriceClient
}

func (sac ShopAppraisalClient) Fetch(
	ctx context.Context,
	params ShopAppraisalParams,
) (*appraisal.ShopAppraisal, error) {
	locationInfoPtr := staticdb.GetShopLocationInfo(params.LocationId)
	if locationInfoPtr == nil {
		return NewRejectedAppraisal(params), nil
	}
	locationInfo := *locationInfoPtr

	// fetch the appraisal items in a goroutine
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := util.NewChanResult[appraisal.ShopItem](ctx).Split()
	for _, item := range params.Items {
		go sac.fetchOne(ctx, item, locationInfo, chnSend)
	}

	// initialize the appraisal
	appraisal := &appraisal.ShopAppraisal{
		// Code: "",
		Items: make([]appraisal.ShopItem, 0, len(params.Items)),
		Price: 0.0,
		// Time: time.Time{},
		Version:     staticdb.SHOP_VERSION,
		LocationId:  params.LocationId,
		CharacterId: params.CharacterId,
	}

	// create a hash code if it's requested
	var hasher hash.Hash64 = nil
	if params.IncludeCode {
		hasher = fnv.New64a()
	}

	// receive the items and add them to the appraisal
	for i := 0; i < len(params.Items); i++ {
		if item, err := chnRecv.Recv(); err != nil {
			return nil, err
		} else {
			appraisal.Price += item.PricePerUnit * float64(
				item.Quantity,
			)
			appraisal.Items = append(appraisal.Items, item)
			if params.IncludeCode {
				hashItem(hasher, item)
			}
		}
	}

	// if we aren't including the code, or everythings rejected, finish
	// (no timestamp, no code)
	if !params.IncludeCode || appraisal.Price <= 0.0 {
		return appraisal, nil
	}

	// hash the appraisal and add the code
	hashAppraisal(hasher, appraisal)
	appraisal.Code = getCode(hasher)

	return appraisal, nil
}

func getCode(hasher hash.Hash64) string {
	// 16 characters: b + 15 hex digits (first digit trimmed)
	return "s" + fmt.Sprintf("%016x", hasher.Sum64())[1:]
}

func hashAppraisal(
	hasher hash.Hash64,
	appraisal *appraisal.ShopAppraisal,
) {
	// put the integer values into a buffer
	buf := make([]byte, 32)
	binary.BigEndian.PutUint64(buf[0:], math.Float64bits(appraisal.Price))
	binary.BigEndian.PutUint64(buf[8:], uint64(appraisal.LocationId))
	binary.BigEndian.PutUint32(buf[16:], uint32(appraisal.CharacterId))
	binary.BigEndian.PutUint64(buf[24:], uint64(len(appraisal.Items)))

	hasher.Write(buf)
	hasher.Write([]byte(appraisal.Version)) // hash strings directly
}

func hashItem(
	hasher hash.Hash64,
	item appraisal.ShopItem,
) {
	// put the integer values into a buffer
	buf := make([]byte, 20)
	binary.BigEndian.PutUint32(buf[0:], uint32(item.TypeId))
	binary.BigEndian.PutUint64(buf[4:], uint64(item.Quantity))
	binary.BigEndian.PutUint64(buf[12:], math.Float64bits(item.PricePerUnit))

	hasher.Write(buf)
	hasher.Write([]byte(item.Description)) // hash strings directly
}

func (sac ShopAppraisalClient) fetchOne(
	ctx context.Context,
	item appraisal.BasicItem,
	locationInfo staticdb.ShopLocationInfo,
	chnSend util.ChanSendResult[appraisal.ShopItem],
) error {
	if apprItem, err := sac.marketClient.Fetch(
		ctx,
		sm.ShopPriceParams{
			ShopLocationInfo: locationInfo,
			TypeId:           item.TypeId,
			Quantity:         item.Quantity,
		},
	); err != nil {
		return chnSend.SendErr(err)
	} else {
		return chnSend.SendOk(*apprItem)
	}
}
