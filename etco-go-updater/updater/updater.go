package updater

import (
	"context"
	"net/http"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"
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
		chanresult.NewChanResult[bool](ctx, 1, 0).Split()
	go transceiveUpdateSDEIfModified(
		ctx,
		bucketClient,
		httpClient,
		userAgent,
		skipSde,
		chnSendSDEModified,
	)

	chnSendCoreModified, chnRecvCoreModified :=
		chanresult.NewChanResult[bool](ctx, 1, 0).Split()
	go transceiveUpdateCoreIfModified(
		ctx,
		bucketClient,
		skipCore,
		chnSendCoreModified,
	)

	updaterResult.SDEModified, err = chnRecvSDEModified.Recv()
	if err != nil {
		return updaterResult, err
	}

	updaterResult.CoreModified, err = chnRecvCoreModified.Recv()
	if err != nil {
		return updaterResult, err
	}

	return updaterResult, nil
}
