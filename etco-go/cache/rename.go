package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"math"
	"time"

	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/localcache"
	"github.com/WiggidyW/etco-go/cache/servercache"
)

type ServerLockParams struct {
	TTL        time.Duration
	MaxBackoff time.Duration
}

func SetAndUnlock[REP any](
	key string,
	typeStr string,
	local bool,
	server bool,
	lock *Lock,
	rep *REP,
	expires *time.Time,
) (
	err error,
) {
	locksAndBufPool := localcache.GetLocksAndBufPool(typeStr)
	buf := locksAndBufPool.BufGet()
	defer locksAndBufPool.BufPut(buf)

	// serialize (don't return error just yet)
	b, err := encode(buf, expirable.New(rep, expires))

	// local cache set + unlock
	if local && err == nil {
		localcache.Set(key, b)
	}
	if lock.local != nil {
		lock.local.Unlock()
	}

	// server cache set + unlock
	if server && err == nil {
		err = servercache.Set(
			context.Background(),
			key,
			b,
			serverSetTTL(expires),
		)
	}
	if lock.server != nil {
		go lock.server.UnlockLogErr()
	}

	return err
}

func LockAndDel(
	ctx context.Context,
	key string,
	typeStr string,
	local bool,
	server *ServerLockParams,
) (
	lock *Lock,
	err error,
) {
	lock = newLock(nil, nil)
	locksAndBufPool := localcache.GetLocksAndBufPool(typeStr)

	// always obtain local lock first
	lock.local, err = locksAndBufPool.ObtainLock(ctx, key)
	if err != nil {
		return nil, err
	}

	// delete from local if requested
	if local {
		localcache.Del(key)
	}

	// lock and del server cache if requested
	if server != nil {
		lock.server, err = serverLockAndDel(
			ctx,
			key,
			server.TTL,
			server.MaxBackoff,
		)
		if err != nil {
			lock.local.Unlock()
			return nil, err
		}
	}

	return lock, nil
}

func GetOrLock[REP any](
	ctx context.Context,
	key string,
	typeStr string,
	local bool,
	server *ServerLockParams,
	newRep *func() *REP,
	slosh *SetLocalOnServerHit[REP],
) (
	rep *expirable.Expirable[REP],
	lock *Lock,
	err error,
) {
	lock = newLock(nil, nil)
	locksAndBufPool := localcache.GetLocksAndBufPool(typeStr)

	// always obtain local lock first
	lock.local, err = locksAndBufPool.ObtainLock(ctx, key)
	if err != nil {
		return nil, nil, err
	}

	// check local cache if requested
	if local {
		rep, err = localGet(key, newRep, locksAndBufPool)
		if err != nil || rep != nil {
			lock.local.Unlock()
			// one of the two will always be nil (1)
			return rep, nil, err
		}
	}

	// check server cache if requested
	if server != nil {
		var repWithBytes *repWithBytes[REP]
		repWithBytes, lock.server, err = serverGetOrLock[REP](
			ctx,
			key,
			newRep,
			server.TTL,
			server.MaxBackoff,
		)
		if err != nil {
			lock.local.Unlock()
			return nil, nil, err
		} else if repWithBytes != nil {
			if local && setLocalOnServerHitOrDefault(slosh, true) {
				localcache.Set(key, repWithBytes.bytes)
			}
			lock.local.Unlock()
			return repWithBytes.rep, nil, nil
		}
	}

	return nil, lock, nil
}

type repWithBytes[REP any] struct {
	rep   *expirable.Expirable[REP]
	bytes []byte
}

func serverLockAndDel(
	ctx context.Context,
	key string,
	ttl time.Duration,
	maxBackoff time.Duration,
) (
	lock *servercache.Lock,
	err error,
) {
	lock, err = servercache.ObtainLock(ctx, key, ttl, maxBackoff)
	if err != nil {
		return nil, err
	}

	// delete from cache, blocking
	err = servercache.Del(ctx, key)
	if err != nil {
		// unlock and return error
		go lock.UnlockLogErr()
		return nil, err
	} else {
		return lock, nil
	}
}

func serverGetOrLock[REP any](
	ctx context.Context,
	key string,
	newRep *func() *REP,
	ttl time.Duration,
	maxBackoff time.Duration,
) (
	rwb *repWithBytes[REP],
	lock *servercache.Lock,
	err error,
) {
	// lock
	lock, err = servercache.ObtainLock(ctx, key, ttl, maxBackoff)
	if err != nil {
		return nil, nil, err
	}

	// get bytes from cache
	b, err := servercache.Get(ctx, key)
	if err != nil {
		// unlock and return error
		go lock.UnlockLogErr()
		return nil, nil, err
	} else if b == nil {
		// return lock
		return nil, lock, nil
	}

	// deserialize and check expired
	rwb = &repWithBytes[REP]{rep: nil, bytes: b}
	rwb.rep, err = decode[REP](b, newRep)
	if err != nil {
		// unlock and return error
		go lock.UnlockLogErr()
		return nil, nil, err
	} else if rwb.rep.Expired() {
		// delete expired and return lock
		go servercache.DelLogErr(key)
		return nil, lock, nil
	} else /* if !repWithBytes.rep.Expired() */ {
		// unlock and return rep
		go lock.UnlockLogErr()
		return rwb, nil, nil
	}
}

// (1) If err != nil, rep will be nil.
func localGet[REP any](
	key string,
	newRep *func() *REP,
	locksAndBufPool localcache.TypeLocksAndBufPool,
) (
	rep *expirable.Expirable[REP],
	err error,
) {
	// obtain buf
	buf := locksAndBufPool.BufGet()
	defer locksAndBufPool.BufPut(buf)

	// get bytes from cache
	b := localcache.Get(key, *buf)
	if b == nil {
		return nil, nil
	}

	// deserialize and check expired
	rep, err = decode(b, newRep)
	if err != nil {
		return nil, err
	} else if rep.Expired() {
		localcache.Del(key)
		return nil, nil
	} else {
		return rep, nil
	}
}

func decode[REP any](
	b []byte,
	newRep *func() *REP,
) (
	rep *expirable.Expirable[REP],
	err error,
) {
	rep = initializeRep(newRep)
	reader := bytes.NewReader(b)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(rep)
	if err != nil {
		return nil, err
	} else {
		return rep, nil
	}
}

func encode[REP any](
	buf *[]byte,
	rep expirable.Expirable[REP],
) (
	b []byte,
	err error,
) {
	writer := bytes.NewBuffer(b)
	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(rep)
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func initializeRep[REP any](
	newRep *func() *REP,
) *expirable.Expirable[REP] {
	var rep *REP
	if newRep != nil {
		rep = (*newRep)()
	} else {
		rep = new(REP)
	}
	return expirable.NewMarshalPtr(rep)
}

func serverSetTTL(expires *time.Time) time.Duration {
	if expires == nil {
		return math.MaxInt64
	} else {
		return time.Until(*expires)
	}
}
