package write

import (
	a "github.com/WiggidyW/eve-trading-co-go/client/appraisal"
	"github.com/WiggidyW/eve-trading-co-go/client/cachekeys"
)

type WriteShopPurchaseParams struct {
	Appraisal a.ShopAppraisal
}

func (p WriteShopPurchaseParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.ReadUserDataCacheKey(
			p.Appraisal.CharacterId,
		),
		cachekeys.ShopQueueReadCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
