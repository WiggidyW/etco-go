package appraisal

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash"
	"hash/fnv"
	"math"

	"github.com/WiggidyW/chanresult"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/client/market"
	rdb "github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

type MakeShopAppraisalParams struct {
	Items       []BasicItem
	LocationId  int64
	CharacterId int32 // optional, just puts it into the response
	IncludeCode bool
}

type MakeShopAppraisalClient struct {
	marketClient market.ShopPriceClient
}

func NewMakeShopAppraisalClient(
	marketClient market.ShopPriceClient,
) MakeShopAppraisalClient {
	return MakeShopAppraisalClient{marketClient}
}

func (sac MakeShopAppraisalClient) Fetch(
	ctx context.Context,
	params MakeShopAppraisalParams,
) (*rdb.ShopAppraisal, error) {
	locationInfoPtr := staticdb.GetShopLocationInfo(params.LocationId)
	if locationInfoPtr == nil {
		return NewRejectedShopAppraisal(params), nil
	}
	locationInfo := *locationInfoPtr

	// fetch the appraisal items in a goroutine
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := chanresult.
		NewChanResult[rdb.ShopItem](ctx, 0, 0).Split()
	for _, item := range params.Items {
		go sac.fetchOne(ctx, item, locationInfo, chnSend)
	}

	// initialize the appraisal
	appraisal := &rdb.ShopAppraisal{
		// Code: "",
		Items: make([]rdb.ShopItem, 0, len(params.Items)),
		Price: 0.0,
		// Time: time.Time{},
		Version:     build.VERSION_SHOP,
		LocationId:  params.LocationId,
		CharacterId: params.CharacterId,
	}

	// create a hash code if it's requested
	var hasher hash.Hash64
	if params.IncludeCode {
		hasher = fnv.New64a() // TODO: Use a sync pool of hashers
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
				sac.hashItem(hasher, item)
			}
		}
	}

	// if we aren't including the code, or everythings rejected, finish
	// (no timestamp, no code)
	if !params.IncludeCode || appraisal.Price <= 0.0 {
		return appraisal, nil
	}

	// hash the appraisal and add the code
	sac.hashAppraisal(hasher, appraisal)
	appraisal.Code = sac.getHashCode(hasher)

	return appraisal, nil
}

func (sac MakeShopAppraisalClient) getHashCode(hasher hash.Hash64) string {
	// 16 characters: s + 15 hex digits (first digit trimmed)
	return "s" + fmt.Sprintf("%016x", hasher.Sum64())[1:]
}

func (sac MakeShopAppraisalClient) hashAppraisal(
	hasher hash.Hash64,
	appraisal *rdb.ShopAppraisal,
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

func (sac MakeShopAppraisalClient) hashItem(
	hasher hash.Hash64,
	item rdb.ShopItem,
) {
	// put the integer values into a buffer
	buf := make([]byte, 20)
	binary.BigEndian.PutUint32(buf[0:], uint32(item.TypeId))
	binary.BigEndian.PutUint64(buf[4:], uint64(item.Quantity))
	binary.BigEndian.PutUint64(buf[12:], math.Float64bits(item.PricePerUnit))

	hasher.Write(buf)
	hasher.Write([]byte(item.Description)) // hash strings directly
}

func (sac MakeShopAppraisalClient) fetchOne(
	ctx context.Context,
	item BasicItem,
	locationInfo staticdb.ShopLocationInfo,
	chnSend chanresult.ChanSendResult[rdb.ShopItem],
) error {
	if apprItem, err := sac.marketClient.Fetch(
		ctx,
		market.ShopPriceParams{
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
