package market

import (
	"fmt"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

const KIND_BUYBACK = "buyback"

func GetBuybackPrice(
	x cache.Context,
	typeId int32,
	quantity int64,
	systemInfo staticdb.BuybackSystemInfo,
) (
	price BuybackPriceParent,
	expires time.Time,
	err error,
) {
	bPricingInfo := systemInfo.GetTypePricingInfo(typeId)
	if bPricingInfo == nil {
		price = newRejectedParent(typeId, quantity)
		expires = fetch.MAX_EXPIRES
	} else if bPricingInfo.ReprocessingEfficiency != 0.0 {
		price, expires, err = bpgReprocessed(
			x,
			typeId,
			quantity,
			systemInfo,
			*bPricingInfo,
		)
	} else if bPricingInfo.PricingInfo != nil {
		price, expires, err = bpgLeaf(
			x,
			typeId,
			quantity,
			systemInfo,
			*bPricingInfo.PricingInfo,
		)
	} else {
		logger.Fatal(fmt.Sprintf(
			"%d: buyback pricing info has neither reprocessed nor leaf pricing",
			typeId,
		))
	}
	return price, expires, err
}

func ProtoGetBuybackPrice(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	typeId int32,
	quantity int64,
	systemInfo staticdb.BuybackSystemInfo,
) (
	price *proto.BuybackParentItem,
	expires time.Time,
	err error,
) {
	var rPrice BuybackPriceParent
	rPrice, expires, err = GetBuybackPrice(x, typeId, quantity, systemInfo)
	if err != nil {
		return nil, expires, err
	}
	return rPrice.ToProto(r), expires, nil
}

func bpgLeaf(
	x cache.Context,
	typeId int32,
	quantity int64,
	systemInfo staticdb.BuybackSystemInfo,
	pricingInfo staticdb.PricingInfo,
) (
	price BuybackPriceParentLeaf,
	expires time.Time,
	err error,
) {
	var positivePrice float64
	positivePrice, expires, err = GetPercentilePrice(x, typeId, pricingInfo)
	if err != nil {
		return price, expires, err
	}
	price = leafUnpackPositivePrice(
		typeId,
		quantity,
		positivePrice,
		pricingInfo,
		systemInfo,
	)
	return price, expires, nil
}

func bpgReprocessed(
	x cache.Context,
	typeId int32,
	quantity int64,
	systemInfo staticdb.BuybackSystemInfo,
	bPricingInfo staticdb.BuybackPricingInfo,
) (
	price BuybackPriceParentRepr,
	expires time.Time,
	err error,
) {
	sdeTypeInfo := sdeTypeInfoOrFatal(KIND_BUYBACK, typeId)
	x, cancel := x.WithCancel()
	defer cancel()
	chnChild := expirable.NewChanResult[remotedb.BuybackChildItem](
		x.Ctx(),
		len(sdeTypeInfo.ReprocessedMaterials),
		0,
	)

	// fetch the buyback price child for each material
	for _, reprMat := range sdeTypeInfo.ReprocessedMaterials {
		go transceiveBPGChild(
			x,
			systemInfo,
			bPricingInfo,
			reprMat.TypeId,
			reprMat.Quantity,
			chnChild,
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
	expires = fetch.MAX_EXPIRES
	var child remotedb.BuybackChildItem
	for i := 0; i < len(sdeTypeInfo.ReprocessedMaterials); i++ {
		child, expires, err = chnChild.RecvExpMin(expires)
		if err != nil {
			return price, expires, err
		} else {
			sumPrice += child.PricePerUnit * child.QuantityPerParent
			children = append(children, child)
		}
	}

	// unpack into a repr variant
	// - accepted
	// - rejected (all children rejected (standard, or no orders))
	// - rejected fee
	price = reprUnpackSumPrice(
		typeId,
		quantity,
		sumPrice,
		children,
		bPricingInfo.ReprocessingEfficiency,
		sdeTypeInfo,
		systemInfo,
	)
	return price, expires, nil
}

func transceiveBPGChild(
	x cache.Context,
	systemInfo staticdb.BuybackSystemInfo,
	parentBPricingInfo staticdb.BuybackPricingInfo,
	typeId int32,
	quantity float64,
	chn expirable.ChanResult[remotedb.BuybackChildItem],
) error {
	return chn.SendExp(bpgChild(
		x,
		systemInfo,
		parentBPricingInfo,
		typeId,
		quantity,
	))
}

func bpgChild(
	x cache.Context,
	systemInfo staticdb.BuybackSystemInfo,
	parentBPricingInfo staticdb.BuybackPricingInfo,
	typeId int32,
	quantity float64,
) (
	price BuybackPriceChild,
	expires time.Time,
	err error,
) {
	quantity = quantity * parentBPricingInfo.ReprocessingEfficiency

	childPricingInfoPtr := bpgChildPricingInfo(
		systemInfo,
		parentBPricingInfo,
		typeId,
	)
	if childPricingInfoPtr == nil {
		price = newRejectedChild(typeId, quantity)
		expires = fetch.MAX_EXPIRES
		return price, expires, nil
	}
	childPricingInfo := *childPricingInfoPtr

	var positivePrice float64
	positivePrice, expires, err = GetPercentilePrice(
		x,
		typeId,
		childPricingInfo,
	)
	if err != nil {
		return price, expires, err
	}

	price = childUnpackPositivePrice(
		typeId,
		quantity,
		positivePrice,
		childPricingInfo,
	)
	return price, expires, nil
}

func bpgChildPricingInfo(
	systemInfo staticdb.BuybackSystemInfo,
	parentBPricingInfo staticdb.BuybackPricingInfo,
	typeId int32,
) *staticdb.PricingInfo {
	// try to get pricing info from parent type info (inherited)
	if parentBPricingInfo.PricingInfo != nil {
		return parentBPricingInfo.PricingInfo
	}

	// try to get it from the system info (unique to the child type)
	childBPricingInfo := systemInfo.GetTypePricingInfo(typeId)
	if childBPricingInfo != nil && childBPricingInfo.PricingInfo != nil {
		return childBPricingInfo.PricingInfo
	}

	// currently, we do not support children with reprocessed-only pricing

	return nil
}

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
	feePerM3 float64,
	typeId int32, // only used if sdeTypeInfo is nil
) (
	accepted bool,
	priceWithFee float64,
	fee float64,
) {
	fee = calculateTypeFee(typeId, sdeTypeInfo, feePerM3)
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
	feePerM3 float64,
) float64 {
	if feePerM3 <= 0.0 {
		return 0.0
	}
	if sdeTypeInfo == nil {
		sdeTypeInfo = sdeTypeInfoOrFatal(KIND_BUYBACK, typeId)
	}
	if sdeTypeInfo.Volume <= 0.0 {
		return 0.0
	}

	fee := sdeTypeInfo.Volume * feePerM3
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

// Leaf = parent with no children
// Reprocessed = parent with children
// Parent = Leaf || Reprocessed

// // Parent (leaf or reprocessed)

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

// // Leaf

type BuybackPriceParentLeaf = BuybackPriceParent

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
	if !accepted {
		return newRejectedLeafNoOrders(
			typeId,
			quantity,
			priceInfo.MarketName,
		)
	}
	accepted, price, fee := priceWithFee(
		positivePrice,
		nil, // won't need if m3 fee is <= 0 / nil
		systemInfo.M3Fee,
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
}

// // Reprocessed

type BuybackPriceParentRepr = BuybackPriceParent

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
			systemInfo.M3Fee,
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
