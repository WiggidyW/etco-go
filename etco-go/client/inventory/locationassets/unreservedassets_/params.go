package unreservedassets_

import "github.com/WiggidyW/etco-go/client/cachekeys"

type UnreservedShopAssetsParams struct {
	ShopQueue []string
}

func (UnreservedShopAssetsParams) CacheKey() string {
	return cachekeys.UnreservedShopAssetsCacheKey()
}
