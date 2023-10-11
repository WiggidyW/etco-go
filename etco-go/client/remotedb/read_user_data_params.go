package remotedb

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
)

type ReadUserDataParams struct {
	CharacterId int32
}

func (p ReadUserDataParams) CacheKey() string {
	return cachekeys.ReadUserDataCacheKey(p.CharacterId)
}
