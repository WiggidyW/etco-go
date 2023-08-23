package cache

import (
	"context"
	"sync"
	"time"

	"github.com/bsm/redislock"
)

type StrongLock struct {
	sLock   *redislock.Lock
	lLock   *sync.Mutex
	expires time.Time
}

func newStrongLock(
	sLock *redislock.Lock,
	lLock *sync.Mutex,
	ttl time.Duration,
) *StrongLock {
	return &StrongLock{sLock, lLock, time.Now().Add(ttl)}
}

func (sl *StrongLock) expired() bool {
	return time.Now().After(sl.expires)
}

func (sl *StrongLock) refresh(
	ctx context.Context,
	ttl time.Duration,
) error {
	if err := sl.sLock.Refresh(ctx, ttl, nil); err != nil {
		return ErrServerLock{err}
	} else {
		sl.expires = time.Now().Add(ttl)
	}
	return nil
}

func (sl *StrongLock) unlock() error {
	sl.lLock.Unlock()
	if err := sl.sLock.Release(context.Background()); err != nil {
		return ErrServerUnlock{err}
	} else {
		return nil
	}
}
