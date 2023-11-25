package cachepostfetch

import (
	"sync"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/logger"
)

type Params struct {
	Namespace *ActionNamespace
	Set       []ActionSet
}

type ActionNamespace struct {
	CacheKey keys.Key
	TypeStr  keys.Key
	Expires  time.Time
}

type ActionSet struct {
	CacheKey  keys.Key
	TypeStr   keys.Key
	Expirable any
	Expires   time.Time
	Local     bool
	Server    bool
}

func Handle(
	x cache.Context,
	params *Params,
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
			go func(action ActionSet) {
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
	action *ActionNamespace,
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
	action ActionSet,
) (err error) {
	return cache.SetAndUnlock(
		x,
		action.CacheKey, action.TypeStr,
		action.Local, action.Server,
		action.Expirable,
		action.Expires,
	)
}
