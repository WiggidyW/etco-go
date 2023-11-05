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
	WEB_B_TYPEMAPSBUILDER_BUF_CAP          int           = 0
	WEB_B_TYPEMAPSBUILDER_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_B_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_B_TYPEMAPSBUILDER_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebBuybackSystemTypeMapsBuilder = localcache.RegisterType[map[b.TypeId]b.WebBuybackSystemTypeBundle](WEB_B_TYPEMAPSBUILDER_BUF_CAP)
}

func GetWebBuybackSystemTypeMapsBuilder(
	ctx context.Context,
) (
	rep map[b.TypeId]b.WebBuybackSystemTypeBundle,
	expires *time.Time,
	err error,
) {
	return webGet(
		ctx,
		client.ReadWebBuybackSystemTypeMapsBuilder,
		keys.TypeStrWebBuybackSystemTypeMapsBuilder,
		keys.CacheKeyWebBuybackSystemTypeMapsBuilder,
		WEB_B_TYPEMAPSBUILDER_LOCK_TTL,
		WEB_B_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF,
		WEB_B_TYPEMAPSBUILDER_EXPIRES_IN,
		build.CAPACITY_WEB_BUYBACK_SYSTEM_TYPE_MAPS_BUILDER,
	)
}

func SetWebBuybackSystemTypeMapsBuilder(
	ctx context.Context,
	rep map[b.TypeId]b.WebBuybackSystemTypeBundle,
) (
	err error,
) {
	return set(
		ctx,
		client.WriteWebBuybackSystemTypeMapsBuilder,
		keys.TypeStrWebBuybackSystemTypeMapsBuilder,
		keys.CacheKeyWebBuybackSystemTypeMapsBuilder,
		WEB_B_TYPEMAPSBUILDER_LOCK_TTL,
		WEB_B_TYPEMAPSBUILDER_LOCK_MAX_BACKOFF,
		WEB_B_TYPEMAPSBUILDER_EXPIRES_IN,
		rep,
	)
}
