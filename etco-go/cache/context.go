package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
)

type ScopedLocks struct {
	ctx    context.Context
	cancel context.CancelFunc

	parent *ScopedLocks
	locks  map[[16]byte]*Lock
	scope  int64
	mu     *sync.RWMutex
}

func newScopedLocks() ScopedLocks {
	ctx, cancel := context.WithCancel(context.Background())
	return ScopedLocks{
		parent: nil,
		ctx:    ctx,
		cancel: cancel,
		locks:  make(map[[16]byte]*Lock),
		scope:  0,
		mu:     new(sync.RWMutex),
	}
}

func (sl ScopedLocks) withNewScope() ScopedLocks {
	sl.mu.RLock()
	locks := make(map[[16]byte]*Lock, len(sl.locks))
	for key, lock := range sl.locks {
		locks[key] = lock
	}
	sl.mu.RUnlock()

	sl.locks = locks
	sl.scope = sl.newScope()
	sl.mu = new(sync.RWMutex)
	return sl
}

func (sl ScopedLocks) newScope() int64 {
	return rand.Int63()
}

func (sl ScopedLocks) unlock(key, typeStr keys.Key) {
	sl.mu.RLock()
	lock, ok := sl.locks[key.Bytes16()]
	sl.mu.RUnlock()
	if ok {
		go sl.localUnlock(lock)
		go sl.serverUnlock(lock)
	}
}

func (sl ScopedLocks) unlockScoped() {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	for _, lock := range sl.locks {
		go sl.localUnlock(lock)
		go sl.serverUnlock(lock)
	}
}

func (sl ScopedLocks) unlockAll() { sl.cancel() }

func (sl ScopedLocks) getLockForChild(key, typeStr keys.Key) (
	lock *Lock,
	ok bool,
) {
	sl.mu.RLock()
	lock, ok = sl.locks[key.Bytes16()]
	sl.mu.RUnlock()
	if !ok && sl.parent != nil {
		return sl.parent.getLockForChild(key, typeStr)
	} else {
		return lock, ok
	}
}

func (sl ScopedLocks) getLock(key, typeStr keys.Key) (lock *Lock) {
	var ok bool
	sl.mu.Lock()
	defer sl.mu.Unlock()
	lock, ok = sl.locks[key.Bytes16()]
	if !ok && sl.parent != nil {
		lock, ok = sl.parent.getLockForChild(key, typeStr)
		if ok {
			sl.locks[key.Bytes16()] = lock
		}
	}
	if !ok {
		lock = newLock(sl.ctx, key, typeStr)
		sl.locks[key.Bytes16()] = lock
	}
	return lock
}

func (sl ScopedLocks) localLock(lock *Lock) error {
	return lock.localLock(sl.scope)
}
func (sl ScopedLocks) serverLock(lock *Lock) error {
	return lock.serverLock(sl.scope)
}
func (sl ScopedLocks) localUnlock(lock *Lock) {
	lock.localUnlock(sl.scope)
}
func (sl ScopedLocks) serverUnlock(lock *Lock) {
	lock.serverUnlock(sl.scope)
}

type Context struct {
	ctx         context.Context
	scopedLocks ScopedLocks
}

func NewContext(ctx context.Context) Context {
	return Context{ctx: ctx, scopedLocks: newScopedLocks()}
}

func (x Context) Scope() int64         { return x.scopedLocks.scope }
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
	x.scopedLocks = x.scopedLocks.withNewScope()
	return x
}

func (x Context) Background() Context {
	x.ctx = context.Background()
	return x
}

func (x Context) Unlock(key, typeStr keys.Key) {
	x.scopedLocks.unlock(key, typeStr)
}
func (x Context) UnlockAll() {
	x.scopedLocks.unlockAll()
}
func (x Context) UnlockScoped() {
	x.scopedLocks.unlockScoped()
}
func (x Context) getLock(key, typeStr keys.Key) (lock *Lock) {
	return x.scopedLocks.getLock(key, typeStr)
}
func (x Context) localLock(lock *Lock) error {
	return x.scopedLocks.localLock(lock)
}
func (x Context) serverLock(lock *Lock) error {
	return x.scopedLocks.serverLock(lock)
}
func (x Context) localUnlock(lock *Lock) {
	x.scopedLocks.localUnlock(lock)
}
func (x Context) serverUnlock(lock *Lock) {
	x.scopedLocks.serverUnlock(lock)
}
