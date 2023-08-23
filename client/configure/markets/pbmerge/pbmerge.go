package pbmerge

import (
	cfg "github.com/WiggidyW/weve-esi/client/configure"
	"github.com/WiggidyW/weve-esi/proto"
	"github.com/WiggidyW/weve-esi/staticdb"
	"github.com/WiggidyW/weve-esi/util"
)

func ConvertPBMarket(
	name string,
	pbMarket *proto.Market,
) (cfg.Market, error) {
	if pbMarket.IsStructure && pbMarket.RefreshToken == "" {
		return cfg.Market{}, newError(
			name,
			"structure market must have refresh token",
		)
	} else if !pbMarket.IsStructure {
		if pbMarket.RefreshToken != "" {
			return cfg.Market{}, newError(
				name,
				"non-structure market must not have refresh token",
			)
		}
		_, ok := staticdb.GetStation(int32(pbMarket.LocationId))
		if !ok {
			return cfg.Market{}, newError(
				name,
				"station does not exist",
			)
		}
	}
	var refreshToken *string
	if pbMarket.RefreshToken != "" {
		refreshToken = &pbMarket.RefreshToken
	}
	return cfg.Market{
		RefreshToken: refreshToken,
		LocationId:   pbMarket.LocationId,
		IsStructure:  pbMarket.IsStructure,
	}, nil
}

func MergeMarkets[HS util.HashSet[string]](
	original cfg.Markets,
	updates map[string]*proto.Market,
	activeMapNames HS,
) error {
	// if updates == nil || len(updates.Inner) == 0 {
	// 	return false, nil
	// }
	for name, pbMarket := range updates {
		if pbMarket == nil {
			if activeMapNames.Has(name) {
				return newError(
					name,
					"market is currently in use",
				)
			} else {
				delete(original, name)
			}
		} else {
			market, err := ConvertPBMarket(name, pbMarket)
			if err != nil {
				return err
			}
			original[name] = market
		}
	}
	return nil
}

// since getting ActiveMapNames is expensive operation, check if it's needed
func ActiveMapNamesNeeded(updates map[string]*proto.Market) bool {
	for _, pbMarket := range updates {
		if pbMarket == nil {
			return true
		}
	}
	return false
}

func extractBuybackBuilderActiveMarkets(
	bbuilder cfg.BuybackSystemTypeMapsBuilder,
) map[string]struct{} {
	activeMarkets := make(map[string]struct{})
	for _, bundle := range bbuilder {
		for _, bTypePricing := range bundle {
			if bTypePricing.Pricing == nil {
				continue
			}
			activeMarkets[bTypePricing.Pricing.Market] = struct{}{}
		}
	}
	return activeMarkets
}

func extractShopBuilderActiveMarkets(
	sbuilder cfg.ShopLocationTypeMapsBuilder,
) map[string]struct{} {
	activeMarkets := make(map[string]struct{})
	for _, bundle := range sbuilder {
		for _, sTypePricing := range bundle {
			activeMarkets[sTypePricing.Market] = struct{}{}
		}
	}
	return activeMarkets
}
