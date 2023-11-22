package localcache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type typeLocksAndBufPools map[[16]byte]TypeLocksAndBufPool

func newTypeLocksAndBufPools() typeLocksAndBufPools {
	return make(typeLocksAndBufPools)
}

func (tlbps typeLocksAndBufPools) register(
	typeStr [16]byte,
	bufPoolCap int,
) {
	tlbps[typeStr] = newTypeLocksAndBufPool(bufPoolCap)
}

func (tlbps typeLocksAndBufPools) get(typeStr [16]byte) TypeLocksAndBufPool {
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

func (tlbp TypeLocksAndBufPool) obtainLock(
	ctx context.Context,
	key [16]byte,
	maxWait time.Duration,
) (
	lock *Lock,
	err error,
) {
	lockAny, _ := tlbp.locks.LoadOrStore(key, new(sync.Mutex))
	rawLock := lockAny.(*sync.Mutex)
	err = lockWithTimeout(rawLock, maxWait)
	if err == nil {
		lock = newLock(ctx, rawLock)
	}
	return lock, err
}

func lockWithTimeout(
	mu *sync.Mutex,
	maxWait time.Duration,
) (err error) {
	endTime := time.Now().Add(maxWait)
	mu.Lock()
	if time.Now().After(endTime) {
		mu.Unlock()
		err = fmt.Errorf("local lock timed out after %ds", maxWait/1e9)
	}
	return err
}
