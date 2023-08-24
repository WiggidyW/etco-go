package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/WiggidyW/eve-trading-co-go/logger"
)

// TODO: Add Local locks on top of server locks to StrongCache
//   so that we can avoid hitting the server for local contention

type StrongAntiCache struct {
	StrongCacheInner
}

func (sac *StrongAntiCache) Del(
	// ctx context.Context, // del should never be cancelled
	key string,
	lock *StrongLock,
) error {
	ctx := context.Background()

	// validate and refresh lock if needed
	if err := sac.handleLock(ctx, lock); err != nil {
		return err
		// delete from server cache
	} else if err := sac.serverCache.del(ctx, key); err != nil {
		return ErrServerDel{err}
		// operation successful
	} else {
		return nil
	}
}

// stores data on server side only that can be invalidated
// uses strong synchronization mechanisms and enforces consistency
// (should still be considered unreliable, failing is just less likely)
type StrongCache[D any, ED Expirable[D]] struct {
	StrongCacheInner
}

func (sc *StrongCache[D, ED]) Get(
	ctx context.Context,
	key string,
	lock *StrongLock,
) ( /*val*/ *ED /*err*/, error) {
	if err := sc.handleLock(ctx, lock); err != nil {
		return nil, err
	}

	// read from server cache
	b, err := sc.serverCache.get(ctx, key)
	if err != nil {
		return nil, ErrServerGet{err}
	} else if b == nil {
		return nil, nil
	}

	// deserialize
	val, err := deserialize[ED](b)
	if err != nil {
		return nil, ErrServerDeserialize{err}
	}

	// check expiration
	if (*val).Expires().Before(time.Now()) {
		// TTL is supposed to prevent this from happening
		logger.Warn(fmt.Errorf(
			"expired key: %s returned from server cache",
			key,
		))
		return nil, nil
	} else {
		return val, nil
	}
}

func (sc *StrongCache[D, ED]) Set(
	// ctx context.Context, // set should never be cancelled
	key string,
	val ED,
	lock *StrongLock,
) error {
	ctx := context.Background()

	if err := sc.handleLock(ctx, lock); err != nil {
		return err
	}

	// get TTL and check expired
	ttl := time.Until(val.Expires())
	if ttl < 0 {
		return ErrInvalidSet{key, ttl}
	}

	// get a buf from the pool
	buf := sc.bufPool.Get()
	defer sc.bufPool.Put(buf)

	// serialize
	b, err := serialize[ED](val, buf)
	if err != nil {
		return ErrSerialize{err}
	}

	// write to server cache
	if err := sc.serverCache.set(ctx, key, b, ttl); err != nil {
		return ErrServerSet{err}
	} else {
		return nil
	}
}

type StrongCacheInner struct {
	serverCache  *ServerCache
	bufPool      *BufferPool
	lLocks       *sync.Map     // prevents local contention from hitting server
	sLockTTL     time.Duration // should be pretty high
	sLockMaxWait time.Duration // should be pretty high
}

func (sci StrongCacheInner) Unlock(lock *StrongLock) error {
	if lock == nil || lock.expired() {
		return ErrInvalidLock{"StrongCacheUnlock"}
	} else {
		return lock.unlock()
	}
}

// converts the key to a lock key and then tries to obtain a lock from it
func (sci StrongCacheInner) Lock(
	ctx context.Context,
	key string,
) (*StrongLock, error) {
	return sci.lock(ctx, lockKey(key))
}

// tries to obtain a lock from the given key
func (sci StrongCacheInner) lock(
	ctx context.Context,
	lockKey string,
) (*StrongLock, error) {
	lLockAny, _ := sci.lLocks.LoadOrStore(lockKey, &sync.Mutex{})
	lLock := lLockAny.(*sync.Mutex)
	lLock.Lock()

	sLock, err := sci.serverCache.lock(
		ctx,
		lockKey,
		sci.sLockTTL,
		sci.sLockMaxWait,
	)
	if err != nil {
		lLock.Unlock()
		return nil, ErrServerLock{err}
	}

	return newStrongLock(sLock, lLock, sci.sLockTTL), nil
}

func (sci StrongCacheInner) handleLock(
	ctx context.Context,
	lock *StrongLock,
) error {
	if lock == nil {
		return ErrInvalidLock{"StrongCacheHandleLock"}
	} else if lock.expired() {
		if err := lock.refresh(ctx, sci.sLockTTL); err != nil {
			return err
		}
	}

	return nil
}
