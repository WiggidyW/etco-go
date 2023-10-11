package remotedb

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
)

type CancelPurchaseParams struct {
	CharacterId   int32
	AppraisalCode string
}

func (p CancelPurchaseParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.ReadUserDataCacheKey(
			p.CharacterId,
		),
		cachekeys.ReadShopQueueCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
