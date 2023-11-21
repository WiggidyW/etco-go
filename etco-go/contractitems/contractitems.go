package contractitems

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/esi"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
	"github.com/WiggidyW/etco-go/proto"
	"github.com/WiggidyW/etco-go/protoregistry"
)

const (
	CONTRACT_ITEMS_BUF_CAP int = 0
)

func init() {
	keys.TypeStrContractItems = cache.RegisterType[[]ContractItem]("contractitems", CONTRACT_ITEMS_BUF_CAP)
}

type ContractItem struct {
	Quantity int64
	TypeId   int32
}

func (ci ContractItem) GetTypeId() int32   { return ci.TypeId }
func (ci ContractItem) GetQuantity() int64 { return ci.Quantity }
func (ci ContractItem) ToProto(
	r *protoregistry.ProtoRegistry,
) *proto.NamedBasicItem {
	return &proto.NamedBasicItem{
		TypeId:   r.AddTypeById(ci.TypeId),
		Quantity: ci.Quantity,
	}
}

type ContractItems = []ContractItem

func fromEntries(entries []esi.ContractItemsEntry) []ContractItem {
	items := make([]ContractItem, 0, len(entries))
	itemsMap := make(map[int32]int64, len(entries))
	for _, entry := range entries {
		itemsMap[entry.TypeId] += int64(entry.Quantity)
	}
	for typeId, quantity := range itemsMap {
		items = append(items, ContractItem{quantity, typeId})
	}
	return items
}

func GetContractItems(x cache.Context, contractId int32) (
	items []ContractItem,
	expires time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyContractItems(contractId)
	return fetch.FetchWithCache(
		x,
		getContractItemsFetchFunc(contractId, cacheKey),
		cacheprefetch.WeakCache(
			cacheKey,
			keys.TypeStrContractItems,
			getContractItemsNewRep,
			cache.SloshTrue[[]ContractItem],
			nil,
		),
	)
}

func ProtoGetContractItems(
	x cache.Context,
	r *protoregistry.ProtoRegistry,
	contractId int32,
) (
	items []*proto.NamedBasicItem,
	expires time.Time,
	err error,
) {
	var rItems []ContractItem
	rItems, expires, err = GetContractItems(x, contractId)
	if err == nil {
		items = proto.P1ToProtoMany(rItems, r)
	}
	return items, expires, err
}

func getContractItemsNewRep() []ContractItem {
	return make([]ContractItem, 0, esi.CONTRACT_ITEMS_ENTRIES_MAKE_CAP)
}

func getContractItemsFetchFunc(
	contractId int32,
	cacheKey string,
) fetch.CachingFetch[[]ContractItem] {
	return func(x cache.Context) (
		rep []ContractItem,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		var entries []esi.ContractItemsEntry
		entries, expires, err = esi.GetContractItemsEntries(x, contractId)
		if err != nil {
			return nil, expires, nil, err
		}
		if entries != nil {
			rep = fromEntries(entries)
		}
		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.DualSetOne[[]ContractItem](
				cacheKey,
				keys.TypeStrContractItems,
				rep,
				expires,
			),
		}
		return rep, expires, postFetch, nil
	}
}
