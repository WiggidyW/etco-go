package cancel

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type CancelShopPurchaseParams struct {
	CharacterId   int32
	AppraisalCode string
}

func (p CancelShopPurchaseParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.ReadUserDataCacheKey(
			p.CharacterId,
		),
		cachekeys.ShopQueueReadCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
