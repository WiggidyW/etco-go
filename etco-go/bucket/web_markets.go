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
	WEB_MARKETS_BUF_CAP          int           = 0
	WEB_MARKETS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_MARKETS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_MARKETS_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebMarkets = localcache.RegisterType[map[b.MarketName]b.WebMarket](WEB_MARKETS_BUF_CAP)
}

func GetWebMarkets(
	ctx context.Context,
) (
	rep map[b.MarketName]b.WebMarket,
	expires *time.Time,
	err error,
) {
	return webGet(
		ctx,
		client.ReadWebMarkets,
		keys.TypeStrWebMarkets,
		keys.CacheKeyWebMarkets,
		WEB_MARKETS_LOCK_TTL,
		WEB_MARKETS_LOCK_MAX_BACKOFF,
		WEB_MARKETS_EXPIRES_IN,
		build.CAPACITY_WEB_MARKETS,
	)
}

func SetWebMarkets(
	ctx context.Context,
	rep map[b.MarketName]b.WebMarket,
) (
	err error,
) {
	return set(
		ctx,
		client.WriteWebMarkets,
		keys.TypeStrWebMarkets,
		keys.CacheKeyWebMarkets,
		WEB_MARKETS_LOCK_TTL,
		WEB_MARKETS_LOCK_MAX_BACKOFF,
		WEB_MARKETS_EXPIRES_IN,
		rep,
	)
}
