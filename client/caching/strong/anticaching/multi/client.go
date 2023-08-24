package multi

import (
	"context"
	"fmt"

	"github.com/WiggidyW/eve-trading-co-go/cache"
	"github.com/WiggidyW/eve-trading-co-go/client"
	"github.com/WiggidyW/eve-trading-co-go/client/caching"
	"github.com/WiggidyW/eve-trading-co-go/logger"
	"github.com/WiggidyW/eve-trading-co-go/util"
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
	chnSend, chnRecv := util.NewChanResult[struct{}](chnCtx).Split()
	for i, antiCacheKey := range antiCacheKeys {
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
	chnSend util.ChanSendResult[struct{}], // uses separate context
	cKey string,
	cIdx int,
) {
	// lock the cache
	lock, err := smacc.antiCaches[cIdx].Lock(ctx, cKey)
	if err != nil { // lock acquisition failed
		// try sending the error, or log it if cancelled
		if ctxErr := chnSend.SendErr(err); ctxErr != nil {
			logger.Err(err)
		}
		// return now, since we don't need to unlock
		return
	}

	// // cache delete

	if err := smacc.antiCaches[cIdx].Del(cKey, lock); err != nil { // failed
		// unlock in a goroutine
		go func() { logger.Err(smacc.antiCaches[cIdx].Unlock(lock)) }()
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
		logger.Err(smacc.antiCaches[cIdx].Unlock(lock))
	}
	//
}
