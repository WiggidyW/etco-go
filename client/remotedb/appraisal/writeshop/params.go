package writeshop

import (
	"github.com/WiggidyW/weve-esi/client/cachekeys"
	a "github.com/WiggidyW/weve-esi/client/remotedb/appraisal"
)

type WriteShopPurchaseParams[
	S a.IShopAppraisal[I],
	I a.IShopItem,
] struct {
	CharacterId   int32
	AppraisalCode string
	IAppraisal    S
}

func (p WriteShopPurchaseParams[S, I]) AntiCacheKeys() []string {
	return []string{
		cachekeys.ReadCharacterAppraisalCodesCacheKey(p.CharacterId),
		cachekeys.ShopQueueReadCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
