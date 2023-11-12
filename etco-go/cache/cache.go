package cache

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"time"

	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/localcache"
	"github.com/WiggidyW/etco-go/cache/servercache"
)

type NamespaceCommand uint8

const (
	NCPanicIfErrNil NamespaceCommand = iota
	NCFetch
	NCRetry
	NCRepNil
)

func timesToBytes(
	modifyTime time.Time,
	expires time.Time,
) []byte {
	buf := make([]byte, binary.MaxVarintLen64*2)
	binary.PutVarint(buf, modifyTime.Unix())
	binary.PutVarint(buf[binary.MaxVarintLen64:], expires.Unix())
	return buf
}

func bytesToTimes(b []byte) (
	modifyTime time.Time,
	expires time.Time,
) {
	modifyTimeUnix, _ := binary.Varint(b[:binary.MaxVarintLen64])
	expiresUnix, _ := binary.Varint(b[binary.MaxVarintLen64:])
	modifyTime = time.Unix(modifyTimeUnix, 0)
	expires = time.Unix(expiresUnix, 0)
	return modifyTime, expires
}

// TODO: set local cache to nil for the respective key of the actual caller
func NamespaceCheck(
	x Context,
	nsKey, nsTypeStr string,
	startTime time.Time,
	// 'false' for invalidatable child-keys like remoteDB UserData
	// 'true' for non-invalidatable data like ESI data
	expiredValid bool,
) (
	cmd NamespaceCommand,
	expires time.Time,
	err error,
) {
	lock := x.GetLock(nsKey, nsTypeStr)

	err = x.LocalLock(lock)
	if err != nil {
		return cmd, expires, err
	}
	defer func() {
		if cmd != NCFetch || err != nil {
			go x.LocalUnlock(lock)
		}
	}()

	bLocal := localcache.Get(nsKey, make([]byte, 0, binary.MaxVarintLen64))
	if bLocal != nil {
		var lModTime time.Time
		lModTime, expires = bytesToTimes(bLocal)
		if lModTime.After(startTime) {
			cmd = NCRetry
			return cmd, expires, nil
		} else if expiredValid && expires.After(startTime) {
			cmd = NCRepNil
			return cmd, expires, nil
		}
	}

	err = x.ServerLock(lock)
	if err != nil {
		return cmd, expires, err
	}
	defer func() {
		if cmd != NCFetch || err != nil {
			go x.ServerUnlock(lock)
		}
	}()

	var bServer []byte
	bServer, err = servercache.Get(x.ctx, nsKey)
	if err != nil {
		return cmd, expires, err
	}

	if bServer != nil {
		var sModTime time.Time
		sModTime, expires = bytesToTimes(bServer)
		if sModTime.After(startTime) {
			cmd = NCRetry
			return cmd, expires, nil
		} else if expiredValid && expires.After(startTime) {
			cmd = NCRepNil
			return cmd, expires, nil
		}
	}

	cmd = NCFetch
	return cmd, expires, nil
}

func NamespaceModify(
	x Context,
	key, typeStr string,
	expires time.Time,
) (
	err error,
) {
	lock := x.GetLock(key, typeStr)
	b := timesToBytes(time.Now(), expires)
	localcache.Set(key, b)
	go x.LocalUnlock(lock)
	err = servercache.Set(x.ctx, key, b, time.Until(expires))
	go x.ServerUnlock(lock)
	return err
}

func LockAndDel(
	x Context,
	key, typeStr string,
	local, server bool,
) (
	err error,
) {
	lock := x.GetLock(key, typeStr)

	// always obtain local lock first
	err = x.LocalLock(lock)
	if err != nil {
		return err
	}

	// delete from local if requested
	if local {
		localcache.Del(key)
	}

	if !server {
		return nil
	}

	// lock and del server cache if requested
	err = x.ServerLock(lock)
	if err != nil {
		go x.LocalUnlock(lock)
		return err
	}

	err = servercache.Del(x.ctx, key)
	if err != nil {
		go x.LocalUnlock(lock)
		go x.ServerUnlock(lock)
	}

	return err
}

