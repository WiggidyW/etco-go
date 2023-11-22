package cacheprefetch

import (
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
)

func AntiCache(lock []ActionOrderedLocks) Params[struct{}] {
	return Params[struct{}]{Get: nil, Namespace: nil, Lock: lock}
}

func LocalLock(cacheKey, typeStr keys.Key) ActionLock {
	return ActionLock{
		CacheKey: cacheKey,
		TypeStr:  typeStr,
		Local:    true,
		Server:   false,
	}
}

func ServerLock(cacheKey, typeStr keys.Key) ActionLock {
	return ActionLock{
		CacheKey: cacheKey,
		TypeStr:  typeStr,
		Local:    false,
		Server:   true,
	}
}

func DualLock(cacheKey, typeStr keys.Key) ActionLock {
	return ActionLock{
		CacheKey: cacheKey,
		TypeStr:  typeStr,
		Local:    true,
		Server:   true,
	}
}

func TransientCache[REP any](
	cacheKey, typeStr keys.Key,
	newRep func() REP, // nil okay
	lock []ActionOrderedLocks,
) Params[REP] {
	return Params[REP]{
		Get: &ActionGet[REP]{
			CacheKey:          cacheKey,
			TypeStr:           typeStr,
			Local:             true,
			Server:            false,
			NewRep:            newRep,
			Slosh:             cache.SloshFalse[REP],
			KeepLockAfterMiss: true,
		},
		Namespace: nil,
		Lock:      lock,
	}
}

func StrongCache[REP any](
	cacheKey, typeStr keys.Key,
	newRep func() REP, // nil okay
	lock []ActionOrderedLocks,
) Params[REP] {
	return Params[REP]{
		Get: &ActionGet[REP]{
			CacheKey:          cacheKey,
			TypeStr:           typeStr,
			Local:             false,
			Server:            true,
			NewRep:            newRep,
			Slosh:             cache.SloshFalse[REP],
			KeepLockAfterMiss: true,
		},
		Namespace: nil,
		Lock:      lock,
	}
}

func WeakCache[REP any](
	cacheKey, typeStr keys.Key,
	newRep func() REP, // nil okay
	slosh cache.SetLocalOnServerHit[REP], // required
	lock []ActionOrderedLocks,
) Params[REP] {
	return Params[REP]{
		Get: &ActionGet[REP]{
			CacheKey:          cacheKey,
			TypeStr:           typeStr,
			Local:             true,
			Server:            true,
			NewRep:            newRep,
			Slosh:             slosh,
			KeepLockAfterMiss: true,
		},
		Namespace: nil,
		Lock:      lock,
	}
}

func StrongMultiCacheKnownKeys[REP any](
	cacheKey, typeStr keys.Key,
	nsCacheKey, nsTypeStr keys.Key,
	newRep func() REP, // nil okay
	lock []ActionOrderedLocks,
) Params[REP] {
	return Params[REP]{
		Get: &ActionGet[REP]{
			CacheKey:          cacheKey,
			TypeStr:           typeStr,
			Local:             false,
			Server:            true,
			NewRep:            newRep,
			Slosh:             cache.SloshFalse[REP],
			KeepLockAfterMiss: false,
		},
		Namespace: &ActionNamespace{
			CacheKey:     nsCacheKey,
			TypeStr:      nsTypeStr,
			ExpiredValid: false,
		},
		Lock: lock,
	}
}

func WeakMultiCacheKnownKeys[REP any](
	cacheKey, typeStr keys.Key,
	nsCacheKey, nsTypeStr keys.Key,
	newRep func() REP, // nil okay
	slosh cache.SetLocalOnServerHit[REP], // required
	lock []ActionOrderedLocks,
) Params[REP] {
	return Params[REP]{
		Get: &ActionGet[REP]{
			CacheKey:          cacheKey,
			TypeStr:           typeStr,
			Local:             true,
			Server:            true,
			NewRep:            newRep,
			Slosh:             slosh,
			KeepLockAfterMiss: false,
		},
		Namespace: &ActionNamespace{
			CacheKey:     nsCacheKey,
			TypeStr:      nsTypeStr,
			ExpiredValid: true,
		},
		Lock: lock,
	}
}

func WeakMultiCacheDynamicKeys[REP any](
	cacheKey, typeStr keys.Key,
	nsCacheKey, nsTypeStr keys.Key,
	newRep func() REP, // nil okay
	slosh cache.SetLocalOnServerHit[REP], // required
	lock []ActionOrderedLocks,
) Params[REP] {
	return Params[REP]{
		Get: &ActionGet[REP]{
			CacheKey:          cacheKey,
			TypeStr:           typeStr,
			Local:             true,
			Server:            true,
			NewRep:            newRep,
			Slosh:             slosh,
			KeepLockAfterMiss: true,
		},
		Namespace: &ActionNamespace{
			CacheKey:     nsCacheKey,
			TypeStr:      nsTypeStr,
			ExpiredValid: true,
		},
		Lock: lock,
	}
}
