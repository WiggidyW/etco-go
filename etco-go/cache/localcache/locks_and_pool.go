package localcache

import (
	"context"
	"sync"
)

type TypeLocksAndBufPools map[string]TypeLocksAndBufPool

func newTypeLocksAndBufPools() TypeLocksAndBufPools {
	return make(TypeLocksAndBufPools)
}

func (tlbps TypeLocksAndBufPools) register(
	typeStr string,
	bufPoolCap int,
) {
	tlbps[typeStr] = newTypeLocksAndBufPool(bufPoolCap)
}

func (tlbps TypeLocksAndBufPools) get(typeStr string) TypeLocksAndBufPool {
	return tlbps[typeStr]
}

type TypeLocksAndBufPool struct {
	locks   *sync.Map
	bufPool *BufferPool
}

func newTypeLocksAndBufPool(bufPoolCap int) TypeLocksAndBufPool {
	return TypeLocksAndBufPool{
		locks:   new(sync.Map),
		bufPool: newBufferPool(bufPoolCap),
	}
}

func (tlbp TypeLocksAndBufPool) ObtainLock(
	ctx context.Context,
	key string,
) (
	lock *sync.Mutex,
	err error,
) {
	lockAny, _ := tlbp.locks.LoadOrStore(key, &sync.Mutex{})
	lock = lockAny.(*sync.Mutex)
	lock.Lock()
	if ctx.Err() != nil {
		lock.Unlock()
		return nil, ctx.Err()
	} else {
		return lock, nil
	}
}

func (tlbp TypeLocksAndBufPool) BufGet() *[]byte {
	return tlbp.bufPool.get()
}

func (tlbp TypeLocksAndBufPool) BufPut(buf *[]byte) {
	tlbp.bufPool.put(buf)
}
