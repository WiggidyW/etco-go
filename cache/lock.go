package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/bsm/redislock"

	"github.com/WiggidyW/weve-esi/logger"
)

type Lock struct {
	localLock  *sync.Mutex
	serverLock *redislock.Lock
}

func (l *Lock) validate() {
	if l.localLock == nil {
		logger.Logger.Fatal("localLock is nil")
	}
}

func (l *Lock) localUnlock() {
	if l.localLock != nil {
		l.localLock.Unlock()
		l.localLock = nil
	}
}

// create a new context rather than using a parameter, because we never want to cancel this
func (l *Lock) serverUnlock() error {
	if l.serverLock != nil {
		ctx := context.Background()
		err := l.serverLock.Release(ctx)
		if err != nil && err != redislock.ErrLockNotHeld {
			return fmt.Errorf(
				"error releasing lock: %w",
				err,
			)
		}
		l.serverLock = nil
	}
	return nil
}

func (l *Lock) serverUnlockLogErr() {
	go func() {
		if err := l.serverUnlock(); err != nil {
			logger.Logger.Error(err.Error())
		}
	}()
}

func (l *Lock) unlockLogErr() {
	l.localUnlock()
	l.serverUnlockLogErr()
}
