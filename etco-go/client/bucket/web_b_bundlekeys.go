package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	SC "github.com/WiggidyW/etco-go/client/caching/strong/caching"
)

const (
	WEB_B_BUNDLEKEYS_MIN_EXPIRES    time.Duration = 0
	WEB_B_BUNDLEKEYS_SLOCK_TTL      time.Duration = 1 * time.Minute
	WEB_B_BUNDLEKEYS_SLOCK_MAX_WAIT time.Duration = 1 * time.Minute
)

type WebBuybackBundleKeysParams struct{}

func (p WebBuybackBundleKeysParams) CacheKey() string {
	return cachekeys.WebBuybackBundleKeysCacheKey()
}

type SC_WebBuybackBundleKeysClient = SC.StrongCachingClient[
	WebBuybackBundleKeysParams,
	map[string]struct{},
	cache.ExpirableData[map[string]struct{}],
	WebBuybackBundleKeysClient,
]

func NewSC_WebBuybackBundleKeysClient(
	webBTypeMapsBuilderReaderClient SC_WebBuybackSystemTypeMapsBuilderReaderClient,
	sCache cache.SharedServerCache,
) SC_WebBuybackBundleKeysClient {
	return SC.NewStrongCachingClient(
		NewWebBuybackBundleKeysClient(webBTypeMapsBuilderReaderClient),
		WEB_B_BUNDLEKEYS_MIN_EXPIRES,
		sCache,
		WEB_B_BUNDLEKEYS_SLOCK_TTL,
		WEB_B_BUNDLEKEYS_SLOCK_MAX_WAIT,
	)
}

type WebBuybackBundleKeysClient struct {
	webBTypeMapsBuilderReaderClient SC_WebBuybackSystemTypeMapsBuilderReaderClient
}

func NewWebBuybackBundleKeysClient(
	webBTypeMapsBuilderReaderClient SC_WebBuybackSystemTypeMapsBuilderReaderClient,
) WebBuybackBundleKeysClient {
	return WebBuybackBundleKeysClient{webBTypeMapsBuilderReaderClient}
}

func (wbbkc WebBuybackBundleKeysClient) Fetch(
	ctx context.Context,
	params WebBuybackBundleKeysParams,
) (
	rep *cache.ExpirableData[map[string]struct{}],
	err error,
) {
	builderRep, err := wbbkc.webBTypeMapsBuilderReaderClient.Fetch(
		ctx,
		WebBuybackSystemTypeMapsBuilderReaderParams{},
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
