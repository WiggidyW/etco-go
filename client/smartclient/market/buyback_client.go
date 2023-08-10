package market

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/desc"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/percentile"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/percentile/orders"
	"github.com/WiggidyW/weve-esi/logger"
	"github.com/WiggidyW/weve-esi/staticdb"
	"github.com/WiggidyW/weve-esi/staticdb/inner/sde"
	"github.com/WiggidyW/weve-esi/util"
)

const MIN_BUYBACK_PRICE = 0.01 // minimum price for non-rejected items

type BuybackClientFetchParams struct {
	TypeId   int32
	SystemId int32
}

func NewBuybackClientFetchParams(
	typeId int32,
	systemId int32,
) BuybackClientFetchParams {
	return BuybackClientFetchParams{
		TypeId:   typeId,
		SystemId: systemId,
	}
}

type BuybackClient struct {
	client *client.CachingClient[
		percentile.MrktPrctileParams,
		percentile.MrktPrctile,
		cache.ExpirableData[percentile.MrktPrctile],
		*percentile.MrktPrctileClient,
	]
}

// gets the BuybackPrice for the given item
func (bc *BuybackClient) Fetch(
	ctx context.Context,
	params BuybackClientFetchParams,
) (BuybackPrice, error) {
	// // static data
	// get the system
	bbSystem := staticdb.GetBuybackSystemInfo(params.SystemId)
	if bbSystem == nil { // location has no shop
		return BuybackPrice{
			0,
			desc.Rejected(),
			[]BuybackPriceChild{},
		}, nil
	}
	// get the type info
	bbTypeInfo := bbSystem.GetTypeInfo(params.TypeId)
	if bbTypeInfo == nil { // item not sold at locations shop
		return BuybackPrice{
			0,
			desc.Rejected(),
			[]BuybackPriceChild{},
		}, nil
	}

	// // determine pricing based on which are present, pricing vs repr eff

	if bbTypeInfo.ReprEff != nil { // reprocessed pricing
		return bc.fetchReprocessed(
			ctx,
			params,
			*bbSystem,
			*bbTypeInfo,
		)

	} else { // childless pricing
		if price, descr, err := bc.fetchInner(
			ctx,
			params,
			*bbSystem,
			*bbTypeInfo.PricingInfo,
		); err != nil {
			return BuybackPrice{
				0,
				"",
				[]BuybackPriceChild{},
			}, err
		} else {
			return BuybackPrice{
				price,
				descr,
				[]BuybackPriceChild{},
			}, nil
		}
	}
}

func (bc *BuybackClient) fetchReprocessed(
	ctx context.Context,
	params BuybackClientFetchParams,
	bbSystem staticdb.BuybackSystemInfo,
	bbTypeInfo staticdb.BuybackTypeInfo,
) (BuybackPrice, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	sdeTypeInfo := sdeTypeInfoOrFatal(params.TypeId)
	chnSend, chnRecv := util.NewChanResult[BuybackPriceChild](ctx).Split()

	// fetch the buyback price child for each material
	for _, reprMat := range sdeTypeInfo.ReprMats {
		go bc.fetchMaterial(
			ctx,
			params,
			bbSystem,
			bbTypeInfo,
			reprMat,
			chnSend,
		)
	}

	// initialize the parent
	parent := &BuybackPrice{
		price:    0,
		desc:     "",
		Children: make([]BuybackPriceChild, 0, len(sdeTypeInfo.ReprMats)),
	}

	// collect the children
	for i := 0; i < len(sdeTypeInfo.ReprMats); i++ {
		if child, err := chnRecv.Recv(); err != nil {
			return *parent, err
		} else {
			parent.price += child.price * child.Quantity
			parent.Children = append(parent.Children, child)
		}
	}

	// set the parent description
	if parent.price > 0 {
		parent.desc = desc.AcceptedReprocessed(*bbTypeInfo.ReprEff)
	} else {
		parent.desc = desc.RejectedReprocessed(*bbTypeInfo.ReprEff)
	}

	// return the parent
	return *parent, nil
}

