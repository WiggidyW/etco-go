package remotedb

import (
	"github.com/WiggidyW/etco-go/client/cachekeys"
)

type DelPurchasesParams struct {
	AppraisalCodes []string
}

// TODO: This should actually only invalidate the locations that the appraisals are for
func (DelPurchasesParams) AntiCacheKeys() []string {
	return []string{
		cachekeys.ReadShopQueueCacheKey(),
		cachekeys.UnreservedShopAssetsCacheKey(),
	}
}
