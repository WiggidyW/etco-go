package webtocore

import (
	"fmt"

	b "github.com/WiggidyW/etco-go-bucket"
)

func convertWebShopLocations(
	webShopLocations map[b.LocationId]b.WebShopLocation,
	coreSLTypeMapsIndexMap map[b.BundleKey]int,
) (
	coreShopLocations map[b.LocationId]b.ShopLocation,
	coreBannedFlagSets []b.BannedFlagSet,
	err error,
) {
	coreShopLocations = make(
		map[b.LocationId]b.ShopLocation,
		len(webShopLocations),
	)
	coreBannedFlagSets = make([]b.BannedFlagSet, 0)

	for locationId, webShopLocation := range webShopLocations {
		typeMapIndex, ok := coreSLTypeMapsIndexMap[webShopLocation.
			BundleKey]
		if !ok {
			return nil, nil, fmt.Errorf(
				"ShopLocation %d has invalid BundleKey %s",
				locationId,
				webShopLocation.BundleKey,
			)
		}
		bannedFlagSet := newBannedFlagSet(webShopLocation.BannedFlags)
		bannedFlagSetIndex := getBannedFlagSetIndex(
			bannedFlagSet,
			&coreBannedFlagSets,
		)
		coreShopLocations[locationId] = b.ShopLocation{
			BannedFlagSetIndex: bannedFlagSetIndex,
			TypeMapIndex:       typeMapIndex,
		}
	}

	return coreShopLocations, coreBannedFlagSets, nil
}

// possibly mutates coreBannedFlagSets
func getBannedFlagSetIndex(
	bannedFlagSet b.BannedFlagSet,
	bannedFlagSets *[]b.BannedFlagSet,
) (index int) {
	for existingIndex, existingBannedFlagSet := range *bannedFlagSets {
		if bannedFlagSetsEqual(
			bannedFlagSet,
			existingBannedFlagSet,
		) {
			return existingIndex
		}
	}
	index = len(*bannedFlagSets)
	*bannedFlagSets = append(*bannedFlagSets, bannedFlagSet)
	return index
}

func newBannedFlagSet(bannedFlags []string) b.BannedFlagSet {
	bannedFlagSet := make(b.BannedFlagSet, len(bannedFlags))
	for _, bannedFlag := range bannedFlags {
		bannedFlagSet[bannedFlag] = struct{}{}
	}
	return bannedFlagSet
}

func bannedFlagSetsEqual(a, b b.BannedFlagSet) bool {
	if len(a) != len(b) {
		return false
	}
	for bannedFlag := range a {
		if _, ok := b[bannedFlag]; !ok {
			return false
		}
	}
	return true
}
