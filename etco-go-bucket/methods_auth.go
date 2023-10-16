package etcogobucket

import (
	"context"
)

func (bc *BucketClient) ReadAttrsAuthHashSet(
	ctx context.Context,
	key string, // careful with this matching constants
) (*Attrs, error) {
	return bc.readAttrs(
		ctx,
		AUTH,
		key,
	)
}

func (bc *BucketClient) ReadAuthHashSet(
	ctx context.Context,
	key string, // careful with this matching constants
) (v AuthHashSet, err error) {
	_, err = read(
		bc,
		ctx,
		AUTH,
		key,
		&v,
	)
	return v, err
}

func (bc *BucketClient) WriteAuthHashSet(
	ctx context.Context,
	v AuthHashSet,
	key string, // careful with this matching constants
) error {
	return write(
		bc,
		ctx,
		AUTH,
		key,
		v,
	)
}
