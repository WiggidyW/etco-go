package write

import (
	a "github.com/WiggidyW/weve-esi/client/appraisal"
	"github.com/WiggidyW/weve-esi/client/cachekeys"
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
