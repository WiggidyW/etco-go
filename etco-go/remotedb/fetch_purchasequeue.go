package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

func purchaseQueueCancel(
	ctx context.Context,
	method func(context.Context) error,
	lockTTL, lockMaxBackoff time.Duration,
	cacheDels *[]prefetch.CacheAction,
) (
	err error,
) {
	if cacheDels == nil {
		cacheDelsVal := make([]prefetch.CacheAction, 1)
		cacheDels = &cacheDelsVal
	}
	*cacheDels = append(*cacheDels, prefetch.ServerCacheDel(
		keys.TypeStrPurchaseQueue, keys.CacheKeyPurchaseQueue,
		lockTTL, lockMaxBackoff,
	))
	_, _, err = fetch.HandleFetch(
		ctx,
		&prefetch.Params[struct{}]{
			CacheParams: &prefetch.CacheParams[struct{}]{
				Del: cacheDels,
			},
		},
		purchaseQueueCancelFetchFunc(method),
	)
	return err
}

func purchaseQueueCancelFetchFunc(
	method func(context.Context) error,
) fetch.Fetch[struct{}] {
	return func(ctx context.Context) (
		*struct{},
		*time.Time,
		*postfetch.Params,
		error,
	) {
		return nil, nil, nil, method(ctx)
	}
}

func purchaseQueueGet(
	ctx context.Context,
	lockTTL, lockMaxBackoff, expiresIn time.Duration,
) (
	rep []string,
	expires *time.Time,
	err error,
) {
	var repPtr *[]string
	repPtr, expires, err = fetch.HandleFetch(
		ctx,
		&prefetch.Params[[]string]{
			CacheParams: &prefetch.CacheParams[[]string]{
				Get: prefetch.ServerCacheGet[[]string](
					keys.TypeStrPurchaseQueue, keys.CacheKeyPurchaseQueue,
					lockTTL, lockMaxBackoff,
					nil,
				),
			},
		},
		purchaseQueueGetFetchFunc(expiresIn),
	)
	if err != nil {
		return nil, nil, err
	} else if repPtr != nil {
		rep = *repPtr
	}
	return rep, expires, nil
}

func purchaseQueueGetFetchFunc(expiresIn time.Duration) fetch.Fetch[[]string] {
	return func(ctx context.Context) (
		rep *[]string,
		expires *time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var rdbRep *PurchaseQueue
		rdbRep, err = client.readPurchaseQueue(ctx)
		if err != nil {
			return nil, nil, nil, err
		}
		if rdbRep != nil {
			rep = &rdbRep.PurchaseQueue
		}
		expires = fetch.ExpiresIn(expiresIn)
		postFetch = &postfetch.Params{
			CacheParams: postfetch.ServerCacheSet(
				keys.TypeStrPurchaseQueue,
				keys.CacheKeyPurchaseQueue,
			),
		}
		return rep, expires, postFetch, nil
	}
}
