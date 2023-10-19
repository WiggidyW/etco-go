package multi

import (
	"context"
	"fmt"

	"github.com/WiggidyW/chanresult"

	"github.com/WiggidyW/etco-go/cache"
	"github.com/WiggidyW/etco-go/client"
	"github.com/WiggidyW/etco-go/client/cachekeys"
	"github.com/WiggidyW/etco-go/client/caching"
	"github.com/WiggidyW/etco-go/logger"
)

// TODO: Evaluate the usage of context.Background vs reusing the context passed in

type StrongMultiAntiCachingClient[
	F caching.MultiAntiCacheableParams, // the inner client params type
	D any, // the inner client response type
	C client.Client[F, D], // the inner client type
] struct {
	Client     C
	antiCaches []*cache.StrongAntiCache
}

func NewStrongMultiAntiCachingClient[
	F caching.MultiAntiCacheableParams,
	D any,
	C client.Client[F, D],
](
	client C,
	antiCaches ...*cache.StrongAntiCache,
) StrongMultiAntiCachingClient[F, D, C] {
	return StrongMultiAntiCachingClient[F, D, C]{
		Client:     client,
		antiCaches: antiCaches,
	}
}

func (smacc StrongMultiAntiCachingClient[F, D, C]) Fetch(
	ctx context.Context,
	params F,
) (*D, error) {
	antiCacheKeys := params.AntiCacheKeys()

	// fatal if lengths don't match
	if len(antiCacheKeys) != len(smacc.antiCaches) {
		logger.Fatal(fmt.Errorf(
			"antiCacheKeys length (%d) != antiCaches length (%d)",
			len(antiCacheKeys),
			len(smacc.antiCaches),
		))
	}

	// start the individual anti-cache threads (one for each anticache)
	chnCtx, chnCancel := context.WithCancel(context.Background())
	chnSend, chnRecv := chanresult.
		NewChanResult[struct{}](chnCtx, 0, 0).Split()
	for i, antiCacheKey := range antiCacheKeys {
		if antiCacheKey == cachekeys.NULL_ANTI_CACHE_KEY {
			continue
		}
		go smacc.fetchOne(ctx, chnSend, antiCacheKey, i)
	}

	// wait for them to finish locking
	if err := chnRecv.RecvAllDiscard(len(antiCacheKeys)); err != nil {
		chnCancel()
		return nil, err
	}

	// fetch
	rep, err := smacc.Client.Fetch(ctx, params)

	// the cached values have been deleted
	// the inner client has either completed its task or failed
	// thus, the locks have no further use.
	// (no reason to block the caller)
	chnCancel()

	// return the fetch result
	if err != nil {
		return nil, err
	} else {
		return rep, nil
	}
}

func (smacc StrongMultiAntiCachingClient[F, D, C]) fetchOne(
	ctx context.Context, // propagated from parent, could cancel the lock
	chnSend chanresult.ChanSendResult[struct{}], // uses separate context
	antiCacheKey string,
	antiCacheIdx int,
) {
	if antiCacheKey == cachekeys.NULL_ANTI_CACHE_KEY {
		return
	}

	// lock the cache
	lock, err := smacc.antiCaches[antiCacheIdx].Lock(ctx, antiCacheKey)
	if err != nil { // lock acquisition failed
		// try sending the error, or log it if cancelled
		if ctxErr := chnSend.SendErr(err); ctxErr != nil {
			logger.Err(err)
		}
		// return now, since we don't need to unlock
		return
	}

	// // cache delete

	err = smacc.antiCaches[antiCacheIdx].Del(antiCacheKey, lock)

	if err != nil { // failed
		// unlock in a goroutine
		go func() {
			logger.Err(smacc.antiCaches[antiCacheIdx].Unlock(lock))
		}()
		// try sending the error, or log it if cancelled
		if ctxErr := chnSend.SendErr(err); ctxErr != nil {
			logger.Err(err)
		}

	} else { // succeeded
		// try to send Ok
		if ctxErr := chnSend.SendOk(struct{}{}); ctxErr == nil {
			// if not cancelled, then wait for cancellation
			<-ctx.Done()
		}
		// unlock synchronously
		logger.Err(smacc.antiCaches[antiCacheIdx].Unlock(lock))
	}
	//
}
