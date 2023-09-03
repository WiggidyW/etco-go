package proto

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"

	"github.com/WiggidyW/etco-go/error/configerror"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/util"
)

type CfgMergeResponse struct {
	Modified   bool
	MergeError error
}

func extractBuilderBundleKeys[V any](
	builder map[int32]map[string]V,
) map[string]struct{} {
	bundleKeys := make(map[string]struct{})
	for _, bundle := range builder {
		for bundleKey := range bundle {
			bundleKeys[bundleKey] = struct{}{}
		}
	}
	return bundleKeys
}

func extractBuybackBuilderActiveMarkets(
	bbuilder map[b.TypeId]b.WebBuybackSystemTypeBundle,
) map[string]struct{} {
	activeMarkets := make(map[string]struct{})
	for _, bundle := range bbuilder {
		for _, bTypePricing := range bundle {
			if bTypePricing.Pricing == nil {
				continue
			}
			activeMarkets[bTypePricing.Pricing.MarketName] =
				struct{}{}
		}
	}
	return activeMarkets
}

func extractShopBuilderActiveMarkets(
	sbuilder map[b.TypeId]b.WebShopLocationTypeBundle,
) map[string]struct{} {
	activeMarkets := make(map[string]struct{})
	for _, bundle := range sbuilder {
		for _, sTypePricing := range bundle {
			activeMarkets[sTypePricing.MarketName] = struct{}{}
		}
	}
	return activeMarkets
}

func PBtoWebTypePricing[HS util.HashSet[string]](
	pbPricing *proto.CfgTypePricing,
	markets HS,
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
	} else if !markets.Has(pbPricing.Market) {
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
