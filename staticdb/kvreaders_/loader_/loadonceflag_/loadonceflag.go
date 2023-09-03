package loadonceflag_

import (
	"sync"
)

type LoadOnceFlag struct {
	flag   bool          // false if loading, or CPU caching is out of date (acceptable data race)
	rwLock *sync.RWMutex // initially write locked by load thread
}

// returns a LoadOnceFlag that is safe to use immediately
func NewLoadOnceFlag() *LoadOnceFlag {
	lof := newUnsafeLoadOnceFlag()
	lof.loadStart()
	return lof
}

// returns a LoadOnceFlag that is not safe to use until loadStart() is called
func newUnsafeLoadOnceFlag() *LoadOnceFlag {
	return &LoadOnceFlag{false, &sync.RWMutex{}}
}

// call when load thread is started
func (lof *LoadOnceFlag) loadStart() {
	lof.rwLock.Lock() // write lock
}

// call when load thread is finished
func (lof *LoadOnceFlag) LoadFinish() {
	lof.flag = true     // update the flag
	lof.rwLock.Unlock() // write unlock
}

// block until LoadFinish() is called
func (lof *LoadOnceFlag) Check() {
	if !lof.flag {
		// sychronize with the load thread
		// (won't block other check threads)
		lof.rwLock.RLock()
		defer lof.rwLock.RUnlock()
		lof.flag = true
	}
}
