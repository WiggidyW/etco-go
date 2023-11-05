package prefetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
)

func DualCacheSet(
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
) *CacheAction {
	return dualCacheAction(typeStr, cacheKey, lockTTL, lockMaxBackoff)
}

func ServerCacheSet(
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
) *CacheAction {
	return serverCacheAction(typeStr, cacheKey, lockTTL, lockMaxBackoff)
}

func LocalCacheSet(
	typeStr, cacheKey string,
) *CacheAction {
	return localCacheAction(typeStr, cacheKey)
}

func DualCacheDel(
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
) CacheAction {
	return *dualCacheAction(typeStr, cacheKey, lockTTL, lockMaxBackoff)
}

func ServerCacheDel(
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
) CacheAction {
	return *serverCacheAction(typeStr, cacheKey, lockTTL, lockMaxBackoff)
}

func LocalCacheDel(
	typeStr, cacheKey string,
) CacheAction {
	return *localCacheAction(typeStr, cacheKey)
}

func DualCacheGet[REP any](
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
	newRep *func() *REP,
	slosh *cache.SetLocalOnServerHit[REP],
) *CacheActionGet[REP] {
	return &CacheActionGet[REP]{
		NewRep: newRep,
		CacheAction: *dualCacheAction(
			typeStr, cacheKey,
			lockTTL, lockMaxBackoff,
		),
		Slosh: slosh,
	}
}

func ServerCacheGet[REP any](
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
	newRep *func() *REP,
) *CacheActionGet[REP] {
	return &CacheActionGet[REP]{
		NewRep: newRep,
		CacheAction: *serverCacheAction(
			typeStr, cacheKey,
			lockTTL, lockMaxBackoff,
		),
	}
}

func LocalCacheGet[REP any](
	typeStr, cacheKey string,
	newRep *func() *REP,
) *CacheActionGet[REP] {
	return &CacheActionGet[REP]{
		NewRep:      newRep,
		CacheAction: *localCacheAction(typeStr, cacheKey),
	}
}

func dualCacheAction(
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
) *CacheAction {
	return &CacheAction{
		TypeStr: typeStr,
		Key:     cacheKey,
		Local:   true,
		Server: &cache.ServerLockParams{
			TTL:        lockTTL,
			MaxBackoff: lockMaxBackoff,
		},
	}
}

func serverCacheAction(
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff time.Duration,
) *CacheAction {
	return &CacheAction{
		TypeStr: typeStr,
		Key:     cacheKey,
		Local:   false,
		Server: &cache.ServerLockParams{
			TTL:        lockTTL,
			MaxBackoff: lockMaxBackoff,
		},
	}
}

func localCacheAction(
	typeStr, cacheKey string,
) *CacheAction {
	return &CacheAction{
		TypeStr: typeStr,
		Key:     cacheKey,
		Local:   true,
		Server:  nil,
	}
}
