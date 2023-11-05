package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

// set to local cache only if not-nil
// 1. once an appraisal is set and is non-nil, only then is it immutable
// 2. local cache is non-invalidatable
// therefore, nil appraisals may end up causing stale data if set in local cache
func appraisalGetCacheSetLocal[A any](appraisal *A) bool {
	return appraisal != nil
}

func appraisalGet[A any](
	ctx context.Context,
	method func(context.Context, string) (*A, error),
	typeStr, code string,
	lockTTL, lockMaxBackoff, expiresIn time.Duration,
) (
	rep *A,
	expires *time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyAppraisal(code)
	rep, expires, err = fetch.HandleFetch(
		ctx,
		&prefetch.Params[A]{
			CacheParams: &prefetch.CacheParams[A]{
				Get: prefetch.DualCacheGet[A](
					typeStr, cacheKey,
					lockTTL, lockMaxBackoff,
					nil,
					cache.NewSloshFunc(appraisalGetCacheSetLocal[A]),
				),
			},
		},
		appraisalGetFetchFunc[A](
			method,
			typeStr, cacheKey, code,
			expiresIn,
		),
	)
	return rep, expires, err
}

func appraisalGetFetchFunc[A any](
	method func(context.Context, string) (*A, error),
	typeStr, cacheKey, code string,
	expiresIn time.Duration,
) fetch.Fetch[A] {
	return func(ctx context.Context) (
		rep *A,
		expires *time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		rep, err = method(ctx, code)
		if err != nil {
			return nil, nil, nil, err
		}
		expires = fetch.ExpiresIn(expiresIn)
		if appraisalGetCacheSetLocal(rep) {
			postFetch = &postfetch.Params{
				CacheParams: postfetch.DualCacheSet(typeStr, cacheKey),
			}
		} else {
			postFetch = &postfetch.Params{
				CacheParams: postfetch.ServerCacheSet(typeStr, cacheKey),
			}
		}
		return rep, expires, postFetch, nil
	}
}

func appraisalSet[A Appraisal](
	ctx context.Context,
	method func(context.Context, A) error,
	typeStr string,
	lockTTL, lockMaxBackoff, expiresIn time.Duration,
	rep A,
	cacheDels *[]prefetch.CacheAction,
) (
	err error,
) {
	cacheKey := keys.CacheKeyAppraisal(rep.GetCode())
	_, _, err = fetch.HandleFetch(
		ctx,
		&prefetch.Params[A]{
			CacheParams: &prefetch.CacheParams[A]{
				Set: prefetch.DualCacheSet(
					typeStr, cacheKey,
					lockTTL, lockMaxBackoff,
				),
				Del: cacheDels,
			},
		},
		appraisalSetFetchFunc[A](
			method,
			typeStr, cacheKey,
			expiresIn,
			rep,
		),
	)
	return err
}

func appraisalSetFetchFunc[A any](
	method func(context.Context, A) error,
	typeStr, cacheKey string,
	expiresIn time.Duration,
	rep A,
) fetch.Fetch[A] {
	return func(ctx context.Context) (
		_ *A,
		expires *time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		err = method(ctx, rep)
		if err != nil {
			return nil, nil, nil, err
		}
		expires = fetch.ExpiresIn(expiresIn)
		postFetch = &postfetch.Params{
			CacheParams: postfetch.DualCacheSet(typeStr, cacheKey),
		}
		return &rep, expires, postFetch, nil
	}
}
