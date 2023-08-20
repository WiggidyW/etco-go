package buyback

import (
	"github.com/WiggidyW/weve-esi/client/market"
	"github.com/WiggidyW/weve-esi/client/market/internal"
	"github.com/WiggidyW/weve-esi/staticdb"
)

// TODO: handle the sdeTypeInfo / typeId parameters better without duplicating code
func priceWithFee(
	price float64,
	sdeTypeInfo *staticdb.SDETypeInfo,
	systemInfo staticdb.BuybackSystemInfo,
	typeId int32, // only used if sdeTypeInfo is nil
) (accepted bool, _ float64, fee float64) {
	feePtr := calculateTypeFee(sdeTypeInfo, systemInfo, typeId)
	if feePtr == nil {
		return true, price, 0.0
	}
	fee = *feePtr
	if fee <= 0.0 {
		return true, price, 0.0
	}
	price -= fee
	if price < 0.0 {
		return false, 0.0, fee
	}

	return true, price, fee
}

func calculateTypeFee(
	sdeTypeInfo *staticdb.SDETypeInfo,
	systemInfo staticdb.BuybackSystemInfo,
	typeId int32,
) *float64 {
	if systemInfo.M3Fee == nil || *systemInfo.M3Fee <= 0.0 {
		return nil
	}
	if sdeTypeInfo == nil {
		sdeTypeInfo = sdeTypeInfoOrFatal(typeId)
	}
	if sdeTypeInfo.Volume == nil || *sdeTypeInfo.Volume <= 0.0 {
		return nil
	}

	fee := *sdeTypeInfo.Volume * *systemInfo.M3Fee
	return &fee
}

type BuybackPriceChild struct {
	market.MarketPrice
	Quantity float64 // number per 1 parent item
}

func newRejectedChild(quantity float64) *BuybackPriceChild {
	return &BuybackPriceChild{
		MarketPrice: market.MarketPrice{
			Price: 0.0,
			Desc:  market.Rejected(),
		},
		Quantity: quantity,
	}
}

func newRejectedChildNoOrders(mrktName string) *BuybackPriceChild {
	return &BuybackPriceChild{
		MarketPrice: market.MarketPrice{
			Price: 0.0,
			Desc:  market.RejectedNoOrders(mrktName),
		},
		Quantity: 0.0,
	}
}

func newAcceptedChild(
	price float64,
	quantity float64,
	priceInfo staticdb.PricingInfo,
) *BuybackPriceChild {
	return &BuybackPriceChild{
		MarketPrice: market.MarketPrice{
			Price: market.RoundedPrice(price),
			Desc: market.Accepted(
				priceInfo.MrktName,
				priceInfo.Prctile,
				priceInfo.Modifier,
				priceInfo.IsBuy,
			),
		},
		Quantity: quantity,
	}
}

func childUnpackPositivePrice(
	positivePrice *internal.PositivePrice,
	priceInfo staticdb.PricingInfo,
	quantity float64,
) *BuybackPriceChild {
	accepted, price := market.UnpackPositivePrice(positivePrice)
	if accepted {
		return newAcceptedChild(price, quantity, priceInfo)
	} else {
		return newRejectedChildNoOrders(priceInfo.MrktName)
	}
}

// leaf = parent with no children
// reprocessed = parent with children
// parent = leaf || reprocessed
//
// Parent with no children: This could be referred to as a "leaf" node (chatGPT)
type BuybackPriceParent struct {
	market.MarketPrice
	Fee      float64
	Children []BuybackPriceChild
}

// // parent (leaf or reprocessed)
func newRejectedParent() *BuybackPriceParent {
	return &BuybackPriceParent{
		MarketPrice: market.MarketPrice{
			Price: 0.0,
			Desc:  market.Rejected(),
		},
		Fee:      0.0,
		Children: []BuybackPriceChild{},
	}
}

// //

// // leaf
func newRejectedLeafNoOrders(mrktName string) *BuybackPriceParent {
	return &BuybackPriceParent{
		MarketPrice: market.MarketPrice{
			Price: 0.0,
			Desc:  market.RejectedNoOrders(mrktName),
		},
		Fee:      0.0,
		Children: []BuybackPriceChild{},
	}
}

func newRejectedLeafFee(fee float64) *BuybackPriceParent {
	return &BuybackPriceParent{
		MarketPrice: market.MarketPrice{
			Price: 0.0,
			Desc:  market.RejectedFee(),
		},
		Fee:      fee,
		Children: []BuybackPriceChild{},
	}
}

func newAcceptedLeaf(
	price float64,
	fee float64,
	priceInfo staticdb.PricingInfo,
) *BuybackPriceParent {
	return &BuybackPriceParent{
		MarketPrice: market.MarketPrice{
			Price: market.RoundedPrice(price),
			Desc: market.Accepted(
				priceInfo.MrktName,
				priceInfo.Prctile,
				priceInfo.Modifier,
				priceInfo.IsBuy,
			),
		},
		Fee:      fee,
		Children: []BuybackPriceChild{},
	}
}

func leafUnpackPositivePrice(
	positivePrice *internal.PositivePrice,
	priceInfo staticdb.PricingInfo,
	systemInfo staticdb.BuybackSystemInfo,
	typeId int32,
) *BuybackPriceParent {
	accepted, price := market.UnpackPositivePrice(positivePrice)
	if accepted {
		accepted, price, fee := priceWithFee(
			price,
			nil, // won't need if m3 fee is <= 0 / nil
			systemInfo,
			typeId,
		)
		if accepted {
			return newAcceptedLeaf(price, fee, priceInfo)
		} else {
			return newRejectedLeafFee(fee)
		}
	} else {
		return newRejectedLeafNoOrders(priceInfo.MrktName)
	}
}

// //

// // reprocessed
func newRejectedRepr(
	repEff float64,
	children []BuybackPriceChild,
) *BuybackPriceParent {
	return &BuybackPriceParent{
		MarketPrice: market.MarketPrice{
			Price: 0.0,
			Desc:  market.RejectedReprocessed(repEff),
		},
		Fee:      0.0,
		Children: children,
	}
}

func newRejectedReprFee(
	fee float64,
	repEff float64,
	children []BuybackPriceChild,
) *BuybackPriceParent {
	return &BuybackPriceParent{
		MarketPrice: market.MarketPrice{
			Price: 0.0,
			Desc:  market.RejectedReprocessedFee(repEff),
		},
		Fee:      fee,
		Children: children,
	}
}

func newAcceptedRepr(
	price float64,
	fee float64,
	repEff float64,
	children []BuybackPriceChild,
) *BuybackPriceParent {
	return &BuybackPriceParent{
		MarketPrice: market.MarketPrice{
			Price: market.RoundedPrice(price),
			Desc:  market.AcceptedReprocessed(repEff),
		},
		Fee:      fee,
		Children: children,
	}
}

func reprUnpackSumPrice(
	sumPrice float64,
	children []BuybackPriceChild,
	repEff float64,
	sdeTypeInfo *staticdb.SDETypeInfo,
	systemInfo staticdb.BuybackSystemInfo,
	typeId int32,
) *BuybackPriceParent {
	accepted := sumPrice > 0.0
	if accepted {
		accepted, price, fee := priceWithFee(
			sumPrice,
			sdeTypeInfo, // already retrieved previously
			systemInfo,
			typeId,
		)
		if accepted {
			return newAcceptedRepr(price, fee, repEff, children)
		} else {
			return newRejectedReprFee(fee, repEff, children)
		}
	} else {
		return newRejectedRepr(repEff, children)
	}
}

// //
