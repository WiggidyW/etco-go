package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
)

type ContextLocks struct {
	locks  map[[16]byte]*Lock
	ctx    context.Context
	cancel context.CancelFunc
	mu     *sync.RWMutex
}

type Context struct {
	locks *ContextLocks
	ctx   context.Context
	scope int64
}

func NewContext(ctx context.Context) Context {
	locksCtx, locksCancel := context.WithCancel(context.Background())
	return Context{
		locks: &ContextLocks{
			locks:  make(map[[16]byte]*Lock),
			mu:     new(sync.RWMutex),
			ctx:    locksCtx,
			cancel: locksCancel,
		},
		ctx:   ctx,
		scope: 0,
	}
}

func (x Context) Ctx() context.Context { return x.ctx }

func (x Context) WithCancel() (
	Context,
	context.CancelFunc,
) {
	var cancel context.CancelFunc
	x.ctx, cancel = context.WithCancel(x.ctx)
	return x, cancel
}

func (x Context) WithTimeout(timeout time.Duration) (
	Context,
	context.CancelFunc,
) {
	var cancel context.CancelFunc
	x.ctx, cancel = context.WithTimeout(x.ctx, timeout)
	return x, cancel
}

func (x Context) WithNewScope() Context {
	x.scope = x.newScope()
	return x
}

func (x Context) Background() Context {
	x.ctx = context.Background()
	return x
}

func (x Context) Unlock(key, typeStr keys.Key) {
	x.locks.mu.RLock()
	lock, ok := x.locks.locks[key.Bytes16()]
	x.locks.mu.RUnlock()
	if ok {
		go x.localUnlock(lock)
		go x.serverUnlock(lock)
	}
}

func (x Context) UnlockAll() { x.locks.cancel() }

func (x Context) UnlockScoped() {
	x.locks.mu.RLock()
	defer x.locks.mu.RUnlock()
	for _, lock := range x.locks.locks {
		go lock.localUnlock(x.scope)
		go lock.serverUnlock(x.scope)
	}
}

func (x Context) getLock(key, typeStr keys.Key) (lock *Lock) {
	var lockOk bool
	x.locks.mu.RLock()
	lock, lockOk = x.locks.locks[key.Bytes16()]
	x.locks.mu.RUnlock()
	if !lockOk {
		lock = newLock(x.locks.ctx, key, typeStr)
		x.locks.mu.Lock()
		defer x.locks.mu.Unlock()
		x.locks.locks[key.Bytes16()] = lock
	}
	return lock
}

func (x Context) localLock(lock *Lock) error  { return lock.localLock(x.scope) }
func (x Context) serverLock(lock *Lock) error { return lock.serverLock(x.scope) }
func (x Context) localUnlock(lock *Lock)      { lock.localUnlock(x.scope) }
func (x Context) serverUnlock(lock *Lock)     { lock.serverUnlock(x.scope) }
func (x Context) newScope() int64             { return rand.Int63() }
