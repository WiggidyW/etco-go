package contracts

import "github.com/WiggidyW/eve-trading-co-go/client/cachekeys"

type ContractsParams struct{}

func (p ContractsParams) CacheKey() string {
	// return fmt.Sprintf("contracts-%d", p.CorporationId)
	return cachekeys.ContractsCacheKey()
}
