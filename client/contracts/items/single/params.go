package single

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
)

type RateLimitingContractItemsParams struct {
	ContractId int32
}

func (p RateLimitingContractItemsParams) CacheKey() string {
	return cachekeys.RateLimitingContractItemsCacheKey(p.ContractId)
}
