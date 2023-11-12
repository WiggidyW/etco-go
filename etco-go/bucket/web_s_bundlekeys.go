package bucket

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_S_BUNDLEKEYS_BUF_CAP          int           = 0
	WEB_S_BUNDLEKEYS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_S_BUNDLEKEYS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

func init() {
	keys.TypeStrWebShopBundleKeys = cache.RegisterType[map[b.BundleKey]struct{}]("webshopbundlekeys", WEB_S_BUNDLEKEYS_BUF_CAP)
}

func GetWebShopBundleKeys(
	x cache.Context,
) (
	rep map[b.BundleKey]struct{},
	expires time.Time,
	err error,
) {
	return bundleKeysGet(
		x,
		GetWebShopLocationTypeMapsBuilder,
		keys.CacheKeyWebShopBundleKeys,
		keys.TypeStrWebShopBundleKeys,
	)
}
