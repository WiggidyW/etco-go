package cacheprefetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
)

type Params[REP any] struct {
	Get       *ActionGet[REP]      // 1st
	Namespace *ActionNamespace     // 2nd
	Lock      []ActionOrderedLocks // 3rd
}

type ActionGet[REP any] struct {
	CacheKey          keys.Key
	TypeStr           keys.Key
	Local             bool
	Server            bool
	NewRep            func() REP
	Slosh             cache.SetLocalOnServerHit[REP] // if server hit and local miss, set server value to local cache
	KeepLockAfterMiss bool                           // if true, keeps lock if cache miss
}

type ActionNamespace struct {
	CacheKey     keys.Key
	TypeStr      keys.Key
	ExpiredValid bool
}

type ActionOrderedLocks struct {
	Locks []ActionLock
	Child *ActionOrderedLocks
}

type ActionLock struct {
	CacheKey keys.Key
	TypeStr  keys.Key
	Local    bool
	Server   bool
}

// always: unlock all scoped locks if an error occurs
//
//  1. try to 'get' the value from cache. If it exists + not expired, return
//  2. check 'namespace'. If it was modified, then 'other' is refreshing the
//     cache and locking our 'get' lock.
//     Try again - re-lock the 'get' Lock. (should lock after 'other' finishes)
//  3. lock the 'lock' locks.
func Handle[REP any](
	x cache.Context,
	params Params[REP],
) (
	namespaceRetry bool,
	rep *expirable.Expirable[REP],
	err error,
) {
	defer func() {
		if err != nil {
			x.UnlockScoped()
		}
	}()

	// run the 'get' action if requested
	if params.Get != nil {
		rep, err = cache.GetOrLock(
			x,
			params.Get.CacheKey,
			params.Get.TypeStr,
			params.Get.Local,
			params.Get.Server,
			params.Get.NewRep,
			params.Get.Slosh,
		)
		if err != nil || rep != nil {
			return false, rep, err
		} else if !params.Get.KeepLockAfterMiss {
			go x.Unlock(params.Get.CacheKey, params.Get.TypeStr)
		}
	}

	// run the 'namespace' check if requested
	// (multi-cache-key lock / synchronization)
	if params.Namespace != nil {
		var cmd cache.NamespaceCommand
		var expires time.Time
		now := time.Now()
		cmd, expires, err = cache.NamespaceCheck(
			x,
			params.Namespace.CacheKey,
			params.Namespace.TypeStr,
			now,
			params.Namespace.ExpiredValid,
		)
		if err != nil {
			return false, nil, err
		} else if cmd == cache.NCRetry {
			return true, nil, nil
		} else if cmd == cache.NCRepEmpty {
			rep = &expirable.Expirable[REP]{
				Expires: expires,
				Data:    *new(REP),
			}
			return false, rep, nil
		}
	}

	// send out locks
	if params.Lock != nil {
		lenOrderedLocks := len(params.Lock)

		if lenOrderedLocks == 1 {
			err = handleCacheLocks(x, params.Lock[0])
			return false, nil, err
		}

		chnErr := make(chan error, lenOrderedLocks)

		for _, locksAction := range params.Lock {
			go func(locksAction ActionOrderedLocks) {
				chnErr <- handleCacheLocks(x, locksAction)
			}(locksAction)
		}

		for i := 0; i < lenOrderedLocks; i++ {
			err = <-chnErr
			if err != nil {
				return false, nil, err
			}
		}
	}

	return false, rep, err
}

func handleCacheLocks(
	x cache.Context,
	action ActionOrderedLocks,
) (
	err error,
) {
	lenLocks := len(action.Locks)

	if lenLocks == 1 {
		return cache.LockAndDel(
			x,
			action.Locks[0].CacheKey,
			action.Locks[0].TypeStr,
			action.Locks[0].Local,
			action.Locks[0].Server,
		)
	} else if lenLocks == 0 {
		return nil
	}

	chnErr := make(chan error, lenLocks)

	for _, lockAction := range action.Locks {
		go func(lockAction ActionLock) {
			chnErr <- cache.LockAndDel(
				x,
				lockAction.CacheKey,
				lockAction.TypeStr,
				lockAction.Local,
				lockAction.Server,
			)
		}(lockAction)
	}

	for i := 0; i < lenLocks; i++ {
		err = <-chnErr
		if err != nil {
			return err
		}
	}

	if action.Child == nil {
		return nil
	}

	return handleCacheLocks(x, *action.Child)
}
