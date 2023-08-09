package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/WiggidyW/weve-esi/logger"
	"github.com/redis/go-redis/v9"
)

type SharedServerCache *redis.Client

type SharedClientCache *fastcache.Cache

type Cache[D any, ED Expirable[D]] struct { // unique per type
	localCache  *LocalCache  // unique per type
	serverCache *ServerCache // shared (1)
	bufPool     *BufferPool  // unique per type
	// ttl for server lock
	sLockTTL time.Duration // unique per type
	// max wait time for server lock acquire (if > ttl, it has no effect))
	sLockMaxWait time.Duration // unique per type
}

func NewCache[D any, ED Expirable[D]](
	bufPool *BufferPool,
	cCache SharedClientCache,
	sCache SharedServerCache,
	sLockTTL time.Duration,
	sLockMaxWait time.Duration,
) *Cache[D, ED] {
	return &Cache[D, ED]{
		localCache:   newLocalCache(cCache),
		serverCache:  newServerCache(sCache),
		bufPool:      bufPool,
		sLockTTL:     sLockTTL,
		sLockMaxWait: sLockMaxWait,
	}
}

func (c *Cache[D, ED]) localCacheGet(key string) (*ED, error) {
	// get a buf from the pool
	buf := c.bufPool.Get()
	defer c.bufPool.Put(buf)

	// read from local cache
	data := c.localCache.get(key, *buf)
	if data == nil {
		return nil, nil
	}

	// deserialize
	val, err := deserialize[ED](data)
	if err != nil {
		return nil, err
	}

	// check expiration
	if (*val).Expires().Before(time.Now()) {
		c.localCache.del(key)
		return nil, nil
	}

	return val, nil
}

// expired results from serverCache won't happen because of TTL
func (c *Cache[D, ED]) serverCacheGet(
	ctx context.Context,
	key string,
) (*ED, error) {
	// read from server cache
	data, err := c.serverCache.get(ctx, key)
	if err != nil {
		return nil, err
	} else if data == nil {
		return nil, nil
	}

	// deserialize
	val, err := deserialize[ED](data)
	if err != nil {
		return nil, err
	}

	// insert into local cache
	c.localCache.set(key, data)

	return val, nil
}

func (c *Cache[D, ED]) GetOrLock(
	ctx context.Context,
	key string,
) (*ED, *Lock, error) {
	// create empty consolidated lock
	cLock := new(Lock)

	// lock local cache
	cLock.localLock = c.localCache.lock(key)

	// try to hit value from local cache
	if lcVal, err := c.localCacheGet(key); err != nil {
		cLock.localUnlock()
		return nil, nil, err
	} else if lcVal != nil {
		cLock.localUnlock()
		return lcVal, nil, nil
	}

	// lock server cache
	if serverLock, err := c.serverCache.lock(
		ctx,
		key,
		c.sLockTTL,
		c.sLockMaxWait,
	); err != nil {
		// if we fail to lock the server, log the error and continue
		logger.Logger.Error(err.Error())
		return nil, cLock, nil
	} else {
		cLock.serverLock = serverLock
	}

	// try to hit value from server cache
	if scVal, err := c.serverCacheGet(ctx, key); err != nil {
		// if we fail to get from server, log the error and continue
		logger.Logger.Error(err.Error())
		cLock.serverReleaseLogErr(ctx)
		return nil, cLock, nil
	} else if scVal != nil {
		return scVal, nil, nil
	}

	// return consolidated lock
	return nil, cLock, nil
}

func (c *Cache[D, ED]) Set(
	ctx context.Context,
	key string,
	val ED,
	lock *Lock,
) error {
	// get ttl from val.Expires()
	ttl := time.Until(val.Expires())
	if ttl < 0 {
		lock.releaseLogErr(ctx)
		return fmt.Errorf(
			"cache: cannot set expired value (key: %s, ttl: %s)",
			key,
			ttl,
		)
	}

	// get a buf from the pool
	buf := c.bufPool.Get()

	// serialize
	data, err := serialize[ED](val, buf)
	if err != nil {
		lock.releaseLogErr(ctx)
		c.bufPool.Put(buf)
		return err
	}

	// set local cache & release local lock
	c.localCache.set(key, data)
	lock.localUnlock()

	// set server cache in a goroutine
	go func() {
		err := c.serverCache.set(ctx, key, data, ttl)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
		c.bufPool.Put(buf)
		err = lock.serverRelease(ctx)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
	}()

	return nil
}

func deserialize[T any](data []byte) (*T, error) {
	// create an empty val
	var val T

	// create decoder
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	// decode
	err := decoder.Decode(&val)
	if err != nil {
		return nil, err
	}

	return &val, nil
}

func serialize[T any](
	val T,
	b *[]byte,
) ([]byte, error) {
	// create encoder
	buffer := bytes.NewBuffer(*b)
	encoder := gob.NewEncoder(buffer)

	// encode
	err := encoder.Encode(val)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
