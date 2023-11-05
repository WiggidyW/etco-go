package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func get[REP any](
	ctx context.Context,
	method func(context.Context) (REP, error),
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff, expiresIn time.Duration,
	newRep *func() *REP,
) (
	rep REP,
	expires *time.Time,
	err error,
) {
	var repPtr *REP
	repPtr, expires, err = fetch.HandleFetch(
		ctx,
		&prefetch.Params[REP]{
			CacheParams: &prefetch.CacheParams[REP]{
				Get: prefetch.ServerCacheGet(
					typeStr, cacheKey,
					lockTTL, lockMaxBackoff,
					newRep,
				),
			},
		},
		getFetchFunc(method, typeStr, cacheKey, expiresIn),
	)
	if err != nil {
		return rep, nil, err
	} else if repPtr != nil {
		rep = *repPtr
	}
	return rep, expires, nil
}

func getFetchFunc[REP any](
	method func(context.Context) (REP, error),
	typeStr, cacheKey string,
	expiresIn time.Duration,
) fetch.Fetch[REP] {
	return func(ctx context.Context) (
		rep *REP,
		expires *time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var repVal REP
		repVal, err = method(ctx)
		if err != nil {
			return nil, nil, nil, err
		}
		rep = &repVal
		expires = fetch.ExpiresIn(expiresIn)
		postFetch = &postfetch.Params{
			CacheParams: postfetch.ServerCacheSet(typeStr, cacheKey),
		}
		return rep, expires, postFetch, nil
	}
}

func set[REP any](
	ctx context.Context,
	method func(context.Context, REP) error,
	typeStr, cacheKey string,
	lockTTL, lockMaxBackoff, expiresIn time.Duration,
	repVal REP,
) (
	err error,
) {
	_, _, err = fetch.HandleFetch(
		ctx,
		&prefetch.Params[REP]{
			CacheParams: &prefetch.CacheParams[REP]{
				Set: prefetch.ServerCacheSet(
					typeStr, cacheKey,
					lockTTL, lockMaxBackoff,
				),
			},
		},
		setFetchFunc(method, typeStr, cacheKey, expiresIn, repVal),
	)
	return err
}

func setFetchFunc[REP any](
	method func(context.Context, REP) error,
	typeStr, cacheKey string,
	expiresIn time.Duration,
	repVal REP,
) fetch.Fetch[REP] {
	return func(ctx context.Context) (
		rep *REP,
		expires *time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		err = method(ctx, repVal)
		if err != nil {
			return nil, nil, nil, err
		}
		rep = &repVal
		expires = fetch.ExpiresIn(expiresIn)
		postFetch = &postfetch.Params{
			CacheParams: postfetch.ServerCacheSet(typeStr, cacheKey),
		}
		return rep, expires, postFetch, nil
	}
}
