package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func get[REP any](
	x cache.Context,
	method func(context.Context) (REP, error),
	cacheKey, typeStr string,
	expiresIn time.Duration,
	newRep func() *REP,
) (
	rep REP,
	expires time.Time,
	err error,
) {
	return fetch.HandleFetchVal(
		x,
		&prefetch.Params[REP]{
			CacheParams: &prefetch.CacheParams[REP]{
				Get: prefetch.ServerCacheGet[REP](
					cacheKey, typeStr,
					true,
					newRep,
				),
			},
		},
		getFetchFunc(method, cacheKey, typeStr, expiresIn),
		nil,
	)
}

func getFetchFunc[REP any](
	method func(context.Context) (REP, error),
	cacheKey, typeStr string,
	expiresIn time.Duration,
) fetch.Fetch[REP] {
	return func(x cache.Context) (
		repPtr *REP,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var rep REP
		rep, err = method(x.Ctx())
		if err != nil {
			return nil, expires, nil, err
		}
		expires = time.Now().Add(expiresIn)
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.ServerCacheSetOne(
					cacheKey, typeStr,
					&rep,
					expires,
				),
			},
		}
		return &rep, expires, postFetch, nil
	}
}

func set[REP any](
	x cache.Context,
	method func(context.Context, REP) error,
	cacheKey, typeStr string,
	expiresIn time.Duration,
	rep REP,
	delDerivative *prefetch.CacheActionLock,
) (
	err error,
) {
	var cacheLocks []prefetch.CacheActionOrderedLocks
	if delDerivative != nil {
		cacheLocks = []prefetch.CacheActionOrderedLocks{
			prefetch.CacheOrderedLocks(
				prefetch.CacheOrderedLocksPtr(
					nil,
					prefetch.ServerCacheLock(cacheKey, typeStr),
				),
				*delDerivative,
			),
		}
	} else {
		cacheLocks = prefetch.ServerCacheOrderedLocksOne(cacheKey, typeStr)
	}
	_, _, err = fetch.HandleFetch(
		x,
		&prefetch.Params[struct{}]{
			CacheParams: &prefetch.CacheParams[struct{}]{
				Lock: cacheLocks,
			},
		},
		setFetchFunc(method, cacheKey, typeStr, expiresIn, rep),
		nil,
	)
	return err
}

func setFetchFunc[REP any](
	method func(context.Context, REP) error,
	cacheKey, typeStr string,
	expiresIn time.Duration,
	rep REP,
) fetch.Fetch[struct{}] {
	return func(x cache.Context) (
		_ *struct{},
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		err = method(x.Ctx(), rep)
		if err != nil {
			return nil, expires, nil, err
		}
		expires = time.Now().Add(expiresIn)
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.ServerCacheSetOne(
					cacheKey, typeStr,
					&rep,
					expires,
				),
			},
		}
		return nil, expires, postFetch, nil
	}
}
