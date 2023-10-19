package updater

import (
	"context"
	"net/http"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	updatersde "github.com/WiggidyW/etco-go-updater/sde"
	updaterwtc "github.com/WiggidyW/etco-go-updater/webtocore"
)

type UpdaterResult struct {
	SDEModified  bool
	CoreModified bool
}

// TODO: break into smaller functions
func Update(
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	skipSde bool,
) (
	updaterResult UpdaterResult,
	err error,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// fetch the SDE checksum
	chnSendSDEChecksum, chnRecvSDEChecksum :=
		chanresult.NewChanResult[string](ctx, 1, 0).Split()
	go updatersde.TransceiveDownloadChecksum(
		ctx,
		httpClient,
		userAgent,
		skipSde,
		chnSendSDEChecksum,
	)

	// fetch the Web Attributes (checksums and version string)
	chnSendWebAttrs, chnRecvWebAttrs :=
		chanresult.NewChanResult[updaterwtc.WebAttrs](ctx, 1, 0).Split()
	go updaterwtc.TransceiveDownloadWebAttrs(
		ctx,
		bucketClient,
		chnSendWebAttrs,
	)

	// fetch the UpdaterData (contains previous checksums, and capacities)
	syncUpdaterData, err := DownloadSyncUpdaterData(ctx, bucketClient)
	if err != nil {
		return updaterResult, err
	}

	// prepare to receive and handle the sde checksum and web attrs
	chnSendDone, chnRecvDone := chanresult.
		NewChanResult[struct{}](ctx, 2, 0).Split()

	// alias the update functions
	updateCoreIfModified, updateSDEIfModified :=
		func(webAttrs updaterwtc.WebAttrs) (updating bool) {
			return UpdateCoreIfModified(
				ctx,
				bucketClient,
				syncUpdaterData,
				webAttrs,
				chnSendDone,
			)
		},
		func(sdeChecksum string) (updating bool) {
			return UpdateSDEIfModified(
				ctx,
				bucketClient,
				httpClient,
				userAgent,
				syncUpdaterData,
				sdeChecksum,
				skipSde,
				chnSendDone,
			)
		}

	// handle the sde checksum and web attrs
	variant, sdeChecksum, webAttrs, err := chanresult.Recv2(
		chnRecvSDEChecksum,
		chnRecvWebAttrs,
	)
	if err != nil {
		return updaterResult, err
	}
	if variant == chanresult.RECV_1 {
		// sde checksum received first
		updaterResult.SDEModified = updateSDEIfModified(sdeChecksum)
		webAttrs, err = chnRecvWebAttrs.Recv()
		if err != nil {
			return updaterResult, err
		}
		updaterResult.CoreModified = updateCoreIfModified(webAttrs)
	} else {
		// web attrs received first
		updaterResult.CoreModified = updateCoreIfModified(webAttrs)
		sdeChecksum, err = chnRecvSDEChecksum.Recv()
		if err != nil {
			return updaterResult, err
		}
		updaterResult.SDEModified = updateSDEIfModified(sdeChecksum)
	}

	// count the updates in progress
	var numUpdating int = 0
	if updaterResult.SDEModified {
		numUpdating++
	}
	if updaterResult.CoreModified {
		numUpdating++
	}

	// return early if there are no updates to perform
	if numUpdating == 0 {
		return updaterResult, nil
	}

	// wait for the updates to finish
	for numUpdating > 0 {
		_, err := chnRecvDone.Recv()
		if err != nil {
			return updaterResult, err
		}
		numUpdating--
	}

	// write the new updater data
	err = UploadUpdaterData(ctx, bucketClient, syncUpdaterData)
	if err != nil {
		return updaterResult, err
	}

	return updaterResult, nil
}
