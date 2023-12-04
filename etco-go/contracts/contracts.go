package contracts

import (
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/contractitems"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/kind"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
)

const (
	BUYBACK_CONTRACTS_BUF_CAP int = 0
	SHOP_CONTRACTS_BUF_CAP    int = 0
	HAUL_CONTRACTS_BUF_CAP    int = 0
)

func init() {
	keys.TypeStrNSContracts = cache.RegisterType[Contracts]("contracts", 0)
	keys.TypeStrBuybackContracts = cache.RegisterType[map[string]Contract]("buybackcontracts", BUYBACK_CONTRACTS_BUF_CAP)
	keys.TypeStrShopContracts = cache.RegisterType[map[string]Contract]("shopcontracts", SHOP_CONTRACTS_BUF_CAP)
	keys.TypeStrHaulContracts = cache.RegisterType[map[string]Contract]("haulcontracts", HAUL_CONTRACTS_BUF_CAP)
}

func getContracts(
	x cache.Context,
	cacheKey, typeStr keys.Key,
	storeKind kind.StoreKind,
) (
	rep map[string]Contract,
	expires time.Time,
	err error,
) {
	return fetch.FetchWithCache(
		x,
		getContractsFetchFunc(storeKind),
		cacheprefetch.WeakMultiCacheKnownKeys(
			cacheKey,
			typeStr,
			keys.CacheKeyNSContracts,
			keys.TypeStrNSContracts,
			nil,
			cache.SloshTrue[map[string]Contract],
			[]cacheprefetch.ActionOrderedLocks{{
				Locks: []cacheprefetch.ActionLock{
					cacheprefetch.DualLock(
						keys.CacheKeyBuybackContracts,
						keys.TypeStrBuybackContracts,
					),
					cacheprefetch.DualLock(
						keys.CacheKeyShopContracts,
						keys.TypeStrShopContracts,
					),
					cacheprefetch.DualLock(
						keys.CacheKeyHaulContracts,
						keys.TypeStrHaulContracts,
					),
				},
				Child: nil,
			}},
		),
	)
}

func getContractsFetchFunc(
	storeKind kind.StoreKind,
) fetch.CachingFetch[map[string]Contract] {
	return func(x cache.Context) (
		rep map[string]Contract,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		x, cancel := x.WithCancel()
		defer cancel()

		var repOrStream esi.RepOrStream[esi.ContractsEntry]
		var pages int
		repOrStream, expires, pages, err = esi.GetContractsEntries(x)
		if err != nil {
			return nil, expires, nil, err
		}

		contracts := newContracts()
		if repOrStream.Rep != nil {
			contracts.filterAddEntries(*repOrStream.Rep)
		} else /* if repOrStream.Stream != nil */ {
			var entries []esi.ContractsEntry
			for i := 0; i < pages; i++ {
				entries, expires, err = repOrStream.Stream.RecvExpMin(expires)
				if err != nil {
					return nil, expires, nil, err
				} else {
					contracts.filterAddEntries(entries)
				}
			}
		}

		if build.BUYBACK_CONTRACT_NOTIFICATIONS ||
			build.SHOP_CONTRACT_NOTIFICATIONS ||
			build.HAUL_CONTRACT_NOTIFICATIONS {
			go func() {
				logger.MaybeErr(getAndNotifyNewContracts(
					x.Background(),
					contracts,
				))
			}()
		}

		switch storeKind {
		case kind.Buyback:
			rep = contracts.BuybackContracts
		case kind.Shop:
			rep = contracts.ShopContracts
		case kind.Haul:
			rep = contracts.HaulContracts
		}
		postFetch = &cachepostfetch.Params{
			Set: []cachepostfetch.ActionSet{
				cachepostfetch.DualSet[map[string]Contract](
					keys.CacheKeyBuybackContracts,
					keys.TypeStrBuybackContracts,
					contracts.BuybackContracts,
					expires,
				),
				cachepostfetch.DualSet[map[string]Contract](
					keys.CacheKeyShopContracts,
					keys.TypeStrShopContracts,
					contracts.ShopContracts,
					expires,
				),
				cachepostfetch.DualSet[map[string]Contract](
					keys.CacheKeyHaulContracts,
					keys.TypeStrHaulContracts,
					contracts.HaulContracts,
					expires,
				),
			},
		}
		return rep, expires, postFetch, nil
	}
}

func GetShopContracts(x cache.Context) (
	contracts map[string]Contract,
	expires time.Time,
	err error,
) {
	return getContracts(
		x,
		keys.CacheKeyShopContracts,
		keys.TypeStrShopContracts,
		kind.Shop,
	)
}

func GetBuybackContracts(x cache.Context) (
	contracts map[string]Contract,
	expires time.Time,
	err error,
) {
	return getContracts(
		x,
		keys.CacheKeyBuybackContracts,
		keys.TypeStrBuybackContracts,
		kind.Buyback,
	)
}

func GetHaulContracts(x cache.Context) (
	contracts map[string]Contract,
	expires time.Time,
	err error,
) {
	return getContracts(
		x,
		keys.CacheKeyHaulContracts,
		keys.TypeStrHaulContracts,
		kind.Haul,
	)
}

func getContract(
	x cache.Context,
	code string,
	getContracts func(cache.Context) (map[string]Contract, time.Time, error),
) (
	contract *Contract,
	expires time.Time,
	err error,
) {
	var contracts map[string]Contract
	contracts, expires, err = getContracts(x)
	if err == nil {
		if contractVal, ok := contracts[code]; ok {
			contract = &contractVal
		}
	}
	return contract, expires, err
}

