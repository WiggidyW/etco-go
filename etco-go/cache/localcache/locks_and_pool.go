package localcache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
)

type typeLocksAndBufPools map[[16]byte]TypeLocksAndBufPool

func newTypeLocksAndBufPools() typeLocksAndBufPools {
	return make(typeLocksAndBufPools)
}

func (tlbps typeLocksAndBufPools) register(
	typeStr keys.Key,
	bufPoolCap int,
) {
	tlbps[typeStr.Bytes16()] = newTypeLocksAndBufPool(bufPoolCap)
}

func (tlbps typeLocksAndBufPools) get(typeStr keys.Key) TypeLocksAndBufPool {
	return tlbps[typeStr.Bytes16()]
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
	key keys.Key,
	maxWait time.Duration,
) (
	lock *Lock,
	err error,
) {
	lockAny, _ := tlbp.locks.LoadOrStore(key.Bytes16(), new(sync.Mutex))
	rawLock := lockAny.(*sync.Mutex)
	err = lockWithTimeout(ctx, rawLock, maxWait)
	if err == nil {
		lock = newLock(ctx, rawLock)
	}
	return lock, err
}

func lockWithTimeout(
	ctx context.Context,
	mu *sync.Mutex,
	maxWait time.Duration,
) (err error) {
	ticker := time.NewTicker(maxWait)
	defer ticker.Stop()
	chnDone := make(chan struct{}, 1)
	go func() {
		mu.Lock()
		chnDone <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-ticker.C:
		err = fmt.Errorf("local lock timed out after %ds", maxWait/1e9)
	case <-chnDone:
	}
	if err != nil {
		go func() {
			<-chnDone
			mu.Unlock()
		}()
	}
	return err
}
