package buyback

import (
	"context"
	"fmt"

	"github.com/WiggidyW/weve-esi/client/market/internal"
	"github.com/WiggidyW/weve-esi/logger"
	"github.com/WiggidyW/weve-esi/staticdb"
	"github.com/WiggidyW/weve-esi/util"
)

type BuybackPriceClient struct {
	Inner    internal.MarketPriceClient
	MinPrice float64
}

// gets the BuybackPrice for the given item
func (bc BuybackPriceClient) Fetch(
	ctx context.Context,
	params BuybackPriceParams,
) (*BuybackPriceParent, error) {
	bTypeInfo := params.BuybackSystemInfo.GetTypeInfo(params.TypeId)
	if bTypeInfo == nil { // item not sold at locations shop
		return newRejectedParent(), nil
	}

	// // determine pricing based on which are present, pricing vs repr eff
	if bTypeInfo.ReprEff != nil {
		// reprocessed pricing
		return bc.fetchReprocessed(
			ctx,
			params,
			*bTypeInfo,
		)
	} else {
		// leaf pricing
		return bc.fetchLeaf(ctx, params, *bTypeInfo.PricingInfo)
	}
}

func (bc *BuybackPriceClient) fetchReprocessed(
	ctx context.Context,
	params BuybackPriceParams,
	bTypeInfo staticdb.BuybackTypeInfo,
) (*BuybackPriceParent, error) {
	sdeTypeInfo := sdeTypeInfoOrFatal(params.TypeId)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := util.NewChanResult[BuybackPriceChild](ctx).Split()

	// fetch the buyback price child for each material
	for _, reprMat := range sdeTypeInfo.ReprMats {
		go bc.fetchChild(
			ctx,
			params,
			bTypeInfo,
			reprMat,
			chnSend,
		)
	}

	// initialize the collections
	sumPrice := 0.0
	children := make([]BuybackPriceChild, 0, len(sdeTypeInfo.ReprMats))

	// collect the children
	for i := 0; i < len(sdeTypeInfo.ReprMats); i++ {
		if child, err := chnRecv.Recv(); err != nil {
			return nil, err
		} else {
			sumPrice += child.MarketPrice.Price * child.Quantity
			children = append(children, child)
		}
	}

	// unpack into a repr variant
	// - accepted
	// - rejected (all children rejected (standard, or no orders))
	// - rejected fee
	return reprUnpackSumPrice(
		sumPrice,
		children,
		*bTypeInfo.ReprEff,
		sdeTypeInfo,
		params.BuybackSystemInfo,
		params.TypeId,
	), nil
}

func (bpc BuybackPriceClient) fetchChild(
	ctx context.Context,
	params BuybackPriceParams,
	parentBTypeInfo staticdb.BuybackTypeInfo,
	reprMat staticdb.ReprocessedMaterial,
	chnSend util.ChanSendResult[BuybackPriceChild],
) error {
	childQuantity := reprMat.Quantity * *parentBTypeInfo.ReprEff

	// in order to fetch the market price, we need to know the pricing info
	priceInfoPtr := getChildPriceInfo(
		parentBTypeInfo,
		params.BuybackSystemInfo,
		reprMat.TypeId,
	)
	if priceInfoPtr == nil {
		// return standard rejected
		return chnSend.SendOk(*newRejectedChild(childQuantity))
	}
	priceInfo := *priceInfoPtr

	price, err := bpc.fetchMarketPrice(ctx, priceInfo, reprMat.TypeId)
	if err != nil {
		return chnSend.SendErr(err)
	} else {
		// unpack into a child variant
		// - accepted
		// - rejected no orders
		return chnSend.SendOk(*childUnpackPositivePrice(
			price,
			priceInfo,
			childQuantity,
		))
	}
}

func (bpc BuybackPriceClient) fetchLeaf(
	ctx context.Context,
	params BuybackPriceParams,
	priceInfo staticdb.PricingInfo,
) (*BuybackPriceParent, error) {
	price, err := bpc.fetchMarketPrice(ctx, priceInfo, params.TypeId)
	if err != nil {
		return nil, err
	} else {
		// unpack into a leaf variant
		// - accepted
		// - rejected no orders
		// - rejected fee
		return leafUnpackPositivePrice(
			price,
			priceInfo,
			params.BuybackSystemInfo,
			params.TypeId,
		), nil
	}
}

func (bpc BuybackPriceClient) fetchMarketPrice(
	ctx context.Context,
	priceInfo staticdb.PricingInfo,
	typeId int32,
) (*internal.PositivePrice, error) {
	return bpc.Inner.Fetch(
		ctx,
		internal.MarketPriceParams{
			PricingInfo: priceInfo,
			TypeId:      typeId,
		},
	)
}

func getChildPriceInfo(
	parentBTypeInfo staticdb.BuybackTypeInfo,
	systemInfo staticdb.BuybackSystemInfo,
	childTypeId int32,
) *staticdb.PricingInfo {
	// try to get pricing info from parent type info (inherited)
	if parentBTypeInfo.PricingInfo != nil {
		return parentBTypeInfo.PricingInfo
	}

	// try to get it from the system info (unique to the child type)
	childBTypeInfo := systemInfo.GetTypeInfo(childTypeId)
	if childBTypeInfo != nil && childBTypeInfo.PricingInfo != nil {
		return childBTypeInfo.PricingInfo
	}

	// currently, we do not support children with reprocessing efficiency

	return nil
}

func sdeTypeInfoOrFatal(typeId int32) *staticdb.SDETypeInfo {
	t := staticdb.GetSDETypeInfo(typeId)
	if t == nil {
		logger.Logger.Fatal(fmt.Sprintf(
			"buyback valid type %d not found in sde type info",
			typeId,
		))
	}
	return t
}
