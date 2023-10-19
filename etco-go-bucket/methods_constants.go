package etcogobucket

import (
	"context"
)

func (bc *BucketClient) ReadConstantsData(
	ctx context.Context,
) (v ConstantsData, err error) {
	_, err = read(
		bc,
		ctx,
		BUILD,
		OBJNAME_CONSTANTS_DATA,
		&v,
	)
	return
}

func (bc *BucketClient) WriteConstantsData(
	ctx context.Context,
	v ConstantsData,
) error {
	return write(
		bc,
		ctx,
		BUILD,
		OBJNAME_CONSTANTS_DATA,
		v,
	)
}
