package bucket

import (
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch/prefetch"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_S_TYPEMAPSBUILDER_BUF_CAP          int           = 0
	WEB_S_TYPEMAPSBUILDER_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_S_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_S_TYPEMAPSBUILDER_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebShopLocationTypeMapsBuilder = cache.RegisterType[map[b.TypeId]b.WebShopLocationTypeBundle]("webshoplocationtypemapsbuilder", WEB_S_TYPEMAPSBUILDER_BUF_CAP)
}

func GetWebShopLocationTypeMapsBuilder(
	x cache.Context,
) (
	rep map[b.TypeId]b.WebShopLocationTypeBundle,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebShopLocationTypeMapsBuilder,
		keys.TypeStrWebShopLocationTypeMapsBuilder,
		keys.CacheKeyWebShopLocationTypeMapsBuilder,
		WEB_S_TYPEMAPSBUILDER_EXPIRES_IN,
		build.CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
	)
}

func SetWebShopLocationTypeMapsBuilder(
	x cache.Context,
	rep map[b.TypeId]b.WebShopLocationTypeBundle,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebShopLocationTypeMapsBuilder,
		keys.CacheKeyWebShopLocationTypeMapsBuilder,
		keys.TypeStrWebShopLocationTypeMapsBuilder,
		WEB_S_TYPEMAPSBUILDER_EXPIRES_IN,
		rep,
		prefetch.ServerCacheLockPtr(
			keys.CacheKeyWebShopBundleKeys,
			keys.TypeStrWebShopBundleKeys,
		),
	)
}
