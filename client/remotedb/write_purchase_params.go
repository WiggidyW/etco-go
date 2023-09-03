package remotedb

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
	rdb "github.com/WiggidyW/etco-go/remotedb"
)

type WritePurchaseParams struct {
	Appraisal rdb.ShopAppraisal
}

func (p WritePurchaseParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.ReadUserDataCacheKey(
			p.Appraisal.CharacterId,
		),
		cachekeys.ReadShopQueueCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
