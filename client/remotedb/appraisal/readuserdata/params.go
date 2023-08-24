package readuserdata

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
)

type ReadUserDataParams struct {
	CharacterId int32
}

func (p ReadUserDataParams) CacheKey() string {
	return cachekeys.ReadUserDataCacheKey(p.CharacterId)
}
