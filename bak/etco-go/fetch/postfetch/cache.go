package postfetch

import (
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/logger"
)

type CacheParams struct {
	Namespace *CacheActionNamespace
	Set       []CacheActionSet
}

type CacheActionNamespace struct {
	CacheKey string
	TypeStr  string
	Expires  time.Time
}

type CacheActionSet struct {
	CacheKey  string
	TypeStr   string
	Expirable any
	Expires   time.Time
	Local     bool
	Server    bool
}

func handleCache(
	x cache.Context,
	params *CacheParams,
	fetchErr error,
) {
	defer x.UnlockScoped()
	if params == nil {
		return
	}

	numSets := len(params.Set)
	if numSets == 1 && params.Namespace == nil {
		logger.MaybeErr(set(x, params.Set[0]))
	} else if numSets != 0 {
		var wg sync.WaitGroup
		wg.Add(numSets)
		for _, action := range params.Set {
			go func(action CacheActionSet) {
				logger.MaybeErr(set(x, action))
				wg.Done()
			}(action)
		}
		defer wg.Wait()
	}

	if params.Namespace != nil {
		logger.MaybeErr(namespaceModify(x, params.Namespace))
	}
}

func namespaceModify(
	x cache.Context,
	action *CacheActionNamespace,
) (err error) {
	return cache.NamespaceModify(
		x,
		action.CacheKey,
		action.TypeStr,
		action.Expires,
	)
}

func set(
	x cache.Context,
	action CacheActionSet,
) (err error) {
	return cache.SetAndUnlock(
		x,
		action.CacheKey, action.TypeStr,
		action.Local, action.Server,
		action.Expirable,
		action.Expires,
	)
}
