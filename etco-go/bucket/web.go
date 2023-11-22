package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_CAPACITY_MULTIPLIER = 3
	WEB_CAPACITY_DIVISOR    = 2
)

func webGet[K comparable, V any](
	x cache.Context,
	method func(context.Context, int) (map[K]V, error),
	cacheKey, typeStr keys.Key,
	expiresIn time.Duration,
	makeCap int,
) (
	rep map[K]V,
	expires time.Time,
	err error,
) {
	return get(
		x,
		func(ctx context.Context) (map[K]V, error) {
			return method(ctx, transformWebCapacity(makeCap))
		},
		cacheKey, typeStr,
		expiresIn,
		makeMapPtrFunc[K, V](makeCap),
	)
}

func makeMapPtrFunc[K comparable, V any](
	capacity int,
) func() map[K]V {
	return func() map[K]V {
		m := make(map[K]V, capacity)
		return m
	}
}

func transformWebCapacity(capacity int) int {
	return capacity * WEB_CAPACITY_MULTIPLIER / WEB_CAPACITY_DIVISOR
}

// // Proto

func newPBCfgTypePricing(
	webTypePricing b.WebTypePricing,
) (
	pbTypePricing *proto.CfgTypePricing,
) {
	return &proto.CfgTypePricing{
		IsBuy:      webTypePricing.IsBuy,
		Percentile: uint32(webTypePricing.Percentile),
		Modifier:   uint32(webTypePricing.Modifier),
		Market:     webTypePricing.MarketName,
	}
}

func PBtoWebTypePricing[HSV any](
	pbPricing *proto.CfgTypePricing,
	markets map[string]HSV,
) (
	webTypePricing *b.WebTypePricing,
	err error,
) {
	if pbPricing.Percentile > 100 {
		return nil, newPBtoWebTypePricingError(
			"percentile must be <= 100",
		)
	} else if pbPricing.Modifier > 255 {
		return nil, newPBtoWebTypePricingError(
			"modifier must be <= 255",
		)
	} else if _, ok := markets[pbPricing.Market]; !ok {
		return nil, newPBtoWebTypePricingError(
			fmt.Sprintf(
				"market '%s' does not exist",
				pbPricing.Market,
			),
		)
	}

	return &b.WebTypePricing{
		IsBuy:      pbPricing.IsBuy,
		Percentile: uint8(pbPricing.Percentile),
		Modifier:   uint8(pbPricing.Modifier),
		MarketName: pbPricing.Market,
	}, nil
}

func newPBtoWebTypePricingError(
	errStr string,
) configerror.ErrPricingInvalid {
	return configerror.ErrPricingInvalid{ErrString: errStr}
}
