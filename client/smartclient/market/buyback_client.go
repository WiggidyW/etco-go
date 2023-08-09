package market

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/cache"
	"github.com/WiggidyW/weve-esi/client"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/desc"
	"github.com/WiggidyW/weve-esi/client/smartclient/market/percentile"
	"github.com/WiggidyW/weve-esi/logger"
	"github.com/WiggidyW/weve-esi/staticdb/sde"
	"github.com/WiggidyW/weve-esi/staticdb/tc"
)

const MIN_BUYBACK_PRICE = 0.01

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
		percentile.MarketPercentileClientFetchParams,
		percentile.MarketPercentile,
		cache.ExpirableData[percentile.MarketPercentile],
		*percentile.MarketPercentileClient,
	]
}

// gets the BuybackPrice for the given item
func (bc *BuybackClient) Fetch(
	ctx context.Context,
	params BuybackClientFetchParams,
) (BuybackPrice, error) {
	// static data
	// get the buyback info
	buybackInfo := tc.KVReaderBuybackInfo.Get(1) // capacity of 1
	// get the system
	system, ok := buybackInfo.GetLocation(params.SystemId)
	if !ok { // system does not buyback anything
		return newBuybackPrice(
			0,
			desc.Rejected(),
			[]BuybackPriceChild{},
		), nil
	}
	// get the type info
	typeInfo, ok := system.GetType(params.TypeId)
	if !ok { // type not bought back at system
		return newBuybackPrice(
			0,
			desc.Rejected(),
			[]BuybackPriceChild{},
		), nil
	}

	// check if the type has pricing and / or repeff
	pricing, hasPricing := typeInfo.Pricing()
	_, hasRepEff := typeInfo.RepEffRaw()

	// determine pricing based on which are present
	if hasRepEff { // reprocessed pricing
		if hasPricing {
			// reprocessed pricing with material-specific pricing
			return bc.fetchReprocessed(
				ctx,
				params,
				pricing,
				system,
				typeInfo,
			)
		} else {
			// reprocessed pricing with override pricing
			return bc.fetchReprocessed(
				ctx,
				params,
				nil,
				system,
				typeInfo,
			)
		}
	} else if hasPricing { // childless pricing
		if price, descr, err := bc.fetchInner(
			ctx,
			params,
			pricing,
			system,
		); err != nil {
			return newBuybackPrice(
				0,
				"",
				[]BuybackPriceChild{},
			), err
		} else {
			return newBuybackPrice(
				price,
				descr,
				[]BuybackPriceChild{},
			), nil
		}
	} else { // no pricing + no repeff = fatal error
		logger.Logger.Fatal(fmt.Sprintf(
			"buyback pricing and repeff for %d at %d are null",
			params.TypeId,
			params.SystemId,
		))
		// logger.Logger.Fatal is synchronous, so this is unreachable
		panic("unreachable")
	}
}

func (bc *BuybackClient) fetchReprocessed(
	ctx context.Context,
	params BuybackClientFetchParams,
	pricing *tc.PricingInfo,
	system *tc.BuybackSystemInfo,
	typeInfo *tc.BuybackTypeInfo,
) (BuybackPrice, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sdeTypeInfo := sdeTypeInfoOrFatal(params.TypeId)
	materials := sdeTypeInfo.ReprocessedMaterials()
	repEff, ok := typeInfo.RepEff()
	if !ok {
		logger.Logger.Fatal("fetchReprocessed: repeff is nil")
	}
	chnOk := make(chan BuybackPriceChild, len(materials))
	chnErr := make(chan error, len(materials))

	// fetch the buyback price child for each material
	for _, m := range sdeTypeInfo.ReprocessedMaterials() {
		go bc.fetchMaterial(
			ctx,
			params,
			system,
			pricing,
			repEff,
			m,
			chnOk,
			chnErr,
		)
	}

	// initialize the parent
	parent := &BuybackPrice{
		price:    0,
		desc:     "",
		Children: make([]BuybackPriceChild, 0, len(materials)),
	}

	// collect the goroutine results
	for i := 0; i < len(materials); i++ {
		select {
		case child := <-chnOk:
			parent.price += child.price * child.Quantity
			parent.Children = append(parent.Children, child)
		case err := <-chnErr:
			return *parent, err
		}
	}

	// set the parent description
	repEffRaw, _ := typeInfo.RepEffRaw()
	if parent.price > 0 {
		parent.desc = desc.AcceptedReprocessed(repEffRaw)
	} else {
		parent.desc = desc.RejectedReprocessed(repEffRaw)
	}

	// return the parent
	return *parent, nil
}

