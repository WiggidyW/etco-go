package purchasequeue

import (
	"slices"
	"time"

	build "github.com/WiggidyW/etco-go/buildconstants"
	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/cache/expirable"
	"github.com/WiggidyW/etco-go/cache/keys"
	"github.com/WiggidyW/etco-go/contracts"
	"github.com/WiggidyW/etco-go/fetch"
	"github.com/WiggidyW/etco-go/fetch/postfetch"
	"github.com/WiggidyW/etco-go/fetch/prefetch"
	"github.com/WiggidyW/etco-go/logger"
	"github.com/WiggidyW/etco-go/remotedb"
)

func userCancel(
	x cache.Context,
	characterId int32,
	code string,
	locationId int64,
) (
	status CancelPurchaseStatus,
	err error,
) {
	x, cancel := x.WithCancel()
	defer cancel()

	// fetch location purchase queue in a goroutine
	chnQueue := expirable.NewChanResult[[]string](x.Ctx(), 1, 0)
	go expirable.Param2Transceive(
		chnQueue,
		x, locationId,
		GetLocationPurchaseQueue,
	)

	// fetch user appraisal codes in a goroutine
	chnUserCodes := expirable.NewChanResult[[]string](x.Ctx(), 1, 0)
	go expirable.Param2Transceive(
		chnUserCodes,
		x, characterId,
		remotedb.GetUserShopAppraisalCodes,
	)

	// fetch user cancelled purchase time and check cooldown
	var userCancelled *time.Time
	userCancelled, _, err = remotedb.GetUserCancelledPurchase(x, characterId)
	if err != nil {
		return status, err
	} else if userCancelled != nil && time.Now().Before(
		(*userCancelled).Add(build.CANCEL_PURCHASE_COOLDOWN),
	) {
		status = CPS_CooldownLimited
		return status, nil
	}

	// recv user codes and check code exists for user
	var userCodes []string
	userCodes, _, err = chnUserCodes.RecvExp()
	if err != nil {
		return status, err
	} else if !slices.Contains(userCodes, code) {
		status = CPS_PurchaseNotFound
		return status, nil
	}

	// cancel the purchase
	status = CPS_Success
	err = remotedb.UserCancelPurchase(x, characterId, code, locationId)
	return status, err
}

func locationGet(
	x cache.Context,
	locationId int64,
) (
	rep []string,
	expires time.Time,
	err error,
) {
	cacheKey := keys.CacheKeyLocationPurchaseQueue(locationId)
	return fetch.HandleFetch(
		x,
		&prefetch.Params[[]string]{
			CacheParams: &prefetch.CacheParams[[]string]{
				Get: prefetch.ServerCacheGet[[]string](
					cacheKey,
					keys.TypeStrLocationPurchaseQueue,
					true,
					nil,
				),
			},
		},
		locationGetFetchFunc(locationId, cacheKey),
		nil,
	)
}

func locationGetFetchFunc(
	locationId int64,
	cacheKey string,
) fetch.Fetch[[]string] {
	return func(x cache.Context) (
		repPtr []string,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var purchaseQueue PurchaseQueue
		purchaseQueue, expires, err = get(x).RecvExp()
		if err != nil {
			return nil, expires, nil, err
		}
		rep := purchaseQueue[locationId]
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.ServerCacheSetOne(
					cacheKey,
					keys.TypeStrLocationPurchaseQueue,
					&rep,
					expires,
				),
			},
		}
		return rep, expires, postFetch, nil
	}
}

func get(
	x cache.Context,
) (
	chn expirable.ChanResult[PurchaseQueue],
) {
	chn = expirable.NewChanResult[PurchaseQueue](x.Ctx(), 1, 0)
	go transceiveGet(x, chn)
	return chn
}

