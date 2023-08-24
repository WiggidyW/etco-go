package read

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
)

type ShopQueueReadParams struct{}

func (ShopQueueReadParams) CacheKey() string {
	return cachekeys.ShopQueueReadCacheKey()
}
