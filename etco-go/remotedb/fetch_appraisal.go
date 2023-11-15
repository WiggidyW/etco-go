package remotedb

import (
	"context"
	"errors"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
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
	typeStr, code string,
	expiresIn time.Duration,
) (
	rep *A,
	expires time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyAppraisal(code)
	return fetch.HandleFetch(
		x,
		&prefetch.Params[*A]{
			CacheParams: &prefetch.CacheParams[*A]{
				Get: prefetch.DualCacheGet[*A](
					cacheKey, typeStr,
					true,
					nil,
					appraisalGetCacheSetLocal,
				),
			},
		},
		appraisalGetFetchFunc[A](
			method,
			cacheKey, typeStr, code,
			expiresIn,
		),
		nil,
	)
}

func appraisalGetFetchFunc[A any](
	method func(context.Context, string) (*A, error),
	cacheKey, typeStr, code string,
	expiresIn time.Duration,
) fetch.Fetch[*A] {
	return func(x cache.Context) (
		rep *A,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		rep, err = method(x.Ctx(), code)
		if err != nil {
			return nil, expires, nil, err
		}
		expires = time.Now().Add(expiresIn)
		var set []postfetch.CacheActionSet
		if appraisalGetCacheSetLocal(rep) {
			set = postfetch.DualCacheSetOne(cacheKey, typeStr, rep, expires)
		} else {
			set = postfetch.ServerCacheSetOne(cacheKey, typeStr, rep, expires)
		}
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: set,
			},
		}
		return rep, expires, postFetch, nil
	}
}

func appraisalSet[A Appraisal](
	x cache.Context,
	method func(context.Context, A) error,
	typeStr string,
	expiresIn time.Duration,
	appraisal A,
	cacheLocks []prefetch.CacheActionOrderedLocks,
) (
	err error,
) {
	code := appraisal.GetCode()
	if code == "" {
		return errors.New("unable to set appraisal without a code")
	}
	cacheKey := keys.CacheKeyAppraisal(code)
	if cacheLocks != nil {
		cacheLocks = append(
			cacheLocks,
			prefetch.CacheOrderedLocks(
				nil,
				prefetch.DualCacheLock(cacheKey, typeStr),
			),
		)
	} else {
		cacheLocks = prefetch.DualCacheOrderedLocksOne(cacheKey, typeStr)
	}
	_, _, err = fetch.HandleFetch(
		x,
		&prefetch.Params[struct{}]{
			CacheParams: &prefetch.CacheParams[struct{}]{
				Lock: cacheLocks,
			},
		},
		appraisalSetFetchFunc[A](
			method,
			cacheKey, typeStr,
			expiresIn,
			appraisal,
		),
		nil,
	)
	return err
}

func appraisalSetFetchFunc[A any](
	method func(context.Context, A) error,
	cacheKey, typeStr string,
	expiresIn time.Duration,
	appraisal A,
) fetch.Fetch[struct{}] {
	return func(x cache.Context) (
		_ struct{},
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		err = method(x.Ctx(), appraisal)
		if err != nil {
			return struct{}{}, expires, nil, err
		}
		expires = time.Now().Add(expiresIn)
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.DualCacheSetOne[*A](
					cacheKey, typeStr,
					&appraisal,
					expires,
				),
			},
		}
		return struct{}{}, expires, postFetch, nil
	}
}
