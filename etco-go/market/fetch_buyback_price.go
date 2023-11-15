package market

import (
	"fmt"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/remotedb"
	"github.com/WiggidyW/etco-go/staticdb"
)

func buybackPriceGet(
	x cache.Context,
	typeId int32,
	quantity int64,
	buybackSystemInfo staticdb.BuybackSystemInfo,
) (
	price BuybackPriceParent,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetch(
		x,
		nil,
		buybackPriceGetFetchFunc(typeId, quantity, buybackSystemInfo),
		nil,
	)
}

func buybackPriceGetFetchFunc(
	typeId int32,
	quantity int64,
	systemInfo staticdb.BuybackSystemInfo,
) fetch.Fetch[BuybackPriceParent] {
	return func(x cache.Context) (
		price BuybackPriceParent,
		expires time.Time,
		_ *postfetch.Params,
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
		return price, expires, nil, err
	}
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
