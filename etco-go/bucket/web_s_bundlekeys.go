package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_S_BUNDLEKEYS_BUF_CAP          int           = 0
	WEB_S_BUNDLEKEYS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_S_BUNDLEKEYS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

func init() {
	keys.TypeStrWebShopBundleKeys = localcache.RegisterType[map[b.BundleKey]struct{}](WEB_S_BUNDLEKEYS_BUF_CAP)
}

func GetWebShopBundleKeys(
	ctx context.Context,
) (
	rep map[b.BundleKey]struct{},
	expires *time.Time,
	err error,
) {
	return bundleKeysGet(
		ctx,
		GetWebShopLocationTypeMapsBuilder,
		keys.TypeStrWebShopBundleKeys,
		keys.CacheKeyWebShopBundleKeys,
		WEB_S_BUNDLEKEYS_LOCK_TTL,
		WEB_S_BUNDLEKEYS_LOCK_MAX_BACKOFF,
	)
}