func GetShopContract(x cache.Context, code string) (
	contract *Contract,
	expires time.Time,
	err error,
) {
	return getContract(x, code, GetShopContracts)
}

func GetBuybackContract(x cache.Context, code string) (
	contract *Contract,
	expires time.Time,
	err error,
) {
	return getContract(x, code, GetBuybackContracts)
}

func GetHaulContract(x cache.Context, code string) (
	contract *Contract,
	expires time.Time,
	err error,
) {
	return getContract(x, code, GetHaulContracts)
}

func protoGetContract(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	getContract func(cache.Context, string) (*Contract, time.Time, error),
) (
	rep *proto.Contract,
	expires time.Time,
	err error,
) {
	// fetch the contract, returning if nil or error
	var rContract *Contract
	rContract, expires, err = getContract(x, code)
	if err != nil || rContract == nil {
		return nil, expires, err
	}

	// fetch location info
	var startLocationInfo *proto.LocationInfo
	var locationInfo *proto.LocationInfo
	var locationInfoExpires time.Time
	startLocationInfo, locationInfo, locationInfoExpires, err =
		ProtoGetLocationInfo(x, r, *rContract)
	if err != nil {
		return nil, expires, err
	} else {
		expires = fetch.CalcExpires(expires, locationInfoExpires)
	}

	return rContract.ToProto(startLocationInfo, locationInfo), expires, nil
}

func ProtoGetShopContract(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
) (*proto.Contract, time.Time, error) {
	return protoGetContract(x, r, code, GetShopContract)
}

func ProtoGetBuybackContract(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
) (*proto.Contract, time.Time, error) {
	return protoGetContract(x, r, code, GetBuybackContract)
}

func ProtoGetHaulContract(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
) (*proto.Contract, time.Time, error) {
	return protoGetContract(x, r, code, GetHaulContract)
}

type ProtoContractWithItemsRep struct {
	Contract *proto.Contract
	Items    []*proto.NamedBasicItem
}

func protoGetContractWithItems(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
	getContract func(cache.Context, string) (*Contract, time.Time, error),
) (
	rep ProtoContractWithItemsRep,
	expires time.Time,
	err error,
) {
	// fetch the contract, return if nil or error
	var rContract *Contract
	rContract, expires, err = getContract(x, code)
	if err != nil || rContract == nil {
		return rep, expires, err
	}

	// fetch the contract items in a goroutine
	x, cancel := x.WithCancel()
	defer cancel()
	chnItems := expirable.NewChanResult[[]*proto.NamedBasicItem](x.Ctx(), 1, 0)
	go expirable.P3Transceive(
		chnItems,
		x, r, rContract.ContractId,
		contractitems.ProtoGetContractItems,
	)

	// fetch location info
	var startLocationInfo *proto.LocationInfo
	var locationInfo *proto.LocationInfo
	var locationInfoExpires time.Time
	startLocationInfo, locationInfo, locationInfoExpires, err =
		ProtoGetLocationInfo(x, r, *rContract)
	if err != nil {
		return rep, expires, err
	} else {
		expires = fetch.CalcExpires(expires, locationInfoExpires)
	}

	rep.Contract = rContract.ToProto(startLocationInfo, locationInfo)
	rep.Items, expires, err = chnItems.RecvExpMin(expires)
	return rep, expires, err
}

func ProtoGetShopContractWithItems(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
) (ProtoContractWithItemsRep, time.Time, error) {
	return protoGetContractWithItems(x, r, code, GetShopContract)
}

func ProtoGetBuybackContractWithItems(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
) (ProtoContractWithItemsRep, time.Time, error) {
	return protoGetContractWithItems(x, r, code, GetBuybackContract)
}

func ProtoGetHaulContractWithItems(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	code string,
) (ProtoContractWithItemsRep, time.Time, error) {
	return protoGetContractWithItems(x, r, code, GetHaulContract)
}

func ProtoGetLocationInfo(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	contract Contract,
) (
	startInfo *proto.LocationInfo,
	info *proto.LocationInfo,
	expires time.Time,
	err error,
) {
	if contract.StartLocationId == 0 {
		info, expires, err =
			esi.ProtoGetLocationInfo(x, r, contract.StartLocationId)
		return nil, info, expires, err
	}

	x, cancel := x.WithCancel()
	defer cancel()
	infoCOV := esi.ProtoGetLocationInfoCOV(x, r, contract.LocationId)

	startInfo, expires, err =
		esi.ProtoGetLocationInfo(x, r, contract.StartLocationId)
	if err != nil {
		return nil, nil, expires, err
	}

	info, expires, err = infoCOV.RecvExpMin(expires)
	return startInfo, info, expires, err
}

func ProtoGetLocationInfoCOV(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	contract Contract,
) (
	startInfoCOVPtr *expirable.ChanOrValue[*proto.LocationInfo],
	infoCOV expirable.ChanOrValue[*proto.LocationInfo],
) {
	infoCOV = esi.ProtoGetLocationInfoCOV(x, r, contract.LocationId)
	if contract.StartLocationId == 0 {
		return nil, infoCOV
	} else {
		startInfoCOV :=
			esi.ProtoGetLocationInfoCOV(x, r, contract.StartLocationId)
		return &startInfoCOV, infoCOV
	}
}
