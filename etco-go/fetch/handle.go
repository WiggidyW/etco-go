package fetch

import (
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

type Retry struct {
	Retries     int
	ShouldRetry func(error) bool // sleep inside this func if you want to sleep
}

func HandleFetchVal[REP any](
	x cache.Context,
	preFetchParams *prefetch.Params[REP],
	fetch Fetch[REP],
	retry *Retry,
) (
	rep REP,
	expires time.Time,
	err error,
) {
	var repPtr *REP
	repPtr, expires, err = HandleFetch(x, preFetchParams, fetch, retry)
	if repPtr != nil {
		rep = *repPtr
	}
	return rep, expires, err
}

func HandleFetch[REP any](
	x cache.Context,
	preFetchParams *prefetch.Params[REP],
	fetch Fetch[REP],
	retry *Retry,
) (
	rep *REP,
	expires time.Time,
	err error,
) {
	return handleFetchInner(x.WithNewScope(), preFetchParams, fetch, retry)
}

func handleFetchInner[REP any](
	x cache.Context,
	preFetchParams *prefetch.Params[REP],
	fetch Fetch[REP],
	retry *Retry,
) (
	rep *REP,
	expires time.Time,
	err error,
) {
	var ncRetry bool
	if preFetchParams != nil {
		var expirableRep *expirable.Expirable[REP]
		ncRetry, expirableRep, err = prefetch.Handle(
			x,
			*preFetchParams,
		)
		if err != nil {
			return nil, expires, err
		} else if ncRetry {
			return handleFetchInner(x, preFetchParams, fetch, retry)
		} else if expirableRep != nil {
			rep, expires = expirableRep.Data, expirableRep.Expires
			return rep, expires, nil
		}
	}

	var postFetchParams *postfetch.Params
	if retry != nil {
		rep, expires, postFetchParams, err = fetchWithRetries(
			x,
			fetch,
			*retry,
			0,
		)
	} else {
		rep, expires, postFetchParams, err = fetch(x)
	}

	go postfetch.Handle[REP](x, postFetchParams, err)

	return rep, expires, err
}

func fetchWithRetries[REP any](
	x cache.Context,
	fetch Fetch[REP],
	retry Retry,
	attempt int,
) (
	rep *REP,
	expires time.Time,
	postFetch *postfetch.Params,
	err error,
) {
	rep, expires, postFetch, err = fetch(x)
	if err != nil {
		if attempt < retry.Retries && retry.ShouldRetry(err) {
			return fetchWithRetries(x, fetch, retry, attempt+1)
		}
		return nil, expires, postFetch, err
	}
	return rep, expires, postFetch, nil
}
