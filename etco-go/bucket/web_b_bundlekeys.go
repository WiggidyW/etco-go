package bucket

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_B_BUNDLEKEYS_BUF_CAP                int           = 0
	WEB_B_BUNDLEKEYS_READ_SLOCK_TTL         time.Duration = 1 * time.Minute
	WEB_B_BUNDLEKEYS_READ_SLOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
)

func init() {
	keys.TypeStrWebBuybackBundleKeys = cache.RegisterType[map[b.BundleKey]struct{}]("webbuybackbundlekeys", WEB_B_BUNDLEKEYS_BUF_CAP)
}

func GetWebBuybackBundleKeys(
	x cache.Context,
) (
	rep map[b.BundleKey]struct{},
	expires time.Time,
	err error,
) {
	return bundleKeysGet(
		x,
		GetWebBuybackSystemTypeMapsBuilder,
		keys.CacheKeyWebBuybackBundleKeys,
		keys.TypeStrWebBuybackBundleKeys,
	)
}

func ProtoGetWebBuybackBundleKeys(
	x cache.Context,
) (
	rep []string,
	expires time.Time,
	err error,
) {
	var m map[b.BundleKey]struct{}
	m, expires, err = GetWebBuybackBundleKeys(x)
	if err == nil {
		rep = mapToKeySlice(m)
	}
	return rep, expires, err
}
