package bucket

import (
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_MARKETS_BUF_CAP          int           = 0
	WEB_MARKETS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_MARKETS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_MARKETS_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebMarkets = cache.RegisterType[map[b.MarketName]b.WebMarket]("webmarkets", WEB_MARKETS_BUF_CAP)
}

func GetWebMarkets(
	x cache.Context,
) (
	rep map[b.MarketName]b.WebMarket,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebMarkets,
		keys.CacheKeyWebMarkets,
		keys.TypeStrWebMarkets,
		WEB_MARKETS_EXPIRES_IN,
		build.CAPACITY_WEB_MARKETS,
	)
}

func SetWebMarkets(
	x cache.Context,
	rep map[b.MarketName]b.WebMarket,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebMarkets,
		keys.CacheKeyWebMarkets,
		keys.TypeStrWebMarkets,
		WEB_MARKETS_EXPIRES_IN,
		rep,
		nil,
	)
}
