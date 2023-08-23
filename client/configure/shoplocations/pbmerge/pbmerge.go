package pbmerge

import (
	"fmt"

	cfg "github.com/WiggidyW/weve-esi/client/configure"
	"github.com/WiggidyW/weve-esi/proto"
	"github.com/WiggidyW/weve-esi/util"
)

func ConvertPBShopLocation(
	pbShopLocation *proto.ShopLocation,
) cfg.ShopLocation {
	return cfg.ShopLocation{
		BundleKey:   pbShopLocation.BundleKey,
		BannedFlags: pbShopLocation.BannedFlags,
	}
}

func MergeShopLocations[HS util.HashSet[string]](
	original cfg.ShopLocations,
	updates map[int64]*proto.ShopLocation,
	typeMapKeys HS,
) error {
	for locationId, pbShopLocation := range updates {
		if pbShopLocation == nil {
			delete(original, locationId)
		} else if !typeMapKeys.Has(pbShopLocation.BundleKey) {
			return newError(locationId, fmt.Sprintf(
				"type map key '%s' does not exist",
				pbShopLocation.BundleKey,
			))
		} else {
			original[locationId] = ConvertPBShopLocation(
				pbShopLocation,
			)
		}
	}
	return nil
}
