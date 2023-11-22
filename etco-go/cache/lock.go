package cache

import (
	"context"
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/cache/localcache"
	"github.com/WiggidyW/etco-go/cache/servercache"
)

const (
	LLOCK_MAX_WAIT    time.Duration = 4 * time.Second
	SLOCK_TTL         time.Duration = 4 * time.Second
	SLOCK_MAX_BACKOFF time.Duration = 4 * time.Second
)

type LockNil struct{}

func (LockNil) Error() string { return "lock is nil" }

type Lock struct {
	local   *innerLockWrapper[*localcache.Lock]
	server  *innerLockWrapper[*servercache.Lock]
	ctx     context.Context
	key     keys.Key
	typeStr keys.Key
}

func newLock(
	ctx context.Context,
	key keys.Key,
	typeStr keys.Key,
) *Lock {
	return &Lock{
		local:   newInnerLockWrapper[*localcache.Lock](),
		server:  newInnerLockWrapper[*servercache.Lock](),
		ctx:     ctx,
		key:     key,
		typeStr: typeStr,
	}
}

func (l *Lock) localUnlock(scope int64)  { l.local.unlock(scope) }
func (l *Lock) serverUnlock(scope int64) { l.server.unlock(scope) }
func (l *Lock) localLock(scope int64) (err error) {
	return l.local.lock(
		l.ctx,
		scope,
		func(ctx context.Context) (*localcache.Lock, error) {
			return localcache.ObtainLock(
				ctx,
				l.key,
				l.typeStr,
				LLOCK_MAX_WAIT,
			)
		},
	)
}
func (l *Lock) serverLock(scope int64) (err error) {
	return l.server.lock(
		l.ctx,
		scope,
		func(ctx context.Context) (*servercache.Lock, error) {
			return servercache.ObtainLock(
				ctx,
				l.key,
				SLOCK_TTL,
				SLOCK_MAX_BACKOFF,
			)
		},
	)
}

type innerLock interface {
	IsNil() bool
	Released() error
}

type innerLockWrapper[L innerLock] struct {
	innerLock L
	cancel    context.CancelFunc
	scopes    map[int64]struct{}
	mu        *sync.RWMutex
}

func newInnerLockWrapper[L innerLock]() *innerLockWrapper[L] {
	var innerLock L
	return &innerLockWrapper[L]{
		innerLock: innerLock,
		cancel:    nil,
		scopes:    make(map[int64]struct{}),
		mu:        new(sync.RWMutex),
	}
}

func (ilw *innerLockWrapper[L]) released() (err error) {
	if ilw.innerLock.IsNil() {
		err = LockNil{}
	} else {
		err = ilw.innerLock.Released()
	}
	return err
}

func (ilw *innerLockWrapper[L]) locked() bool {
	return ilw.released() == nil
}

func (ilw *innerLockWrapper[L]) unlock(scope int64) {
	ilw.mu.Lock()
	defer ilw.mu.Unlock()
	delete(ilw.scopes, scope)
	if len(ilw.scopes) == 0 {
		ilw.cancel()
	}
}

func (ilw *innerLockWrapper[L]) lock(
	ctx context.Context,
	scope int64,
	obtain func(context.Context) (L, error),
) (err error) {
	ilw.mu.Lock()
	defer ilw.mu.Unlock()
	ilw.scopes[scope] = struct{}{}
	if ilw.locked() {
		return nil
	} else if ilw.cancel != nil {
		ilw.cancel() // I think this does nothing, but it's cheap and sound
	}
	ctx, cancel := context.WithCancel(ctx)
	ilw.innerLock, err = obtain(ctx)
	if err != nil {
		cancel() // I think this also does nothing, but it's cheap and sound
	} else {
		ilw.cancel = cancel
	}
	return err
}
