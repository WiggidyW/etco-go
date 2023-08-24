package cancel

import (
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
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
