package prefetch

import (
	"github.com/WiggidyW/etco-go/cache"
)

func CacheNamespace(
	cacheKey, typeStr string,
	expiredValid bool,
) *CacheActionNamespace {
	return &CacheActionNamespace{
		CacheKey:     cacheKey,
		TypeStr:      typeStr,
		ExpiredValid: expiredValid,
	}
}

func DualCacheOrderedLocksOne(
	cacheKey, typeStr string,
) []CacheActionOrderedLocks {
	return []CacheActionOrderedLocks{
		cacheLockOne(cacheKey, typeStr, true, true),
	}
}
func ServerCacheOrderedLocksOne(
	cacheKey, typeStr string,
) []CacheActionOrderedLocks {
	return []CacheActionOrderedLocks{
		cacheLockOne(cacheKey, typeStr, false, true),
	}
}
func LocalCacheOrderedLocksOne(
	cacheKey, typeStr string,
) []CacheActionOrderedLocks {
	return []CacheActionOrderedLocks{
		cacheLockOne(cacheKey, typeStr, true, false),
	}
}

func CacheOrderedLocks(
	child *CacheActionOrderedLocks,
	locks ...CacheActionLock,
) CacheActionOrderedLocks {
	return CacheActionOrderedLocks{Locks: locks, Child: child}
}
func CacheOrderedLocksPtr(
	child *CacheActionOrderedLocks,
	locks ...CacheActionLock,
) *CacheActionOrderedLocks {
	return &CacheActionOrderedLocks{Locks: locks, Child: child}
}

func CacheOrderedLocksNoFamily(
	child *CacheActionOrderedLocks,
	locks ...CacheActionLock,
) []CacheActionOrderedLocks {
	return []CacheActionOrderedLocks{{Locks: locks, Child: nil}}
}

func cacheLockOne(
	cacheKey, typeStr string,
	local, server bool,
) CacheActionOrderedLocks {
	return CacheActionOrderedLocks{
		Locks: []CacheActionLock{cacheLock(cacheKey, typeStr, local, server)},
		Child: nil,
	}
}

func DualCacheLockPtr(cacheKey, typeStr string) *CacheActionLock {
	l := DualCacheLock(cacheKey, typeStr)
	return &l
}
func ServerCacheLockPtr(cacheKey, typeStr string) *CacheActionLock {
	l := ServerCacheLock(cacheKey, typeStr)
	return &l
}
func LocalCacheLockPtr(cacheKey, typeStr string) *CacheActionLock {
	l := LocalCacheLock(cacheKey, typeStr)
	return &l
}

func DualCacheLock(cacheKey, typeStr string) CacheActionLock {
	return cacheLock(cacheKey, typeStr, true, true)
}
func ServerCacheLock(cacheKey, typeStr string) CacheActionLock {
	return cacheLock(cacheKey, typeStr, false, true)
}
func LocalCacheLock(cacheKey, typeStr string) CacheActionLock {
	return cacheLock(cacheKey, typeStr, true, false)
}

func cacheLock(cacheKey, typeStr string, local, server bool) CacheActionLock {
	return CacheActionLock{
		CacheKey: cacheKey,
		TypeStr:  typeStr,
		Local:    local,
		Server:   server,
	}
}

func DualCacheGet[REP any](
	cacheKey, typeStr string,
	keepLockAfterMiss bool,
	newRep func() REP,
	slosh cache.SetLocalOnServerHit[REP],
) *CacheActionGet[REP] {
	return cacheGet(
		cacheKey, typeStr,
		keepLockAfterMiss,
		newRep,
		slosh,
		true, true,
	)
}
func ServerCacheGet[REP any](
	cacheKey, typeStr string,
	keepLockAfterMiss bool,
	newRep func() REP,
) *CacheActionGet[REP] {
	return cacheGet(
		cacheKey, typeStr,
		keepLockAfterMiss,
		newRep,
		cache.SloshFalse[REP],
		false, true,
	)
}
func LocalCacheGet[REP any](
	cacheKey, typeStr string,
	keepLockAfterMiss bool,
	newRep func() REP,
) *CacheActionGet[REP] {
	return cacheGet(
		cacheKey, typeStr,
		keepLockAfterMiss,
		newRep,
		cache.SloshFalse[REP],
		true, false,
	)
}

func cacheGet[REP any](
	cacheKey, typeStr string,
	keepLockAfterMiss bool,
	newRep func() REP,
	slosh cache.SetLocalOnServerHit[REP],
	local, server bool,
) *CacheActionGet[REP] {
	return &CacheActionGet[REP]{
		CacheKey:          cacheKey,
		TypeStr:           typeStr,
		NewRep:            newRep,
		Slosh:             slosh,
		KeepLockAfterMiss: keepLockAfterMiss,
		Local:             local,
		Server:            server,
	}
}
