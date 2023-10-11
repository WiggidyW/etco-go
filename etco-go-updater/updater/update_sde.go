package updater

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/WiggidyW/chanresult"
	b "github.com/WiggidyW/etco-go-bucket"

	updatersde "github.com/WiggidyW/etco-go-updater/sde"
)

func UpdateSDEIfModified(
	ctx context.Context,
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	syncUpdaterData *SyncUpdaterData,
	sdeChecksum string,
	chnSendDone chanresult.ChanSendResult[struct{}],
) (updating bool) {
	syncUpdaterData.RLock()
	updating = sdeChecksum != syncUpdaterData.CHECKSUM_SDE
	syncUpdaterData.RUnlock()

	if !updating {
		return false
	}

	go TransceiveUpdateSDE(
		ctx,
		bucketClient,
		httpClient,
		userAgent,
		syncUpdaterData,
		sdeChecksum,
		chnSendDone,
	)
	return true
}

func TransceiveUpdateSDE(
	ctx context.Context,
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	syncUpdaterData *SyncUpdaterData,
	sdeChecksum string,
	chnSendDone chanresult.ChanSendResult[struct{}],
) error {
	err := UpdateSDE(
		ctx,
		bucketClient,
		httpClient,
		userAgent,
		syncUpdaterData,
		sdeChecksum,
	)
	if err != nil {
		return chnSendDone.SendErr(err)
	} else {
		return chnSendDone.SendOk(struct{}{})
	}
}

func UpdateSDE(
	ctx context.Context,
	bucketClient *b.BucketClient,
	httpClient *http.Client,
	userAgent string,
	syncUpdaterData *SyncUpdaterData,
	sdeChecksum string,
) error {
	// create a temporary directory
	tempDir, err := os.MkdirTemp("", "etco-go-updater-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	sdeBucketData, err := updatersde.DownloadAndConvert(
		ctx,
		httpClient,
		userAgent,
		tempDir,
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		updateSDESyncUpdaterData(
			syncUpdaterData,
			sdeChecksum,
			sdeBucketData,
		)
		wg.Done()
	}()

	// avoid writing anything if an error occured somewhere
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// write the data to the bucket
	err = bucketClient.WriteSDEData(ctx, sdeBucketData)
	if err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func updateSDESyncUpdaterData(
	syncUpdaterData *SyncUpdaterData,
	sdeChecksum string,
	sdeBucketData b.SDEBucketData,
) {
	syncUpdaterData.WLock()
	defer syncUpdaterData.WUnlock()

	syncUpdaterData.CHECKSUM_SDE = sdeChecksum

	syncUpdaterData.CAPACITY_SDE_CATEGORIES =
		len(sdeBucketData.Categories)
	syncUpdaterData.CAPACITY_SDE_GROUPS =
		len(sdeBucketData.Groups)
	syncUpdaterData.CAPACITY_SDE_MARKET_GROUPS =
		len(sdeBucketData.MarketGroups)
	syncUpdaterData.CAPACITY_SDE_NAME_TO_TYPE_ID =
		len(sdeBucketData.NameToTypeId)
	syncUpdaterData.CAPACITY_SDE_REGIONS =
		len(sdeBucketData.Regions)
	syncUpdaterData.CAPACITY_SDE_SYSTEMS =
		len(sdeBucketData.Systems)
	syncUpdaterData.CAPACITY_SDE_STATIONS =
		len(sdeBucketData.Stations)
	syncUpdaterData.CAPACITY_SDE_TYPE_DATA_MAP =
		len(sdeBucketData.TypeDataMap)
	syncUpdaterData.CAPACITY_SDE_TYPE_VOLUMES =
		len(sdeBucketData.TypeVolumes)
}
