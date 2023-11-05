package bucket

import (
	"context"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_S_TYPEMAPSBUILDER_BUF_CAP          int           = 0
	WEB_S_TYPEMAPSBUILDER_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_S_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_S_TYPEMAPSBUILDER_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebShopLocationTypeMapsBuilder = localcache.RegisterType[map[b.TypeId]b.WebShopLocationTypeBundle](WEB_S_TYPEMAPSBUILDER_BUF_CAP)
}

func GetWebShopLocationTypeMapsBuilder(
	ctx context.Context,
) (
	rep map[b.TypeId]b.WebShopLocationTypeBundle,
	expires *time.Time,
	err error,
) {
	return webGet(
		ctx,
		client.ReadWebShopLocationTypeMapsBuilder,
		keys.TypeStrWebShopLocationTypeMapsBuilder,
		keys.CacheKeyWebShopLocationTypeMapsBuilder,
		WEB_S_TYPEMAPSBUILDER_LOCK_TTL,
		WEB_S_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF,
		WEB_S_TYPEMAPSBUILDER_EXPIRES_IN,
		build.CAPACITY_WEB_SHOP_LOCATION_TYPE_MAPS_BUILDER,
	)
}

func SetWebShopLocationTypeMapsBuilder(
	ctx context.Context,
	rep map[b.TypeId]b.WebShopLocationTypeBundle,
) (
	err error,
) {
	return set(
		ctx,
		client.WriteWebShopLocationTypeMapsBuilder,
		keys.TypeStrWebShopLocationTypeMapsBuilder,
		keys.CacheKeyWebShopLocationTypeMapsBuilder,
		WEB_S_TYPEMAPSBUILDER_LOCK_TTL,
		WEB_S_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF,
		WEB_S_TYPEMAPSBUILDER_EXPIRES_IN,
		rep,
	)
}
