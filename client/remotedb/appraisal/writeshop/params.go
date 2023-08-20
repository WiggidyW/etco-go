package writeshop

import (
	a "github.com/WiggidyW/weve-esi/client/appraisal"
	"github.com/WiggidyW/weve-esi/client/cachekeys"
)

type WriteShopPurchaseParams struct {
	CharacterId   int32
	AppraisalCode string
	Appraisal     a.ShopAppraisal
}

func (p WriteShopPurchaseParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.ReadCharacterAppraisalCodesCacheKey(p.CharacterId),
		cachekeys.ShopQueueReadCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
