package prefetch

import (
	"context"
	"time"

	"github.com/WiggidyW/chanresult"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
)

type CacheParams[REP any] struct {
	// Set is ignored if get is present.
	// ^ wish I had discriminated unions
	Get *CacheActionGet[REP]
	Set *CacheAction
	Del *[]CacheAction
}

type CacheAction struct {
	TypeStr string
	Key     string
	Local   bool
	Server  *cache.ServerLockParams
}

type CacheActionGet[REP any] struct {
	NewRep *func() *REP
	Slosh  *cache.SetLocalOnServerHit[REP]
	CacheAction
}

type CacheActionServer struct {
	LockTTL        time.Duration
	LockMaxBackoff time.Duration
}

type CacheLocks struct {
	Del *[]*cache.Lock
	Set *cache.Lock
}

func handleCache[REP any](
	ctx context.Context,
	params CacheParams[REP],
) (
	rep *expirable.Expirable[REP],
	locks *CacheLocks,
	err error,
) {
	locks = &CacheLocks{}

	// run the 'get' action if requested
	if params.Get != nil {
		rep, locks.Set, err = cache.GetOrLock(
			ctx,
			params.Get.Key,
			params.Get.TypeStr,
			params.Get.Local,
			params.Get.Server,
			params.Get.NewRep,
			params.Get.Slosh,
		)
		if params.Del == nil || err != nil || rep != nil {
			// only one of these will ever not be nil
			// wish I had discriminated unions
			return rep, locks, err
		}
	}

	// send out deletes
	if params.Del != nil {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		actions := *params.Del
		chn := chanresult.NewChanResult[*[]*cache.Lock](ctx, len(actions), 0)
		go handleCacheDels(ctx, actions, chn)
		defer func() {
			if err == nil {
				// if all subsequent operations succeed, receive the locks
				locks.Del, err = chn.Recv()
				if err != nil {
					// if that fails, unlock the set lock if it's present
					// and set locks to nil
					if locks.Set != nil {
						go locks.Set.UnlockLogErr()
					}
					locks = nil
				}
			}
		}()
	}

	// run the 'set' action if requested
	if params.Get == nil && params.Set != nil {
		locks.Set, err = cache.LockAndDel(
			ctx,
			params.Set.Key,
			params.Set.TypeStr,
			params.Set.Local,
			params.Set.Server,
		)
		if err != nil {
			return nil, nil, err
		}
	}

	return nil, locks, err
}

func handleCacheDels(
	ctx context.Context,
	actions []CacheAction,
	chn chanresult.ChanResult[*[]*cache.Lock],
) (
	err error,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// unbuffered so we can cancel if we need to
	chnDels := chanresult.NewChanResult[*cache.Lock](ctx, 0, 0)

	// send out the dels, making them wait for cancellation or recv
	// if they are cancelled, they will unlock
	for _, action := range actions {
		go func(action CacheAction) error {
			var lock *cache.Lock
			var err error
			lock, err = cache.LockAndDel(
				ctx,
				action.Key,
				action.TypeStr,
				action.Local,
				action.Server,
			)
			if err != nil {
				return chnDels.SendErr(err)
			}
			err = chnDels.SendOk(lock)
			if err != nil {
				go lock.UnlockLogErr()
			}
			return err
		}(action)
	}

	// receive the dels
	// if there is an error, cancel and break
	locks := make([]*cache.Lock, 0, len(actions))
	for i := 0; i < len(actions); i++ {
		var lock *cache.Lock
		lock, err = chnDels.Recv()
		if err != nil {
			cancel()
			break
		}
		locks = append(locks, lock)
	}

	// if checks out, send the locks
	if err == nil {
		err = chn.SendOk(&locks)
	}

	// if cancelled, or error, unlock all the locks
	if err != nil {
		for _, lock := range locks {
			go lock.UnlockLogErr()
		}
		err = chn.SendErr(err)
	}

	return err
}
