package etcogobucket

import (
	"context"
)

func (bc *BucketClient) ReadUpdaterData(
	ctx context.Context,
) (v UpdaterData, err error) {
	_, err = read(
		bc,
		ctx,
		BUILD,
		OBJNAME_UPDATER_DATA,
		&v,
	)
	return
}

func (bc *BucketClient) WriteUpdaterData(
	ctx context.Context,
	v UpdaterData,
) error {
	return write(
		bc,
		ctx,
		BUILD,
		OBJNAME_UPDATER_DATA,
		v,
	)
}
