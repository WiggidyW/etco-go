package single

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type RateLimitingContractItemsParams struct {
	ContractId int32
}

func (p RateLimitingContractItemsParams) CacheKey() string {
	return cachekeys.RateLimitingContractItemsCacheKey(p.ContractId)
}
