package updater

import (
	"context"
	"sync"

	b "github.com/WiggidyW/etco-go-bucket"
)

type SyncUpdaterData struct {
	rwLock *sync.RWMutex
	b.UpdaterData
}

func (sud *SyncUpdaterData) RLock() {
	sud.rwLock.RLock()
}

func (sud *SyncUpdaterData) RUnlock() {
	sud.rwLock.RUnlock()
}

func (sud *SyncUpdaterData) WLock() {
	sud.rwLock.Lock()
}

func (sud *SyncUpdaterData) WUnlock() {
	sud.rwLock.Unlock()
}

func DownloadSyncUpdaterData(
	ctx context.Context,
	bucketClient *b.BucketClient,
) (
	syncUpdaterData *SyncUpdaterData,
	err error,
) {
	updaterData, err := bucketClient.ReadUpdaterData(ctx)
	if err != nil {
		return nil, err
	}
	return &SyncUpdaterData{
		rwLock:      &sync.RWMutex{},
		UpdaterData: updaterData,
	}, nil
}

func UploadUpdaterData(
	ctx context.Context,
	bucketClient *b.BucketClient,
	syncUpdaterData *SyncUpdaterData,
) error {
	syncUpdaterData.RLock()
	updaterData := syncUpdaterData.UpdaterData
	syncUpdaterData.RUnlock()
	return bucketClient.WriteUpdaterData(ctx, updaterData)
}
