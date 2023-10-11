package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SC "github.com/WiggidyW/etco-go/client/caching/strong/caching"
)

const (
	WEB_S_BUNDLEKEYS_MIN_EXPIRES    time.Duration = 0
	WEB_S_BUNDLEKEYS_SLOCK_TTL      time.Duration = 1 * time.Minute
	WEB_S_BUNDLEKEYS_SLOCK_MAX_WAIT time.Duration = 1 * time.Minute
)

type WebShopBundleKeysParams struct{}

func (p WebShopBundleKeysParams) CacheKey() string {
	return cachekeys.WebShopBundleKeysCacheKey()
}

type SC_WebShopBundleKeysClient = SC.StrongCachingClient[
	WebShopBundleKeysParams,
	map[string]struct{},
	cache.ExpirableData[map[string]struct{}],
	WebShopBundleKeysClient,
]

func NewSC_WebShopBundleKeysClient(
	webSTypeMapsBuilderReaderClient SC_WebShopLocationTypeMapsBuilderReaderClient,
	sCache cache.SharedServerCache,
) SC_WebShopBundleKeysClient {
	return SC.NewStrongCachingClient(
		NewWebShopBundleKeysClient(webSTypeMapsBuilderReaderClient),
		WEB_S_BUNDLEKEYS_MIN_EXPIRES,
		sCache,
		WEB_S_BUNDLEKEYS_SLOCK_TTL,
		WEB_S_BUNDLEKEYS_SLOCK_MAX_WAIT,
	)
}

type WebShopBundleKeysClient struct {
	webSTypeMapsBuilderReaderClient SC_WebShopLocationTypeMapsBuilderReaderClient
}

func NewWebShopBundleKeysClient(
	webSTypeMapsBuilderReaderClient SC_WebShopLocationTypeMapsBuilderReaderClient,
) WebShopBundleKeysClient {
	return WebShopBundleKeysClient{webSTypeMapsBuilderReaderClient}
}

func (wsbkc WebShopBundleKeysClient) Fetch(
	ctx context.Context,
	params WebShopBundleKeysParams,
) (
	rep *cache.ExpirableData[map[string]struct{}],
	err error,
) {
	builderRep, err := wsbkc.webSTypeMapsBuilderReaderClient.Fetch(
		ctx,
		WebShopLocationTypeMapsBuilderReaderParams{},
	)
	if err != nil {
		return nil, err
	}

	bundleKeys := extractBuilderBundleKeys(builderRep.Data())

	return cache.NewExpirableDataPtr(
		bundleKeys,
		builderRep.Expires(),
	), nil
}
