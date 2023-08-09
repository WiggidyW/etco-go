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

func (l *Lock) localUnlock() {
	if l.localLock != nil {
		l.localLock.Unlock()
		l.localLock = nil
	}
}

func (l *Lock) serverRelease(ctx context.Context) error {
	if l.serverLock != nil {
		err := l.serverLock.Release(ctx)
		if err != nil && err != redislock.ErrLockNotHeld {
			return fmt.Errorf(
				"error releasing lock: %w",
				err,
			)
		}
	}
	return nil
}

func (l *Lock) serverReleaseLogErr(ctx context.Context) {
	go func() {
		if err := l.serverRelease(ctx); err != nil {
			logger.Logger.Error(err.Error())
		}
	}()
}

func (l *Lock) releaseLogErr(ctx context.Context) {
	l.localUnlock()
	l.serverReleaseLogErr(ctx)
}
