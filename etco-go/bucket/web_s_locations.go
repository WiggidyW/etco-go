package bucket

import (
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_SHOP_LOCATIONS_BUF_CAP          int           = 0
	WEB_SHOP_LOCATIONS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_SHOP_LOCATIONS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_SHOP_LOCATIONS_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebShopLocations = cache.RegisterType[map[b.LocationId]b.WebShopLocation]("webshoplocations", WEB_SHOP_LOCATIONS_BUF_CAP)
}

func GetWebShopLocations(
	x cache.Context,
) (
	rep map[b.LocationId]b.WebShopLocation,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebShopLocations,
		keys.CacheKeyWebShopLocations,
		keys.TypeStrWebShopLocations,
		WEB_SHOP_LOCATIONS_EXPIRES_IN,
		build.CAPACITY_WEB_SHOP_LOCATIONS,
	)
}

func SetWebShopLocations(
	x cache.Context,
	rep map[b.LocationId]b.WebShopLocation,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebShopLocations,
		keys.CacheKeyWebShopLocations,
		keys.TypeStrWebShopLocations,
		WEB_SHOP_LOCATIONS_EXPIRES_IN,
		rep,
		nil,
	)
}
