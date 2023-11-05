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
	WEB_BUYBACK_SYSTEMS_BUF_CAP          int           = 0
	WEB_BUYBACK_SYSTEMS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_BUYBACK_SYSTEMS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_BUYBACK_SYSTEMS_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebBuybackSystems = localcache.RegisterType[map[b.SystemId]b.WebBuybackSystem](WEB_BUYBACK_SYSTEMS_BUF_CAP)
}

func GetWebBuybackSystems(
	ctx context.Context,
) (
	rep map[b.SystemId]b.WebBuybackSystem,
	expires *time.Time,
	err error,
) {
	return webGet(
		ctx,
		client.ReadWebBuybackSystems,
		keys.TypeStrWebBuybackSystems,
		keys.CacheKeyWebBuybackSystems,
		WEB_BUYBACK_SYSTEMS_LOCK_TTL,
		WEB_BUYBACK_SYSTEMS_LOCK_MAX_BACKOFF,
		WEB_BUYBACK_SYSTEMS_EXPIRES_IN,
		build.CAPACITY_WEB_BUYBACK_SYSTEMS,
	)
}

func SetWebBuybackSystems(
	ctx context.Context,
	rep map[b.SystemId]b.WebBuybackSystem,
) (
	err error,
) {
	return set(
		ctx,
		client.WriteWebBuybackSystems,
		keys.TypeStrWebBuybackSystems,
		keys.CacheKeyWebBuybackSystems,
		WEB_BUYBACK_SYSTEMS_LOCK_TTL,
		WEB_BUYBACK_SYSTEMS_LOCK_MAX_BACKOFF,
		WEB_BUYBACK_SYSTEMS_EXPIRES_IN,
		rep,
	)
}
