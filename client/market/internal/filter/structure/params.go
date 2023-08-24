package filterstructure

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
)

type FilterStructureMarketParams struct {
	WebRefreshToken string
	StructureId     int64
	TypeId          int32
	IsBuy           bool
}

func (p FilterStructureMarketParams) CacheKey() string {
	return cachekeys.FilterStructureMarketCacheKey(
		p.StructureId,
		p.TypeId,
		p.IsBuy,
	)
}
