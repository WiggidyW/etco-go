package read

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type ShopQueueReadParams struct{}

func (ShopQueueReadParams) CacheKey() string {
	return cachekeys.ShopQueueReadCacheKey()
}
