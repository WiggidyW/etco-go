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
	WEB_SHOP_LOCATIONS_BUF_CAP          int           = 0
	WEB_SHOP_LOCATIONS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_SHOP_LOCATIONS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_SHOP_LOCATIONS_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebShopLocations = localcache.RegisterType[map[b.LocationId]b.WebShopLocation](WEB_SHOP_LOCATIONS_BUF_CAP)
}

func GetWebShopLocations(
	ctx context.Context,
) (
	rep map[b.LocationId]b.WebShopLocation,
	expires *time.Time,
	err error,
) {
	return webGet(
		ctx,
		client.ReadWebShopLocations,
		keys.TypeStrWebShopLocations,
		keys.CacheKeyWebShopLocations,
		WEB_SHOP_LOCATIONS_LOCK_TTL,
		WEB_SHOP_LOCATIONS_LOCK_MAX_BACKOFF,
		WEB_SHOP_LOCATIONS_EXPIRES_IN,
		build.CAPACITY_WEB_SHOP_LOCATIONS,
	)
}

func SetWebShopLocations(
	ctx context.Context,
	rep map[b.LocationId]b.WebShopLocation,
) (
	err error,
) {
	return set(
		ctx,
		client.WriteWebShopLocations,
		keys.TypeStrWebShopLocations,
		keys.CacheKeyWebShopLocations,
		WEB_SHOP_LOCATIONS_LOCK_TTL,
		WEB_SHOP_LOCATIONS_LOCK_MAX_BACKOFF,
		WEB_SHOP_LOCATIONS_EXPIRES_IN,
		rep,
	)
}
