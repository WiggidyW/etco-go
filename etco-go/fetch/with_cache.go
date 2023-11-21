package fetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
)

func FetchWithCache[REP any](
	x cache.Context,
	fetch CachingFetch[REP],
	preFetch cacheprefetch.Params[REP],
) (
	rep REP,
	expires time.Time,
	err error,
) {
	x = x.WithNewScope() // no new scope across namespace retries
	return fetchWithCacheInner(x, fetch, preFetch)
}

func fetchWithCacheInner[REP any](
	x cache.Context,
	fetch CachingFetch[REP],
	preFetch cacheprefetch.Params[REP],
) (
	rep REP,
	expires time.Time,
	err error,
) {
	var namespaceRetry bool
	var expirableRep *expirable.Expirable[REP]
	namespaceRetry, expirableRep, err = cacheprefetch.Handle(x, preFetch)
	if err != nil {
		return rep, expires, err
	} else if namespaceRetry {
		return fetchWithCacheInner(x, fetch, preFetch)
	} else if expirableRep != nil {
		rep, expires = expirableRep.Data, expirableRep.Expires
		return rep, expires, nil
	}

	var postFetch *cachepostfetch.Params
	rep, expires, postFetch, err = fetch(x)
	go cachepostfetch.Handle(x, postFetch, err)

	return rep, expires, err
}