func (bc *BuybackClient) fetchMaterial(
	ctx context.Context,
	params BuybackClientFetchParams,
	system staticdb.BuybackSystemInfo,
	parentTypeInfo staticdb.BuybackTypeInfo,
	reprMat sde.ReprocessedMaterial,
	chn util.ChanSendResult[BuybackPriceChild],
) {
	// // select pricing for this material
	var pricing staticdb.PricingInfo

	// if it's not nil, just dereference it from parent type info
	if parentTypeInfo.PricingInfo != nil {
		pricing = *parentTypeInfo.PricingInfo

		// otherwise, try to get it from static data
	} else {
		reprMatTypeInfo := system.GetTypeInfo(reprMat.TypeId)

		// if it's not nil, use the material's unique pricing
		if reprMatTypeInfo != nil &&
			reprMatTypeInfo.PricingInfo != nil {
			pricing = *reprMatTypeInfo.PricingInfo

			// if there's no pricing to be found, reject it
		} else {
			chn.SendOk(BuybackPriceChild{
				price:    0,
				desc:     desc.Rejected(),
				Quantity: reprMat.Quantity * *parentTypeInfo.ReprEff,
			})
			return
		}
	}

	// using the selected pricing, fetch the price
	if price, descr, err := bc.fetchInner(
		ctx,
		NewBuybackClientFetchParams(
			reprMat.TypeId,
			params.SystemId,
		),
		system,
		pricing,
	); err != nil {
		chn.SendErr(err)
	} else {
		chn.SendOk(BuybackPriceChild{
			price:    price,
			desc:     descr,
			Quantity: reprMat.Quantity * *parentTypeInfo.ReprEff,
		})
	}
}

func (bc *BuybackClient) fetchInner(
	ctx context.Context,
	params BuybackClientFetchParams,
	system staticdb.BuybackSystemInfo,
	pricing staticdb.PricingInfo,
) (price float64, descr string, err error) {
	sdeTypeInfo := sdeTypeInfoOrFatal(params.TypeId)

	// fetch the prctile price
	var prctile percentile.MrktPrctile
	if prctileRep, err := bc.client.Fetch(
		ctx,
		percentile.MrktPrctileParams{
			MrktOrdersParams: orders.MrktOrdersParams{
				PricingInfo: pricing,
				TypeId:      params.TypeId,
			},
		},
	); err != nil {
		return 0, "", err
	} else {
		prctile = prctileRep.Data()
	}

	// return early if it's rejected
	if prctile.Rejected != "" {
		return 0, prctile.Rejected, nil
	}

	// get the fee
	var fee float64
	if sdeTypeInfo.Volume != nil && system.M3Fee != nil {
		fee = *sdeTypeInfo.Volume * *system.M3Fee
	} else {
		fee = 0
	}

	// set the price and desc
	price = minPriced( // minned(rounded(multed - fee))
		roundedToCents(
			multedByModifier(
				prctile.Price,
				pricing.Modifier,
			)-fee,
		),
		MIN_BUYBACK_PRICE,
	)

	// set the desc
	if fee > 0 { // add fee to desc
		descr = desc.AcceptedWithFee(
			pricing.MrktName,
			pricing.Prctile,
			pricing.Modifier,
			pricing.IsBuy,
			fee,
		)
	} else { // desc with no fee
		descr = desc.Accepted(
			pricing.MrktName,
			pricing.Prctile,
			pricing.Modifier,
			pricing.IsBuy,
		)
	}

	return price, descr, nil
}

func sdeTypeInfoOrFatal(typeId int32) staticdb.SDETypeInfo {
	t := staticdb.GetSDETypeInfo(typeId)
	if t == nil {
		logger.Logger.Fatal(fmt.Sprintf(
			"buyback valid type %d not found in sde type info",
			typeId,
		))
	}
	return *t
}
