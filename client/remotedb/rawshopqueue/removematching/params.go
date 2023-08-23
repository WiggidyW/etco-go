package removematching

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type ShopQueueRemoveMatchingParams []string

// TODO: This should actually only invalidate the locations that the appraisals are for
func (ShopQueueRemoveMatchingParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.ShopQueueReadCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
