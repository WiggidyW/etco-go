package cache

import (
	"context"
	"time"

	"github.com/bsm/redislock"
)

type StrongLock struct {
	lock    *redislock.Lock
	expires time.Time
}

func newStrongLock(
	lock *redislock.Lock,
	ttl time.Duration,
) *StrongLock {
	return &StrongLock{lock, time.Now().Add(ttl)}
}

func (sl *StrongLock) expired() bool {
	return time.Now().After(sl.expires)
}

func (sl *StrongLock) refresh(
	ctx context.Context,
	ttl time.Duration,
) error {
	if err := sl.lock.Refresh(ctx, ttl, nil); err != nil {
		return ErrServerLock{err}
	} else {
		sl.expires = time.Now().Add(ttl)
	}
	return nil
}

func (sl *StrongLock) unlock() error {
	if err := sl.lock.Release(context.Background()); err != nil {
		return ErrServerUnlock{err}
	} else {
		return nil
	}
}
