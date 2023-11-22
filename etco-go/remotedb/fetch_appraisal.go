package remotedb

import (
	"context"
	"errors"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/cachepostfetch"
	"github.com/WiggidyW/etco-go/fetch/cacheprefetch"
)

// set to local cache only if not-nil
// 1. once an appraisal is set and is non-nil, only then is it immutable
// 2. local cache is non-invalidatable across instances
// therefore, nil appraisals may end up causing stale data if set in local cache
func appraisalGetCacheSetLocal[A any](appraisal *A) bool {
	return appraisal != nil
}

func appraisalGet[A any](
	x cache.Context,
	method func(context.Context, string) (*A, error),
	typeStr keys.Key,
	code string,
	expiresIn time.Duration,
) (
	rep *A,
	expires time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyAppraisal(code)
	return fetch.FetchWithCache(
		x,
		appraisalGetFetchFunc[A](
			method,
			cacheKey,
			typeStr,
			code,
			expiresIn,
		),
		cacheprefetch.WeakCache[*A](
			cacheKey,
			typeStr,
			nil,
			appraisalGetCacheSetLocal,
			nil,
		),
	)
}

func appraisalGetFetchFunc[A any](
	method func(context.Context, string) (*A, error),
	cacheKey, typeStr keys.Key,
	code string,
	expiresIn time.Duration,
) fetch.CachingFetch[*A] {
	return func(x cache.Context) (
		rep *A,
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		rep, err = method(x.Ctx(), code)
		if err != nil {
			return nil, expires, nil, err
		}
		expires = time.Now().Add(expiresIn)
		if appraisalGetCacheSetLocal(rep) {
			postFetch = &cachepostfetch.Params{Set: cachepostfetch.DualSetOne(
				cacheKey,
				typeStr,
				rep,
				expires,
			)}
		} else {
			postFetch = &cachepostfetch.Params{Set: cachepostfetch.ServerSetOne(
				cacheKey,
				typeStr,
				rep,
				expires,
			)}
		}
		return rep, expires, postFetch, nil
	}
}

func appraisalSet[A Appraisal](
	x cache.Context,
	method func(context.Context, A) error,
	typeStr keys.Key,
	expiresIn time.Duration,
	appraisal A,
	cacheLocks []cacheprefetch.ActionOrderedLocks,
) (
	err error,
) {
	code := appraisal.GetCode()
	if code == "" {
		return errors.New("unable to set appraisal without a code")
	}
	cacheKey := keys.CacheKeyAppraisal(code)
	_, _, err = fetch.FetchWithCache(
		x,
		appraisalSetFetchFunc[A](
			method,
			cacheKey, typeStr,
			expiresIn,
			appraisal,
		),
		cacheprefetch.AntiCache[struct{}](append(
			cacheLocks,
			cacheprefetch.ActionOrderedLocks{
				Locks: []cacheprefetch.ActionLock{
					cacheprefetch.DualLock(cacheKey, typeStr),
				},
				Child: nil,
			},
		)),
	)
	return err
}

func appraisalSetFetchFunc[A any](
	method func(context.Context, A) error,
	cacheKey, typeStr keys.Key,
	expiresIn time.Duration,
	appraisal A,
) fetch.CachingFetch[struct{}] {
	return func(x cache.Context) (
		_ struct{},
		expires time.Time,
		postFetch *cachepostfetch.Params,
		err error,
	) {
		err = method(x.Ctx(), appraisal)
		if err != nil {
			return struct{}{}, expires, nil, err
		}
		expires = time.Now().Add(expiresIn)
		postFetch = &cachepostfetch.Params{
			Set: cachepostfetch.DualSetOne[*A](
				cacheKey, typeStr,
				&appraisal,
				expires,
			),
		}
		return struct{}{}, expires, postFetch, nil
	}
}
