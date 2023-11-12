package cache

import (
	"context"
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/cache/localcache"
	"github.com/WiggidyW/etco-go/cache/servercache"
)

const (
	LLOCK_MAX_WAIT    time.Duration = 4 * time.Second
	SLOCK_TTL         time.Duration = 4 * time.Second
	SLOCK_MAX_BACKOFF time.Duration = 4 * time.Second
)

// type LockId struct {
// 	scope   int64
// 	key     string
// 	typeStr string
// }

type LocalLockNil struct{}

func (LocalLockNil) Error() string { return "local lock is nil" }

type ServerLockNil struct{}

func (ServerLockNil) Error() string { return "server lock is nil" }

type Lock struct {
	local       *localcache.Lock
	localCancel context.CancelFunc
	localMu     *sync.RWMutex

	server       *servercache.Lock
	serverCancel context.CancelFunc
	serverMu     *sync.RWMutex

	scope   int64
	key     string
	typeStr string
}

func newLock(scope int64, key, typeStr string) *Lock {
	return &Lock{
		local:       nil,
		localCancel: nil,
		localMu:     new(sync.RWMutex),

		server:       nil,
		serverCancel: nil,
		serverMu:     new(sync.RWMutex),

		scope:   scope,
		key:     key,
		typeStr: typeStr,
	}
}

func (l *Lock) Key() string     { return l.key }
func (l *Lock) TypeStr() string { return l.typeStr }

func (l *Lock) LocalReleased() (err error) {
	l.localMu.RLock()
	defer l.localMu.RUnlock()
	if l.local == nil {
		err = LocalLockNil{}
	} else {
		err = l.local.Released()
	}
	return err
}
func (l *Lock) LocalLocked() bool {
	return l.LocalReleased() == nil
}

func (l *Lock) ServerReleased() (err error) {
	l.serverMu.RLock()
	defer l.serverMu.RUnlock()
	if l.server == nil {
		err = ServerLockNil{}
	} else {
		err = l.server.Released()
	}
	return err
}
func (l *Lock) ServerLocked() bool {
	return l.ServerReleased() == nil
}

func (l *Lock) localUnlock() {
	l.localMu.RLock()
	defer l.localMu.RUnlock()
	if l.localCancel != nil {
		l.localCancel()
	}
}
func (l *Lock) serverUnlock() {
	l.serverMu.RLock()
	defer l.serverMu.RUnlock()
	if l.serverCancel != nil {
		l.serverCancel()
	}
}

func (l *Lock) localLock(ctx context.Context) (err error) {
	if l.LocalLocked() {
		return nil
	} else if l.localCancel != nil {
		l.localCancel() // doesn't really do anything
	}
	ctx, cancel := context.WithCancel(ctx)
	l.localMu.Lock()
	defer l.localMu.Unlock()
	l.local, err = localcache.ObtainLock(
		ctx,
		l.key,
		l.typeStr,
		LLOCK_MAX_WAIT,
	)
	if err != nil {
		cancel()
	} else {
		l.localCancel = cancel
	}
	return err
}
func (l *Lock) serverLock(ctx context.Context) (err error) {
	if l.ServerLocked() {
		return nil
	} else if l.serverCancel != nil {
		l.serverCancel() // doesn't really do anything
	}
	ctx, cancel := context.WithCancel(ctx)
	l.serverMu.Lock()
	defer l.serverMu.Unlock()
	l.server, err = servercache.ObtainLock(
		ctx,
		l.key,
		SLOCK_TTL,
		SLOCK_MAX_BACKOFF,
	)
	if err != nil {
		cancel()
	} else {
		l.serverCancel = cancel
	}
	return err
}
