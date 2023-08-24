package all

import "github.com/WiggidyW/eve-trading-co-go/client/cachekeys"

type AllShopAssetsParams struct{}

func (p AllShopAssetsParams) CacheKey() string {
	return cachekeys.AllShopAssetsCacheKey()
}
