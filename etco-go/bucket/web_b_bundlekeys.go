package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_B_BUNDLEKEYS_BUF_CAP                int           = 0
	WEB_B_BUNDLEKEYS_READ_SLOCK_TTL         time.Duration = 1 * time.Minute
	WEB_B_BUNDLEKEYS_READ_SLOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

func init() {
	keys.TypeStrWebBuybackBundleKeys = localcache.RegisterType[map[b.BundleKey]struct{}](WEB_B_BUNDLEKEYS_BUF_CAP)
}

func GetWebBuybackBundleKeys(
	ctx context.Context,
) (
	rep map[b.BundleKey]struct{},
	expires *time.Time,
	err error,
) {
	return bundleKeysGet(
		ctx,
		GetWebBuybackSystemTypeMapsBuilder,
		keys.TypeStrWebBuybackBundleKeys,
		keys.CacheKeyWebBuybackBundleKeys,
		WEB_B_BUNDLEKEYS_READ_SLOCK_TTL,
		WEB_B_BUNDLEKEYS_READ_SLOCK_MAX_BACKOFF,
	)
}
