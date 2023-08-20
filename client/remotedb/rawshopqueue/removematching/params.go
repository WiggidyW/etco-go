package removematching

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type ShopQueueRemoveMatchingParams []string

func (ShopQueueRemoveMatchingParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.ShopQueueReadCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
