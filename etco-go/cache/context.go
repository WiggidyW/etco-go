package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
)

type CtxWithCancel struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type ContextLocks struct {
	locks    map[[16]byte]*Lock
	mu       *sync.RWMutex
	ctxAll   CtxWithCancel
	ctxScope map[int64]CtxWithCancel
}

type Context struct {
	locks *ContextLocks
	ctx   context.Context
	scope int64
}

func NewContext(ctx context.Context) Context {
	locksCtxAll, locksCancelAll := context.WithCancel(context.Background())
	return Context{
		locks: &ContextLocks{
			locks:    make(map[[16]byte]*Lock),
			mu:       new(sync.RWMutex),
			ctxScope: make(map[int64]CtxWithCancel),
			ctxAll: CtxWithCancel{
				ctx:    locksCtxAll,
				cancel: locksCancelAll,
			},
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

func (x Context) GetLock(key, typeStr keys.Key) (lock *Lock) {
	var lockOk bool

	x.locks.mu.RLock()
	lock, lockOk = x.locks.locks[key.Buf]
	x.locks.mu.RUnlock()

	if !lockOk {
		lock = newLock(x.scope, key, typeStr)
		x.locks.mu.Lock()
		defer x.locks.mu.Unlock()

		x.locks.locks[key.Buf] = lock
		_, ctxScopeOk := x.locks.ctxScope[x.scope]
		if !ctxScopeOk {
			ctxScope, ctxScopeCancel := context.WithCancel(x.locks.ctxAll.ctx)
			x.locks.ctxScope[x.scope] = CtxWithCancel{
				ctx:    ctxScope,
				cancel: ctxScopeCancel,
			}
		}
	}

	return lock
}

func (x Context) Background() Context {
	x.ctx = context.Background()
	return x
}

func (x Context) LocalLock(lock *Lock) error {
	return lock.localLock(x.locks.ctxScope[lock.scope].ctx)
}
func (x Context) ServerLock(lock *Lock) error {
	return lock.serverLock(x.locks.ctxScope[lock.scope].ctx)
}

func (x Context) Unlock(key, typeStr keys.Key) {
	x.locks.mu.RLock()
	lock, ok := x.locks.locks[key.Buf]
	x.locks.mu.RUnlock()
	if ok && x.scope == lock.scope {
		go lock.localUnlock()
		go lock.serverUnlock()
	}
}
func (x Context) LocalUnlock(lock *Lock) {
	if x.scope == lock.scope {
		lock.localUnlock()
	}
}
func (x Context) ServerUnlock(lock *Lock) {
	if x.scope == lock.scope {
		lock.serverUnlock()
	}
}

func (x Context) UnlockAll() {
	x.locks.ctxAll.cancel()
}

func (x Context) UnlockScoped() {
	ctxScope, ok := x.locks.ctxScope[x.scope]
	if ok {
		ctxScope.cancel()
	}
}

func (x Context) newScope() int64 {
	return rand.Int63()
}
