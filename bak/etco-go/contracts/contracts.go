package contracts

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
)

const (
	CONTRACTS_BUF_CAP int = 0
)

func init() {
	keys.TypeStrContracts = localcache.RegisterType[Contracts]("contracts", CONTRACTS_BUF_CAP)
}

func GetContracts(x cache.Context) (
	rep Contracts,
	expires time.Time,
	err error,
) {
	return get(x)
}

func GetShopContracts(x cache.Context) (
	rep map[string]Contract,
	expires time.Time,
	err error,
) {
	var contracts Contracts
	contracts, expires, err = GetContracts(x)
	if err != nil {
		return nil, expires, err
	}
	rep = contracts.ShopContracts
	return rep, expires, err
}

func GetBuybackContracts(x cache.Context) (
	rep map[string]Contract,
	expires time.Time,
	err error,
) {
	var contracts Contracts
	contracts, expires, err = GetContracts(x)
	if err != nil {
		return nil, expires, err
	}
	rep = contracts.BuybackContracts
	return rep, expires, err
}
