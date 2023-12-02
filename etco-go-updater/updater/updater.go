package updater

import (
	"context"
	"net/http"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	updaterwtc "github.com/WiggidyW/etco-go-updater/webtocore"
)

type UpdaterResult struct {
	SDEModified  bool
	CoreModified bool
}

func Update(
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	skipSde bool,
	skipCore bool,
) (
	updaterResult UpdaterResult,
	err error,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chnSendSDEModified, chnRecvSDEModified :=
		chanresult.NewChanResult[*b.SDEBucketData](ctx, 1, 0).Split()
	go transceiveUpdateSDEIfModified(
		ctx,
		bucketClient,
		httpClient,
		userAgent,
		skipSde,
		chnSendSDEModified,
	)

	chnSendCoreModified, chnRecvCoreModified :=
		chanresult.NewChanResult[*updaterwtc.WebAttrs](ctx, 1, 0).Split()
	go transceiveCoreModified(
		ctx,
		bucketClient,
		skipCore,
		chnSendCoreModified,
	)

	var sdeBucketDataPtr *b.SDEBucketData
	sdeBucketDataPtr, err = chnRecvSDEModified.Recv()
	if err != nil {
		updaterResult.SDEModified = false
		return updaterResult, err
	} else {
		updaterResult.SDEModified = sdeBucketDataPtr != nil
	}

	var coreWebAttrs *updaterwtc.WebAttrs
	coreWebAttrs, err = chnRecvCoreModified.Recv()
	if err != nil || coreWebAttrs == nil {
		updaterResult.CoreModified = false
		return updaterResult, err
	}

	var sdeBucketData b.SDEBucketData
	if sdeBucketDataPtr != nil {
		sdeBucketData = *sdeBucketDataPtr
	} else {
		sdeBucketData, err = bucketClient.ReadSDEData(
			ctx, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		)
		if err != nil {
			return updaterResult, err
		}
	}

	err = UpdateCore(ctx, bucketClient, *coreWebAttrs, sdeBucketData.Systems)
	if err == nil {
		updaterResult.CoreModified = true
	}
	return updaterResult, err
}
