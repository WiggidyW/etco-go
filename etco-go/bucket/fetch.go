package bucket

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
)

func get[REP any](
	x cache.Context,
	method func(context.Context) (REP, error),
	cacheKey, typeStr keys.Key,
	expiresIn time.Duration,
	newRep func() REP,
) (
	rep REP,
	expires time.Time,
	err error,
) {
	return fetch.FetchWithCache(
		x,
		getFetchFunc(method, cacheKey, typeStr, expiresIn),
		cacheprefetch.StrongCache(cacheKey, typeStr, newRep, nil),
	)
}

func getFetchFunc[REP any](
	method func(context.Context) (REP, error),
	cacheKey, typeStr keys.Key,
	expiresIn time.Duration,
) fetch.CachingFetch[REP] {
	return func(x cache.Context) (
		rep REP,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		rep, err = method(x.Ctx())
		if err != nil {
			return rep, expires, nil, err
		}
		expires = time.Now().Add(expiresIn)
		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.ServerSetOne[REP](
				cacheKey,
				typeStr,
				rep,
				expires,
			),
		}
		return rep, expires, postFetch, nil
	}
}

func set[REP any](
	x cache.Context,
	method func(context.Context, REP) error,
	cacheKey, typeStr keys.Key,
	expiresIn time.Duration,
	rep REP,
	delDerivative *cacheprefetch.ActionLock,
) (
	err error,
) {
	var cacheLocks []cacheprefetch.ActionOrderedLocks
	if delDerivative != nil {
		cacheLocks = []cacheprefetch.ActionOrderedLocks{{
			Locks: []cacheprefetch.ActionLock{*delDerivative},
			Child: &cacheprefetch.ActionOrderedLocks{
				Locks: []cacheprefetch.ActionLock{
					cacheprefetch.DualLock(cacheKey, typeStr),
				},
			},
		}}
	} else {
		cacheLocks = []cacheprefetch.ActionOrderedLocks{{
			Locks: []cacheprefetch.ActionLock{
				cacheprefetch.DualLock(cacheKey, typeStr),
			},
			Child: nil,
		}}
	}
	_, _, err = fetch.FetchWithCache(
		x,
		setFetchFunc(method, cacheKey, typeStr, expiresIn, rep),
		cacheprefetch.AntiCache[struct{}](cacheLocks),
	)
	return err
}

func setFetchFunc[REP any](
	method func(context.Context, REP) error,
	cacheKey, typeStr keys.Key,
	expiresIn time.Duration,
	rep REP,
) fetch.CachingFetch[struct{}] {
	return func(x cache.Context) (
		_ struct{},
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		err = method(x.Ctx(), rep)
		if err != nil {
			return struct{}{}, expires, nil, err
		}
		expires = time.Now().Add(expiresIn)
		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.ServerSetOne[REP](
				cacheKey,
				typeStr,
				rep,
				expires,
			),
		}
		return struct{}{}, expires, postFetch, nil
	}
}
