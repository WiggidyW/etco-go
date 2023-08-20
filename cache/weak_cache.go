package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/WiggidyW/weve-esi/logger"
)

// stores data that will never be invalidated
// not trustworthy, but fast
type WeakCache[D any, ED Expirable[D]] struct { // unique per type
	localCache  *LocalCache  // unique per type
	serverCache *ServerCache // shared (1)
	bufPool     *BufferPool  // unique per type
	// ttl for server lock
	sLockTTL time.Duration // unique per type
	// max wait time for server lock acquire (if > ttl, it has no effect))
	sLockMaxWait time.Duration // unique per type
}

func NewWeakCache[D any, ED Expirable[D]](
	bufPool *BufferPool,
	cCache SharedClientCache,
	sCache SharedServerCache,
	sLockTTL time.Duration,
	sLockMaxWait time.Duration,
) *WeakCache[D, ED] {
	return &WeakCache[D, ED]{
		localCache:   newLocalCache(cCache),
		serverCache:  newServerCache(sCache),
		bufPool:      bufPool,
		sLockTTL:     sLockTTL,
		sLockMaxWait: sLockMaxWait,
	}
}

// logs errors and returns nil if deserialization fails
func (wc *WeakCache[D, ED]) localCacheGet(key string) *ED { // should lock before calling
	// get a buf from the pool
	buf := wc.bufPool.Get()
	defer wc.bufPool.Put(buf)

	// read from local cache
	data := wc.localCache.get(key, *buf)
	if data == nil {
		return nil
	}

	// deserialize
	val, err := deserialize[ED](data)
	if err != nil {
		logger.Err(ErrLocalDeserialize{err})
		return nil
	}

	// check expiration
	if (*val).Expires().Before(time.Now()) {
		// delete expired key
		wc.localCache.del(key)
		return nil
	} else {
		return val
	}
}

// logs errors and returns nil if deserialization fails
// inserts into local cache if server cache contains value
func (wc *WeakCache[D, ED]) serverCacheGet( // should lock before calling
	ctx context.Context,
	key string,
) *ED {
	// read from server cache
	data, err := wc.serverCache.get(ctx, key)
	if err != nil {
		logger.Err(ErrServerGet{err})
		return nil
	} else if data == nil {
		return nil
	}

	// deserialize
	val, err := deserialize[ED](data)
	if err != nil {
		logger.Err(ErrServerDeserialize{err})
		return nil
	} else if val == nil {
		return nil
	}

	// check expiration
	if (*val).Expires().Before(time.Now()) {
		// TTL is supposed to prevent this from happening
		logger.Warn(fmt.Errorf(
			"expired key: %s returned from server cache",
			key,
		))
		return nil
	} else {
		// insert into local cache and return
		wc.localCache.set(key, data)
		return val
	}
}

func (wc *WeakCache[D, ED]) Lock(ctx context.Context, key string) *WeakLock {
	lockKey := lockKey(key)
	cLock := new(WeakLock)

	// lock local cache
	cLock.localLock = wc.localCache.lock(lockKey)

	// lock server cache
	if serverLock, err := wc.serverCache.lock(
		ctx,
		lockKey,
		wc.sLockTTL,
		wc.sLockMaxWait,
	); err != nil {
		// return the partial lock
		logger.Err(ErrServerLock{err})
	} else {
		cLock.serverLock = serverLock
	}

	return cLock
}

// should be called with 'go'
func (wc *WeakCache[D, ED]) Unlock(lock *WeakLock) error {
	// unlock local cache
	if err := lock.localUnlock(); err != nil {
		return err
	}

	// unlock server cache
	if err := lock.serverUnlock(); err != nil {
		return err
	}

	return nil
}

// always returns either a lock or a value
// a partial lock may be returned (local cache only) if the server fails or is locked
// if there's a cache hit, this will release its locks in the background
func (wc *WeakCache[D, ED]) GetOrLock(
	ctx context.Context,
	key string,
) (*ED, *WeakLock) {
	lockKey := lockKey(key)
	cLock := new(WeakLock)

	// lock local cache
	cLock.localLock = wc.localCache.lock(lockKey)

	// try to hit value from local cache
	if lcVal := wc.localCacheGet(key); lcVal != nil {
		// local cache hit, unlock and return value
		logger.Fatal(cLock.localUnlock())
		return lcVal, nil
	}

	// lock server cache
	if serverLock, err := wc.serverCache.lock(
		ctx,
		lockKey,
		wc.sLockTTL,
		wc.sLockMaxWait,
	); err != nil {
		// return the partial lock
		logger.Err(ErrServerLock{err})
		return nil, cLock
	} else {
		cLock.serverLock = serverLock
	}

	// try to hit value from server cache
	if scVal := wc.serverCacheGet(ctx, key); scVal != nil {
		// server cache hit, unlock and return value
		logger.Fatal(cLock.localUnlock())
		go func() { logger.Err(cLock.serverUnlock()) }()
		return scVal, nil
	} else {
		// server cache miss, return the lock
		return nil, cLock
	}
}

// should be called with 'go'
func (wc *WeakCache[D, ED]) Set(
	key string,
	val ED,
	cLock *WeakLock,
) error {
	if cLock == nil || cLock.localLock == nil {
		return ErrInvalidLock{"WeakCacheSet"}
	}

	// get TTL and check expired
	ttl := time.Until(val.Expires())
	if ttl < 0 {
		logger.Err(cLock.localUnlock())
		if cLock.serverLock != nil {
			logger.Err(cLock.serverUnlock())
		}
		return ErrInvalidSet{key, ttl}
	}

	// get a buf from the pool
	buf := wc.bufPool.Get()

	// serialize
	data, err := serialize[ED](val, buf)
	if err != nil {
		wc.bufPool.Put(buf) // no longer needed
		logger.Err(cLock.localUnlock())
		if cLock.serverLock != nil {
			logger.Err(cLock.serverUnlock())
		}
		return ErrSerialize{err}
	}

	// set local cache & release local lock
	wc.localCache.set(key, data)
	if err := cLock.localUnlock(); err != nil {
		return err
	}

	// set server cache in a goroutine (if serverLock isn't nil)
	if cLock.serverLock != nil {
		err := wc.serverCache.set(
			context.Background(),
			key,
			data,
			ttl,
		)
		wc.bufPool.Put(buf) // no longer needed
		if err != nil {
			logger.Err(cLock.serverUnlock())
			return ErrServerSet{err}
		} else if err := cLock.serverUnlock(); err != nil {
			return err
		}
	} else {
		wc.bufPool.Put(buf) // no longer needed
	}

	return nil
}

// should be called with 'go'
func (wc *WeakCache[D, ED]) Del(
	ctx context.Context,
	key string,
	cLock *WeakLock,
) error {
	if cLock == nil || cLock.localLock == nil {
		return ErrInvalidLock{"WeakCacheSet"}
	}

	// del local cache & release local lock
	wc.localCache.del(key)
	logger.Fatal(cLock.localUnlock())

	// del server cache in a goroutine (if serverLock isn't nil)
	if cLock.serverLock != nil {
		if err := wc.serverCache.del(ctx, key); err != nil {
			logger.Err(cLock.serverUnlock())
			return ErrServerDel{err}
		}
		if err := cLock.serverUnlock(); err != nil {
			return err
		}
	}

	return nil
}
