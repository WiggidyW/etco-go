package unreserved

import "github.com/WiggidyW/eve-trading-co-go/client/cachekeys"

const (
	CACHE_KEY string = "unres-assets"
)

type UnreservedShopAssetsParams struct {
	ShopQueue []string
}

func (UnreservedShopAssetsParams) CacheKey() string {
	return cachekeys.UnreservedShopAssetsCacheKey()
}
