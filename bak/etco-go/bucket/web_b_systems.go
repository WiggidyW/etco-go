package bucket

import (
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"

	b "github.com/WiggidyW/etco-go-bucket"
)

const (
	WEB_BUYBACK_SYSTEMS_BUF_CAP          int           = 0
	WEB_BUYBACK_SYSTEMS_LOCK_TTL         time.Duration = 1 * time.Minute
	WEB_BUYBACK_SYSTEMS_LOCK_MAX_BACKOFF time.Duration = 1 * time.Minute
	WEB_BUYBACK_SYSTEMS_EXPIRES_IN       time.Duration = 24 * time.Hour
)

func init() {
	keys.TypeStrWebBuybackSystems = cache.RegisterType[map[b.SystemId]b.WebBuybackSystem]("webbuybacksystems", WEB_BUYBACK_SYSTEMS_BUF_CAP)
}

func GetWebBuybackSystems(
	x cache.Context,
) (
	rep map[b.SystemId]b.WebBuybackSystem,
	expires time.Time,
	err error,
) {
	return webGet(
		x,
		client.ReadWebBuybackSystems,
		keys.CacheKeyWebBuybackSystems,
		keys.TypeStrWebBuybackSystems,
		WEB_BUYBACK_SYSTEMS_EXPIRES_IN,
		build.CAPACITY_WEB_BUYBACK_SYSTEMS,
	)
}

func SetWebBuybackSystems(
	x cache.Context,
	rep map[b.SystemId]b.WebBuybackSystem,
) (
	err error,
) {
	return set(
		x,
		client.WriteWebBuybackSystems,
		keys.CacheKeyWebBuybackSystems,
		keys.TypeStrWebBuybackSystems,
		WEB_BUYBACK_SYSTEMS_EXPIRES_IN,
		rep,
		nil,
	)
}
