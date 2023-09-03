package contracts

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
)

type SingleContractItemsParams struct {
	ContractId int32
}

func (p SingleContractItemsParams) CacheKey() string {
	return cachekeys.ContractItemsCacheKey(p.ContractId)
}
