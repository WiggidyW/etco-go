package market

import (
	"context"
	"fmt"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/client/market/marketprice"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/staticdb"
)

type BuybackPriceClient struct {
	marketPriceClient marketprice.MarketPriceClient
}

func NewBuybackPriceClient(
	marketPriceClient marketprice.MarketPriceClient,
) BuybackPriceClient {
	return BuybackPriceClient{marketPriceClient}
}

// gets the BuybackPrice for the given item
func (bc BuybackPriceClient) Fetch(
	ctx context.Context,
	params BuybackPriceParams,
) (*BuybackPriceParent, error) {
	bTypeInfo := params.BuybackSystemInfo.GetTypePricingInfo(params.TypeId)
	if bTypeInfo == nil { // item not sold at locations shop
		return newRejectedParent(params.TypeId, params.Quantity), nil
	}

	// // determine pricing based on which are present, pricing vs repr eff
	if bTypeInfo.ReprocessingEfficiency != 0.0 {
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
	bPricingInfo staticdb.BuybackPricingInfo,
) (*BuybackPriceParentRepr, error) {
	sdeTypeInfo := sdeTypeInfoOrFatal(params.TypeId)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chnSend, chnRecv := chanresult.
		NewChanResult[BuybackPriceChild](ctx, 0, 0).Split()

	// fetch the buyback price child for each material
	for _, reprMat := range sdeTypeInfo.ReprocessedMaterials {
		go bc.fetchChild(
			ctx,
			params,
			bPricingInfo,
			reprMat,
			chnSend,
		)
	}

	// initialize the collections
	sumPrice := 0.0
	children := make(
		[]BuybackPriceChild,
		0,
		len(sdeTypeInfo.ReprocessedMaterials),
	)

	// collect the children
	for i := 0; i < len(sdeTypeInfo.ReprocessedMaterials); i++ {
		if child, err := chnRecv.Recv(); err != nil {
			return nil, err
		} else {
			sumPrice += child.PricePerUnit * child.QuantityPerParent
			children = append(children, child)
		}
	}

	// unpack into a repr variant
	// - accepted
	// - rejected (all children rejected (standard, or no orders))
	// - rejected fee
	return reprUnpackSumPrice(
		params.TypeId,
		params.Quantity,
		sumPrice,
		children,
		bPricingInfo.ReprocessingEfficiency,
		sdeTypeInfo,
		params.BuybackSystemInfo,
	), nil
}

func (bpc BuybackPriceClient) fetchChild(
	ctx context.Context,
	params BuybackPriceParams,
	parentBPricingInfo staticdb.BuybackPricingInfo,
	reprMat b.ReprocessedMaterial,
	chnSend chanresult.ChanSendResult[BuybackPriceChild],
) error {
	childQuantity := reprMat.Quantity *
		parentBPricingInfo.ReprocessingEfficiency

	// in order to fetch the market price, we need to know the pricing info
	priceInfoPtr := getChildPriceInfo(
		parentBPricingInfo,
		params.BuybackSystemInfo,
		reprMat.TypeId,
	)
	if priceInfoPtr == nil {
		// return standard rejected
		return chnSend.SendOk(*newRejectedChild(
			reprMat.TypeId,
			childQuantity,
		))
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
			reprMat.TypeId,
			childQuantity,
			price,
			priceInfo,
		))
	}
}

func (bpc BuybackPriceClient) fetchLeaf(
	ctx context.Context,
	params BuybackPriceParams,
	priceInfo staticdb.PricingInfo,
) (*BuybackPriceParentLeaf, error) {
	price, err := bpc.fetchMarketPrice(ctx, priceInfo, params.TypeId)
	if err != nil {
		return nil, err
	} else {
		// unpack into a leaf variant
		// - accepted
		// - rejected no orders
		// - rejected fee
		return leafUnpackPositivePrice(
			params.TypeId,
			params.Quantity,
			price,
			priceInfo,
			params.BuybackSystemInfo,
		), nil
	}
}

func (bpc BuybackPriceClient) fetchMarketPrice(
	ctx context.Context,
	priceInfo staticdb.PricingInfo,
	typeId int32,
) (*marketprice.PositivePrice, error) {
	return bpc.marketPriceClient.Fetch(
		ctx,
		marketprice.MarketPriceParams{
			PricingInfo: priceInfo,
			TypeId:      typeId,
		},
	)
}

func getChildPriceInfo(
	parentBPricingInfo staticdb.BuybackPricingInfo,
	systemInfo staticdb.BuybackSystemInfo,
	childTypeId int32,
) *staticdb.PricingInfo {
	// try to get pricing info from parent type info (inherited)
	if parentBPricingInfo.PricingInfo != nil {
		return parentBPricingInfo.PricingInfo
	}

	// try to get it from the system info (unique to the child type)
	childBPricingInfo := systemInfo.GetTypePricingInfo(childTypeId)
	if childBPricingInfo != nil && childBPricingInfo.PricingInfo != nil {
		return childBPricingInfo.PricingInfo
	}

	// currently, we do not support children with reprocessed-only pricing

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
