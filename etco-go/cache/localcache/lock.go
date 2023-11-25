package localcache

import (
	"context"
	"sync"
)

type Lock struct {
	inner    *sync.Mutex
	released error
	mu       *sync.RWMutex
}

func newLock(
	ctx context.Context,
	l *sync.Mutex,
) (lock *Lock) {
	lock = &Lock{
		inner:    l,
		released: nil,
		mu:       new(sync.RWMutex),
	}
	go lock.holdUntilCancelled(ctx)
	return lock
}

func (ll *Lock) IsNil() bool {
	return ll == nil
}

func (ll *Lock) Released() (err error) {
	ll.mu.RLock()
	err = ll.released
	ll.mu.RUnlock()
	return err
}

func (ll *Lock) holdUntilCancelled(ctx context.Context) (err error) {
	<-ctx.Done()
	err = ctx.Err()
	ll.mu.Lock()
	ll.inner.Unlock()
	ll.released = err
	ll.mu.Unlock()
	return err
}
