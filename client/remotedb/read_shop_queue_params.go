package remotedb

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
)

type ReadShopQueueParams struct{}

func (ReadShopQueueParams) CacheKey() string {
	return cachekeys.ReadShopQueueCacheKey()
}
