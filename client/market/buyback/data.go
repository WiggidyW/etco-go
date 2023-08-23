package buyback

import (
	"github.com/WiggidyW/weve-esi/client/appraisal"
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
	feePtr := calculateTypeFee(typeId, sdeTypeInfo, systemInfo)
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
	typeId int32,
	sdeTypeInfo *staticdb.SDETypeInfo,
	systemInfo staticdb.BuybackSystemInfo,
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

// // Child
type BuybackPriceChild = appraisal.BuybackChildItem

func newRejectedChild(typeId int32, quantity float64) *BuybackPriceChild {
	return &BuybackPriceChild{
		PricePerUnit:      0.0,
		Description:       market.Rejected(),
		TypeId:            typeId,
		QuantityPerParent: quantity,
	}
}

func newRejectedChildNoOrders(
	typeId int32,
	mrktName string,
) *BuybackPriceChild {
	return &BuybackPriceChild{
		PricePerUnit:      0.0,
		Description:       market.RejectedNoOrders(mrktName),
		TypeId:            typeId,
		QuantityPerParent: 0.0,
	}
}

func newAcceptedChild(
	typeId int32,
	quantity float64,
	price float64,
	priceInfo staticdb.PricingInfo,
) *BuybackPriceChild {
	return &BuybackPriceChild{
		PricePerUnit: market.RoundedPrice(price),
		Description: market.Accepted(
			priceInfo.MrktName,
			priceInfo.Prctile,
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
	positivePrice *internal.PositivePrice,
	priceInfo staticdb.PricingInfo,
) *BuybackPriceChild {
	accepted, price := market.UnpackPositivePrice(positivePrice)
	if accepted {
		return newAcceptedChild(typeId, quantity, price, priceInfo)
	} else {
		return newRejectedChildNoOrders(typeId, priceInfo.MrktName)
	}
}

// //

// leaf = parent with no children
// reprocessed = parent with children
// parent = leaf || reprocessed
//
// Parent with no children: This could be referred to as a "leaf" node (chatGPT)

// // parent (leaf or reprocessed)
type BuybackPriceParent = appraisal.BuybackParentItem

func newRejectedParent(typeId int32, quantity int64) *BuybackPriceParent {
	return &BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		Fee:          0.0,
		Description:  market.Rejected(),
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
) *BuybackPriceParentLeaf {
	return &BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		Fee:          0.0,
		Description:  market.RejectedNoOrders(mrktName),
		Children:     []BuybackPriceChild{},
	}
}

func newRejectedLeafFee(
	typeId int32,
	quantity int64,
	fee float64,
) *BuybackPriceParentLeaf {
	return &BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		Fee:          fee,
		Description:  market.RejectedFee(),
		Children:     []BuybackPriceChild{},
	}
}

func newAcceptedLeaf(
	typeId int32,
	quantity int64,
	price float64,
	fee float64,
	priceInfo staticdb.PricingInfo,
) *BuybackPriceParentLeaf {
	return &BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: market.RoundedPrice(price),
		Fee:          fee,
		Description: market.Accepted(
			priceInfo.MrktName,
			priceInfo.Prctile,
			priceInfo.Modifier,
			priceInfo.IsBuy,
		),
		Children: []BuybackPriceChild{},
	}
}

func leafUnpackPositivePrice(
	typeId int32,
	quantity int64,
	positivePrice *internal.PositivePrice,
	priceInfo staticdb.PricingInfo,
	systemInfo staticdb.BuybackSystemInfo,
) *BuybackPriceParentLeaf {
	accepted, price := market.UnpackPositivePrice(positivePrice)
	if accepted {
		accepted, price, fee := priceWithFee(
			price,
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
			priceInfo.MrktName,
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
) *BuybackPriceParentRepr {
	return &BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		Fee:          0.0,
		Description:  market.RejectedReprocessed(repEff),
		Children:     children,
	}
}

func newRejectedReprFee(
	typeId int32,
	quantity int64,
	fee float64,
	repEff float64,
	children []BuybackPriceChild,
) *BuybackPriceParentRepr {
	return &BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: 0.0,
		Fee:          fee,
		Description:  market.RejectedReprocessedFee(repEff),
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
) *BuybackPriceParentRepr {
	return &BuybackPriceParent{
		TypeId:       typeId,
		Quantity:     quantity,
		PricePerUnit: market.RoundedPrice(price),
		Fee:          fee,
		Description:  market.AcceptedReprocessed(repEff),
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
) *BuybackPriceParentRepr {
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
