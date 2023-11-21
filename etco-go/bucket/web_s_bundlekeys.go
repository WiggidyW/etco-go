package bucket

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_S_BUNDLEKEYS_BUF_CAP int = 0
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

func ProtoGetWebShopBundleKeys(
	x cache.Context,
) (
	rep []string,
	expires time.Time,
	err error,
) {
	var m map[b.BundleKey]struct{}
	m, expires, err = GetWebShopBundleKeys(x)
	if err == nil {
		rep = mapToKeySlice(m)
	}
	return rep, expires, err
}