func SetAndUnlock(
	x Context,
	key, typeStr string,
	local, server bool,
	expirable any,
	expires time.Time,
) (
	err error,
) {
	lock := x.GetLock(key, typeStr)
	bufPool := BufPool(typeStr)

	var b []byte
	if time.Now().Before(expires) {
		// serialize (don't return error just yet)
		buf := bufPool.Get()
		defer bufPool.Put(buf)
		b, err = encode(buf, expirable)
	}

	// local cache set + unlock
	if local && err == nil && b != nil {
		localcache.Set(key, b)
	}
	go x.LocalUnlock(lock)

	// server cache set + unlock
	if server && err == nil && b != nil {
		err = servercache.Set(
			context.Background(), // never allow these to be cancelled
			key,
			b,
			time.Until(expires),
		)
	}
	go x.ServerUnlock(lock)

	return err
}

func GetOrLock[REP any](
	x Context,
	key, typeStr string,
	local, server bool,
	newRep func() *REP,
	slosh SetLocalOnServerHit[REP],
) (
	rep *expirable.Expirable[REP],
	err error,
) {
	lock := x.GetLock(key, typeStr)
	bufPool := BufPool(typeStr)

	// always obtain local lock first
	err = x.LocalLock(lock)
	if err != nil {
		return nil, err
	}

	// check local cache if requested
	if local {
		rep, err = localGet(key, newRep, bufPool)
		if !server || err != nil || rep != nil {
			go x.LocalUnlock(lock)
			return rep, err
		}
	}

	// lock and check server cache if requested
	err = x.ServerLock(lock)
	if err != nil {
		go x.LocalUnlock(lock)
		return nil, err
	}

	var repWithBytes *repWithBytes[REP]
	repWithBytes, err = serverGet(x.ctx, key, newRep)
	if err == nil && repWithBytes != nil {
		if local && slosh(repWithBytes.rep.Data) {
			localcache.Set(key, repWithBytes.bytes)
		}
		rep = repWithBytes.rep
	}

	// always unlock
	if err != nil || rep != nil {
		go x.LocalUnlock(lock)
		go x.ServerUnlock(lock)
	}

	return rep, err
}

type repWithBytes[REP any] struct {
	rep   *expirable.Expirable[REP]
	bytes []byte
}

func serverGet[REP any](
	ctx context.Context,
	key string,
	newRep func() *REP,
) (
	rwb *repWithBytes[REP],
	err error,
) {
	// get bytes from cache
	var b []byte
	b, err = servercache.Get(ctx, key)
	if err != nil || b == nil {
		return nil, err
	}

	// deserialize and check expired
	rwb = &repWithBytes[REP]{rep: nil, bytes: b}
	rwb.rep, err = decode[REP](b, newRep)
	if err != nil || !rwb.rep.Expired() { // unlock and return rep / error
		return rwb, err
	} else /* if rwb.rep.Expired() */ { // delete expired and return lock
		go servercache.DelLogErr(key)
		return nil, nil
	}
}

// (1) If err != nil, rep will be nil.
func localGet[REP any](
	key string,
	newRep func() *REP,
	BufPool *BufferPool,
) (
	rep *expirable.Expirable[REP],
	err error,
) {
	// obtain buf
	buf := BufPool.Get()
	defer BufPool.Put(buf)

	// get bytes from cache
	b := localcache.Get(key, *buf)
	if b == nil {
		return nil, nil
	}

	// deserialize and check expired
	rep, err = decode(b, newRep)
	if err == nil && rep.Expired() {
		localcache.Del(key)
		rep = nil
	}

	return rep, err
}

func decode[REP any](
	b []byte,
	newRep func() *REP,
) (
	rep *expirable.Expirable[REP],
	err error,
) {
	rep = initializeRep(newRep)
	reader := bytes.NewReader(b)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(rep)
	return rep, err
}

func encode(
	buf *[]byte,
	rep any,
) (
	b []byte,
	err error,
) {
	writer := bytes.NewBuffer(b)
	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(rep)
	if err == nil {
		b = writer.Bytes()
	}
	return b, err
}

func initializeRep[REP any](
	newRep func() *REP,
) *expirable.Expirable[REP] {
	var rep *REP
	if newRep != nil {
		rep = newRep()
	} else {
		rep = new(REP)
	}
	return expirable.NewMarshalPtr(rep)
}
