package contractitems

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
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

func GetContractItems(x cache.Context, contractId int32) (
	rep []ContractItem,
	expires time.Time,
	err error,
) {
	return contractItemsGet(x, contractId)
}
