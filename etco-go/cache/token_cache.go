package cache

import (
	"sync"
	"time"
)

type RawAuth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenCache struct {
	localCache *LocalCache
	bufPool    *BufferPool
}

func NewTokenCache(
	cCache SharedClientCache,
) *TokenCache {
	return &TokenCache{
		localCache: newLocalCache(cCache),
		bufPool:    NewBufferPool(0),
	}
}

func (tc *TokenCache) get(
	key string,
	lock *sync.Mutex,
) (*RawAuth, error) {
	// get a buf from the pool
	buf := tc.bufPool.Get()

	// get bytes from the cache
	data := tc.localCache.get(key, *buf)
	if data == nil {
		tc.bufPool.Put(buf)
		return nil, nil
	}

	// deserialize
	val, err := deserialize[ExpirableData[RawAuth]](data)
	tc.bufPool.Put(buf)
	if err != nil {
		return nil, ErrLocalDeserialize{err}
	}

	// check expired
	if (*val).Expires().Before(time.Now()) {
		return nil, nil
	}

	rep := val.Data()
	return &rep, nil
}

func (tc *TokenCache) GetOrLock(
	key string,
) (*RawAuth, *sync.Mutex, error) {
	lock := tc.localCache.lock(key)
	val, err := tc.get(key, lock)

	if err != nil {
		tc.Unlock(lock)
		return nil, nil, err
	} else if val != nil {
		tc.Unlock(lock)
		return val, nil, nil
	} else {
		return nil, lock, nil
	}
}

// unlocks the lock before returning
func (tc *TokenCache) Set(
	key string,
	val RawAuth,
	expires time.Time,
	lock *sync.Mutex,
) error {
	// get a buf from the pool
	buf := tc.bufPool.Get()

	// serialize
	data, err := serialize[ExpirableData[RawAuth]](
		NewExpirableData(val, expires),
		buf,
	)
	if err != nil {
		tc.bufPool.Put(buf)
		lock.Unlock()
		return ErrSerialize{err}
	}

	// write to local cache
	tc.localCache.set(key, data)
	tc.bufPool.Put(buf)
	lock.Unlock()

	return nil
}

func (tc *TokenCache) Unlock(lock *sync.Mutex) {
	lock.Unlock()
}
