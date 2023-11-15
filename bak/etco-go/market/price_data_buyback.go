package market

import (
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

const KIND_BUYBACK = "buyback"

func NewRejectedBuybackItem(
	typeId int32,
	quantity int64,
) remotedb.BuybackParentItem {
	return newRejectedParent(typeId, quantity)
}

// TODO: handle the sdeTypeInfo / typeId parameters better without duplicating code
func priceWithFee(
	price float64,
	sdeTypeInfo *staticdb.SDETypeInfo,
	systemInfo staticdb.BuybackSystemInfo,
	typeId int32, // only used if sdeTypeInfo is nil
) (
	accepted bool,
	priceWithFee float64,
	fee float64,
) {
	fee = calculateTypeFee(typeId, sdeTypeInfo, systemInfo)
	if fee <= 0.0 {
		return true, price, 0.0
	}
	if price-fee < 0.0 {
		return false, 0.0, fee
	}
	return true, price, fee
}

func calculateTypeFee(
	typeId int32,
	sdeTypeInfo *staticdb.SDETypeInfo,
	systemInfo staticdb.BuybackSystemInfo,
) float64 {
	if systemInfo.M3Fee <= 0.0 {
		return 0.0
	}
	if sdeTypeInfo == nil {
		sdeTypeInfo = sdeTypeInfoOrFatal(KIND_BUYBACK, typeId)
	}
	if sdeTypeInfo.Volume <= 0.0 {
		return 0.0
	}

	fee := sdeTypeInfo.Volume * systemInfo.M3Fee
	return fee
}

// // Child
type BuybackPriceChild = remotedb.BuybackChildItem

func newRejectedChild(typeId int32, quantity float64) BuybackPriceChild {
	return BuybackPriceChild{
		PricePerUnit:      0.0,
		Description:       Rejected(),
		TypeId:            typeId,
		QuantityPerParent: quantity,
	}
}

func newRejectedChildNoOrders(
	typeId int32,
	mrktName string,
) BuybackPriceChild {
	return BuybackPriceChild{
		PricePerUnit:      0.0,
		Description:       RejectedNoOrders(mrktName),
		TypeId:            typeId,
		QuantityPerParent: 0.0,
	}
}

func newAcceptedChild(
	typeId int32,
	quantity float64,
	price float64,
	priceInfo staticdb.PricingInfo,
) BuybackPriceChild {
	return BuybackPriceChild{
		PricePerUnit: RoundedPrice(price),
		Description: Accepted(
			priceInfo.MarketName,
			priceInfo.Percentile,
			priceInfo.Modifier,
			priceInfo.IsBuy,
		),
		TypeId:            typeId,
		QuantityPerParent: quantity,
	}
}

func childUnpackPositivePrice(
	typeId int32,
	quantity float64,
	positivePrice float64,
	priceInfo staticdb.PricingInfo,
) BuybackPriceChild {
	accepted := positivePrice > 0.0
	if accepted {
		return newAcceptedChild(typeId, quantity, positivePrice, priceInfo)
	} else {
		return newRejectedChildNoOrders(typeId, priceInfo.MarketName)
	}
}

// //

// leaf = parent with no children
// reprocessed = parent with children
// parent = leaf || reprocessed
//
// Parent with no children: This could be referred to as a "leaf" node (chatGPT)

// // parent (leaf or reprocessed)
type BuybackPriceParent = remotedb.BuybackParentItem

func newRejectedParent(typeId int32, quantity int64) BuybackPriceParent {
	return BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		FeePerUnit:   0.0,
		Description:  Rejected(),
		Children:     []BuybackPriceChild{},
	}
}

// //

type BuybackPriceParentLeaf = BuybackPriceParent

// // leaf
func newRejectedLeafNoOrders(
	typeId int32,
	quantity int64,
	mrktName string,
) BuybackPriceParentLeaf {
	return BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		FeePerUnit:   0.0,
		Description:  RejectedNoOrders(mrktName),
		Children:     []BuybackPriceChild{},
	}
}

func newRejectedLeafFee(
	typeId int32,
	quantity int64,
	fee float64,
) BuybackPriceParentLeaf {
	return BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		FeePerUnit:   fee,
		Description:  RejectedFee(),
		Children:     []BuybackPriceChild{},
	}
}

func newAcceptedLeaf(
	typeId int32,
	quantity int64,
	price float64,
	fee float64,
	priceInfo staticdb.PricingInfo,
) BuybackPriceParentLeaf {
	return BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: RoundedPrice(price),
		FeePerUnit:   fee,
		Description: Accepted(
			priceInfo.MarketName,
			priceInfo.Percentile,
			priceInfo.Modifier,
			priceInfo.IsBuy,
		),
		Children: []BuybackPriceChild{},
	}
}

func leafUnpackPositivePrice(
	typeId int32,
	quantity int64,
	positivePrice float64,
	priceInfo staticdb.PricingInfo,
	systemInfo staticdb.BuybackSystemInfo,
) BuybackPriceParentLeaf {
	accepted := positivePrice > 0.0
	if accepted {
		accepted, price, fee := priceWithFee(
			positivePrice,
			nil, // won't need if m3 fee is <= 0 / nil
			systemInfo,
			typeId,
		)
		if accepted {
			return newAcceptedLeaf(
				typeId,
				quantity,
				price,
				fee,
				priceInfo,
			)
		} else {
			return newRejectedLeafFee(
				typeId,
				quantity,
				fee,
			)
		}
	} else {
		return newRejectedLeafNoOrders(
			typeId,
			quantity,
			priceInfo.MarketName,
		)
	}
}

// //

type BuybackPriceParentRepr = BuybackPriceParent

// // reprocessed
func newRejectedRepr(
	typeId int32,
	quantity int64,
	repEff float64,
	children []BuybackPriceChild,
) BuybackPriceParentRepr {
	return BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		FeePerUnit:   0.0,
		Description:  RejectedReprocessed(repEff),
		Children:     children,
	}
}

func newRejectedReprFee(
	typeId int32,
	quantity int64,
	fee float64,
	repEff float64,
	children []BuybackPriceChild,
) BuybackPriceParentRepr {
	return BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		FeePerUnit:   fee,
		Description:  RejectedReprocessedFee(repEff),
		Children:     children,
	}
}

func newAcceptedRepr(
	typeId int32,
	quantity int64,
	price float64,
	fee float64,
	repEff float64,
	children []BuybackPriceChild,
) BuybackPriceParentRepr {
	return BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: RoundedPrice(price),
		FeePerUnit:   fee,
		Description:  AcceptedReprocessed(repEff),
		Children:     children,
	}
}

func reprUnpackSumPrice(
	typeId int32,
	quantity int64,
	sumPrice float64,
	children []BuybackPriceChild,
	repEff float64,
	sdeTypeInfo *staticdb.SDETypeInfo,
	systemInfo staticdb.BuybackSystemInfo,
) BuybackPriceParentRepr {
	accepted := sumPrice > 0.0
	if accepted {
		accepted, price, fee := priceWithFee(
			sumPrice,
			sdeTypeInfo, // already retrieved previously
			systemInfo,
			typeId,
		)
		if accepted {
			return newAcceptedRepr(
				typeId,
				quantity,
				price,
				fee,
				repEff,
				children,
			)
		} else {
			return newRejectedReprFee(
				typeId,
				quantity,
				fee,
				repEff,
				children,
			)
		}
	} else {
		return newRejectedRepr(
			typeId,
			quantity,
			repEff,
			children,
		)
	}
}

// //
