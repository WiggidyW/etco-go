package unreserved

import "github.com/WiggidyW/weve-esi/client/cachekeys"

const (
	CACHE_KEY string = "unres-assets"
)

type UnreservedShopAssetsParams struct {
	ShopQueue []string
}

func (UnreservedShopAssetsParams) CacheKey() string {
	return cachekeys.UnreservedShopAssetsCacheKey()
}