func transceiveGet(
	x cache.Context,
	chn expirable.ChanResult[PurchaseQueue],
) (
	err error, // context errors only
) {
	chnFetched := make(chan struct{}) // track whether rep came from cache or fetch

	var purchaseQueue PurchaseQueue
	var expires time.Time
	purchaseQueue, expires, err = fetch.HandleFetch(
		x,
		&prefetch.Params[PurchaseQueue]{
			CacheParams: &prefetch.CacheParams[PurchaseQueue]{
				Get: prefetch.ServerCacheGet[PurchaseQueue](
					keys.CacheKeyPurchaseQueue,
					keys.TypeStrPurchaseQueue,
					true,
					nil,
				),
			},
		},
		transceiveGetFetchFunc(chn, chnFetched),
		nil,
	)

	select {
	case <-chnFetched:
		// channel closed - fetch func called, 'Get' result sent
		// log error if any and return nil
		logger.MaybeErr(err)
		return nil
	default:
		// channel open - fetch func not called, result not sent
		// send result
		return chn.SendExp(purchaseQueue, expires, err)
	}
}

// accept a "fast rep" channel that does not wait for any raw queue deletions
// - if rawpurchasequeue deletion fails, rawpurchasequeue will be out of sync with purchasequeue
// - therefore, do not cache unless deletion succeeds, and thus, do not return from this function until then
//
// even if post-rep fails, the result as well as its expiry will be valid
func transceiveGetFetchFunc(
	chnRep expirable.ChanResult[PurchaseQueue],
	chnRepDone chan<- struct{},
) fetch.Fetch[PurchaseQueue] {
	return func(x cache.Context) (
		purchaseQueuePtr PurchaseQueue,
		expires time.Time,
		postFetch *postfetch.Params,
		err error,
	) {
		var purchaseQueue PurchaseQueue
		var removed []remotedb.CodeAndLocationId
		purchaseQueue, removed, expires, err = transceiveGetFetchFuncInner(x)

		// close the 'Done' channel and send the result to the 'Rep' channel
		close(chnRepDone) // closing results in non-blocking receive
		if err != nil {
			go chnRep.SendErr(err)
			return nil, expires, nil, nil
		} else {
			go chnRep.SendExpOk(purchaseQueue, expires)
		}

		// even if we fail after this point, the sent queue+expires are valid

		// make deletions if necessary
		if len(removed) > 0 {
			err = remotedb.DelPurchases(
				x.Background(), // never cancel this
				removed...,
			)
			if err != nil {
				return nil, expires, nil, err
			}
		}

		// finally, return alongside cache commands
		postFetch = &postfetch.Params{
			CacheParams: &postfetch.CacheParams{
				Set: postfetch.ServerCacheSetOne(
					keys.CacheKeyPurchaseQueue,
					keys.TypeStrPurchaseQueue,
					&purchaseQueue,
					expires,
				),
			},
		}
		return purchaseQueue, expires, postFetch, nil
	}
}

func transceiveGetFetchFuncInner(
	x cache.Context,
) (
	purchaseQueue PurchaseQueue,
	removed []remotedb.CodeAndLocationId,
	expires time.Time,
	err error,
) {
	// send out a goroutine to get the raw purchase queue
	x, cancel := x.WithCancel()
	defer cancel()
	chnQueue := expirable.NewChanResult[remotedb.RawPurchaseQueue](
		x.Ctx(), 1, 0,
	)
	go expirable.Param1Transceive(chnQueue, x, remotedb.GetRawPurchaseQueue)

	// get the shop contracts
	var shopContracts map[string]contracts.Contract
	shopContracts, expires, err = contracts.GetShopContracts(x)
	if err != nil {
		return nil, nil, expires, err
	}

	// recv the raw purchase queue
	purchaseQueue, expires, err = chnQueue.RecvExpMin(expires)
	if err != nil {
		return nil, nil, expires, err
	}

	// split raw queue into kept and removed and return
	purchaseQueue, removed = newPurchaseQueue(purchaseQueue, shopContracts)
	return purchaseQueue, removed, expires, nil
}
