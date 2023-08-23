package get

import (
	"github.com/WiggidyW/weve-esi/client/authingfwding"
	a "github.com/WiggidyW/weve-esi/client/authingfwding/authing"
	"github.com/WiggidyW/weve-esi/client/caching"
	cfg "github.com/WiggidyW/weve-esi/client/configure"
	fkbucketreader "github.com/WiggidyW/weve-esi/client/configure/internal/fixedkeybucket/reader"
)

type A_GetShopLocationTypeMapsBuilderClient = a.AuthingClient[
	authingfwding.WithAuthableParams[struct{}],
	struct{},
	caching.CachingResponse[cfg.ShopLocationTypeMapsBuilder],
	GetShopLocationTypeMapsBuilderClient,
]

type GetShopLocationTypeMapsBuilderClient = fkbucketreader.
	SC_FixedKeyBucketReaderClient[cfg.ShopLocationTypeMapsBuilder]
