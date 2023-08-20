package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/bsm/redislock"
)

type WeakLock struct {
	localLock  *sync.Mutex
	serverLock *redislock.Lock
}

func (wl *WeakLock) localUnlock() error {
	if wl.localLock == nil {
		return ErrLocalUnlock{
			fmt.Errorf("local lock already unlocked"),
		}
	}

	wl.localLock.Unlock()
	wl.localLock = nil
	return nil
}

// create a new context rather than using a parameter, because we never want to cancel this
func (wl *WeakLock) serverUnlock() error {
	if wl.serverLock == nil {
		return ErrServerUnlock{
			fmt.Errorf("server lock already unlocked"),
		}
	}

	if err := wl.serverLock.Release(
		context.Background(),
	); err != nil && err != redislock.ErrLockNotHeld {
		return ErrServerUnlock{err}
	} else {
		wl.serverLock = nil
		return nil
	}
}
