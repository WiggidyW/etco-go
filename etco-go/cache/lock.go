package cache

import (
	"sync"

	"github.com/WiggidyW/etco-go/cache/servercache"
	"github.com/WiggidyW/etco-go/logger"
)

type Lock struct {
	released bool
	local    *sync.Mutex
	server   *servercache.Lock
	mu       *sync.RWMutex
}

func newLock(
	local *sync.Mutex,
	server *servercache.Lock,
) *Lock {
	return &Lock{
		released: false,
		local:    local,
		server:   server,
		mu:       &sync.RWMutex{},
	}
}

func (l *Lock) Unlock() (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.released {
		l.local.Unlock()
		err = l.server.Unlock()
		l.released = true
	}
	return err
}

func (l *Lock) UnlockLogErr() {
	err := l.Unlock()
	if err != nil {
		logger.Err(err.Error())
	}
}

func UnlockManyLogErr(locks []*Lock) {
	for _, lock := range locks {
		go lock.UnlockLogErr()
	}
}
