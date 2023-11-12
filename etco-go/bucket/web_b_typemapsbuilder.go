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
	WEB_B_TYPEMAPSBUILDER_BUF_CAP          int           = 0
	WEB_B_TYPEMAPSBUILDER_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_B_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_B_TYPEMAPSBUILDER_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebBuybackSystemTypeMapsBuilder = cache.RegisterType[map[b.TypeId]b.WebBuybackSystemTypeBundle]("webbuybacksystemtypemapsbuilder", WEB_B_TYPEMAPSBUILDER_BUF_CAP)
}

func GetWebBuybackSystemTypeMapsBuilder(
	x cache.Context,
) (
	rep map[b.TypeId]b.WebBuybackSystemTypeBundle,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebBuybackSystemTypeMapsBuilder,
		keys.CacheKeyWebBuybackSystemTypeMapsBuilder,
		keys.TypeStrWebBuybackSystemTypeMapsBuilder,
		WEB_B_TYPEMAPSBUILDER_EXPIRES_IN,
		build.CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
	)
}

func SetWebBuybackSystemTypeMapsBuilder(
	x cache.Context,
	rep map[b.TypeId]b.WebBuybackSystemTypeBundle,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebBuybackSystemTypeMapsBuilder,
		keys.CacheKeyWebBuybackSystemTypeMapsBuilder,
		keys.TypeStrWebBuybackSystemTypeMapsBuilder,
		WEB_B_TYPEMAPSBUILDER_EXPIRES_IN,
		rep,
		prefetch.ServerCacheLockPtr(
			keys.CacheKeyWebBuybackBundleKeys,
			keys.TypeStrWebBuybackBundleKeys,
		),
	)
}