func (bc *BuybackClient) fetchMaterial(
	ctx context.Context,
	params BuybackClientFetchParams,
	system *tc.BuybackSystemInfo,
	mPricing *tc.PricingInfo,
	repEff float64,
	m sde.ReprocessedMaterial,
	chnOk chan<- BuybackPriceChild,
	chnErr chan<- error,
) {
	// try to find the pricing for this material if it's nil
	if mPricing == nil {
		if tcTypeInfo, ok := system.GetType(m.TypeId); ok {
			if mPricing, ok = tcTypeInfo.Pricing(); !ok {
				mPricing = nil
			}
		}
	}
	// fetch the price for this material
	if mPricing == nil {
		// no pricing for this material, reject it sychronously
		chnOk <- newBuybackPriceChild(
			0,
			desc.Rejected(),
			m.Quantity*repEff,
		)
	} else {
		// fetch the price for this material asynchronously
		go func() {
			if price, descr, err := bc.fetchInner(
				ctx,
				NewBuybackClientFetchParams(
					m.TypeId,
					params.SystemId,
				),
				mPricing,
				system,
			); err != nil {
				chnErr <- err
			} else {
				chnOk <- newBuybackPriceChild(
					price,
					descr,
					m.Quantity*repEff,
				)
			}
		}()
	}
}

func (bc *BuybackClient) fetchInner(
	ctx context.Context,
	params BuybackClientFetchParams,
	pricing *tc.PricingInfo,
	system *tc.BuybackSystemInfo,
) (price float64, descr string, err error) {
	// validate the modifier
	modifier := pricing.Modifier()
	if modifier == 0 {
		logger.Logger.Fatal(fmt.Sprintf(
			"buyback pricing modifier for %d at %d is 0",
			params.TypeId,
			params.SystemId,
		))
	}
	// fetch the percentile price
	percentile, err := bc.client.Fetch(
		ctx,
		percentile.NewFetchParams(pricing, params.TypeId),
	)
	if err != nil {
		return 0, "", err
	} else if percentile.Data().Rejected != "" {
		return 0, percentile.Data().Rejected, nil
	}
	// get the fee
	var fee float64
	if volume, ok := sdeTypeInfoOrFatal(params.TypeId).Volume(); ok {
		fee = system.M3Fee() * volume
	} else {
		fee = 0
	}
	// return the price and desc
	price = minPriced( // minned(rounded(multed - fee))
		roundedToCents(
			multedByModifier(
				percentile.Data().Price,
				modifier,
			)-fee,
		),
		MIN_BUYBACK_PRICE,
	)
	if fee > 0 { // add fee to desc
		descr = desc.AcceptedWithFee(
			pricing.MarketName(),
			pricing.Percentile(),
			pricing.Modifier(),
			pricing.IsBuy(),
			fee,
		)
	} else { // desc with no fee
		descr = desc.Accepted(
			pricing.MarketName(),
			pricing.Percentile(),
			pricing.Modifier(),
			pricing.IsBuy(),
		)
	}
	return
}

func sdeTypeInfoOrFatal(typeId int32) *sde.TypeInfo {
	t, ok := sde.KVReaderTypeInfo.Get(typeId)
	if !ok {
		logger.Logger.Fatal(fmt.Sprintf(
			"buyback valid type %d not found in sde type info",
			typeId,
		))
	}
	return t
}
