package get

import (
	"github.com/WiggidyW/etco-go/client/authingfwding"
	a "github.com/WiggidyW/etco-go/client/authingfwding/authing"
	"github.com/WiggidyW/etco-go/client/caching"
	cfg "github.com/WiggidyW/etco-go/client/configure"
	fkbucketreader "github.com/WiggidyW/etco-go/client/configure/internal/fixedkeybucket/reader"
)

type A_GetShopLocationTypeMapsBuilderClient = a.AuthingClient[
	authingfwding.WithAuthableParams[struct{}],
	struct{},
	caching.CachingResponse[cfg.ShopLocationTypeMapsBuilder],
	GetShopLocationTypeMapsBuilderClient,
]

type GetShopLocationTypeMapsBuilderClient = fkbucketreader.
	SC_FixedKeyBucketReaderClient[cfg.ShopLocationTypeMapsBuilder]
