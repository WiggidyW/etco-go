package shopassets

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
)

const (
	RAW_BUF_CAP        int = 0
	UNRESERVED_BUF_CAP int = 0
)

func init() {
	keys.TypeStrNSRawShopAssets = localcache.RegisterType[struct{}]("rawshopassets", 0)
	keys.TypeStrRawShopAssets = localcache.RegisterType[map[int32]int64]("shopassets", RAW_BUF_CAP)
	keys.TypeStrUnreservedShopAssets = localcache.RegisterType[map[int32]int64]("unreservedshopassets", UNRESERVED_BUF_CAP)
}

func getRawShopAssets(
	x cache.Context,
	locationId int64,
) (
	rep map[int32]int64,
	expires time.Time,
	err error,
) {
	return rawShopAssetsGet(x, locationId)
}

func GetUnreservedShopAssets(
	x cache.Context,
	locationId int64,
) (
	assets map[int32]int64,
	expires time.Time,
	err error,
) {
	return unreservedShopAssetsGet(x, locationId)
}
