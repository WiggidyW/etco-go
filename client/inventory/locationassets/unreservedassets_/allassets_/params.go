package allassets_

import "github.com/WiggidyW/etco-go/client/cachekeys"

type AllShopAssetsParams struct{}

func (p AllShopAssetsParams) CacheKey() string {
	return cachekeys.AllShopAssetsCacheKey()
}
