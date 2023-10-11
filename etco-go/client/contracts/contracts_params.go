package contracts

import "github.com/WiggidyW/etco-go/client/cachekeys"

type ContractsParams struct{}

func (p ContractsParams) CacheKey() string {
	// return fmt.Sprintf("contracts-%d", p.CorporationId)
	return cachekeys.ContractsCacheKey()
}
