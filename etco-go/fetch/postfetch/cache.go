package postfetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

type CacheParams struct {
	TypeStr   string
	Key       string
	LocalSet  bool
	ServerSet bool
}

func handleCache[REP any](
	locks prefetch.CacheLocks,
	params *CacheParams,
	expires *time.Time,
	rep *REP,
) (err error) {
	if locks.Del != nil {
		// defer to not unlock until set op has gone through (if there is one)
		defer cache.UnlockManyLogErr(*locks.Del)
	}

	if params != nil {
		return cache.SetAndUnlock(
			params.Key,
			params.TypeStr,
			params.LocalSet,
			params.ServerSet,
			locks.Set,
			rep,
			expires,
		)
	} else if locks.Set != nil {
		go locks.Set.UnlockLogErr()
	}

	return nil
}
