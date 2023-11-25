package servercache

import (
	"context"
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/logger"
	"github.com/bsm/redislock"
)

const (
	UNLOCK_RETRY_INTERVAL time.Duration = 10 * time.Millisecond
	MAX_UNLOCK_ATTEMPTS   int           = 10
)

type Lock struct {
	inner    *redislock.Lock
	ttl      time.Duration
	expires  time.Time
	released error
	mu       *sync.RWMutex
}

func newLock(
	ctx context.Context,
	l *redislock.Lock,
	ttl time.Duration,
) (lock *Lock) {
	lock = &Lock{
		inner:    l,
		ttl:      ttl,
		expires:  time.Now().Add(ttl),
		released: nil,
		mu:       new(sync.RWMutex),
	}
	go lock.holdUntilCancelled(ctx)
	return lock
}

func (sl *Lock) IsNil() bool {
	return sl == nil
}

func (sl *Lock) Released() (err error) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	err = sl.released
	if err == nil && sl.expires.Before(time.Now()) {
		err = redislock.ErrLockNotHeld
		go sl.markReleased(err)
	}
	return err
}

func (sl *Lock) markReleased(reason error) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	if sl.released == nil {
		sl.released = reason
	}
}

func (sl *Lock) unlockInner(attempt int) (err error) {
	err = sl.inner.Release(context.Background())
	if err == nil || err == redislock.ErrLockNotHeld {
		return nil
	} else if attempt >= MAX_UNLOCK_ATTEMPTS {
		return ErrServerUnlock{err}
	}
	return sl.unlockInner(attempt + 1)
}

func (sl *Lock) holdUntilCancelled(ctx context.Context) (err error) {
	ticker := time.NewTicker(sl.ttl / 2)

	defer ticker.Stop()
	defer sl.markReleased(err)
	defer func() { go logger.MaybeErr(sl.unlockInner(1)) }()

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return err
		case <-ticker.C:
			err = sl.refresh(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (sl *Lock) refresh(ctx context.Context) (err error) {
	err = sl.inner.Refresh(ctx, sl.ttl, nil)
	if err != nil {
		err = ErrServerRefreshLock{err}
	} else {
		sl.mu.Lock()
		sl.expires = time.Now().Add(sl.ttl)
		sl.mu.Unlock()
	}
	return nil
}
