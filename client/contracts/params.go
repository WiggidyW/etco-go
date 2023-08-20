package contracts

import "github.com/WiggidyW/weve-esi/client/cachekeys"

type ContractsParams struct{}

func (p ContractsParams) CacheKey() string {
	// return fmt.Sprintf("contracts-%d", p.CorporationId)
	return cachekeys.ContractsCacheKey()
}
