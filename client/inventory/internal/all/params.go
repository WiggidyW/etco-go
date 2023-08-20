package all

import "github.com/WiggidyW/weve-esi/client/cachekeys"

type AllShopAssetsParams struct{}

func (p AllShopAssetsParams) CacheKey() string {
	return cachekeys.AllShopAssetsCacheKey()
}
