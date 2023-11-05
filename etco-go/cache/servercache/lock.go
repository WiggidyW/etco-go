package servercache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/logger"
	"github.com/bsm/redislock"
)

const INCREMENTAL_RETRY_INTERVAL = 10 * time.Millisecond

type Lock struct {
	inner         *redislock.Lock
	ttl           time.Duration
	expires       time.Time
	released      bool
	mu            *sync.RWMutex
	cancelRefresh context.CancelFunc
}

func newLock(
	l *redislock.Lock,
	ttl time.Duration,
	expires time.Time,
) (lock *Lock) {
	ctx, cancel := context.WithCancel(context.Background())
	lock = &Lock{
		inner:         l,
		ttl:           ttl,
		expires:       expires,
		released:      false,
		mu:            &sync.RWMutex{},
		cancelRefresh: cancel,
	}
	go lock.refreshUntilCancelled(ctx)
	return lock
}

func (sl *Lock) Expired() bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.released || time.Now().After(sl.expires)
}

func (sl *Lock) Unlock() (err error) {
	sl.cancelRefresh()

	sl.mu.Lock()
	defer sl.mu.Unlock()

	if sl.released {
		return nil
	}

	err = sl.inner.Release(context.Background())
	if err != nil {
		if err == redislock.ErrLockNotHeld {
			sl.released = true
		}
		err = ErrServerUnlock{fmt.Errorf("%s: %w", sl.inner.Key(), err)}
	} else {
		sl.released = true
	}

	return err
}

func (sl *Lock) UnlockLogErr() {
	err := sl.Unlock()
	if err != nil {
		logger.Err(err.Error())
	}
}

func (sl *Lock) refreshUntilCancelled(ctx context.Context) {
	ticker := time.NewTicker(sl.ttl / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := sl.refresh(ctx)
			if err != nil {
				if !sl.released {
					logger.Err(err.Error())
				}
				return
			}
		}
	}
}

func (sl *Lock) refresh(ctx context.Context) (err error) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	if sl.released {
		err = redislock.ErrLockNotHeld
		return ErrServerRefreshLock{
			Released: true,
			err:      fmt.Errorf("%s: %w", sl.inner.Key(), err),
		}
	}

	err = sl.inner.Refresh(ctx, sl.ttl, nil)
	if err != nil {
		return ErrServerRefreshLock{
			Released: false,
			err:      fmt.Errorf("%s: %w", sl.inner.Key(), err),
		}
	} else {
		sl.expires = time.Now().Add(sl.ttl)
		sl.released = false
	}
	return nil
}
