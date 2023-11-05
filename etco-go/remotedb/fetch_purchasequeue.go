package remotedb

import (
	"context"
	"time"

	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
)

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
