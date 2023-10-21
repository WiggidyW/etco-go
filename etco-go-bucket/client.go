package etcogobucket

import (
	"context"
	"encoding/gob"
	"fmt"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Basal Bucket client
// singleton
// has methods to read and write auth lists
type BucketClient struct { // single one for the program
	_client            *storage.Client
	_webBucketHandle   *storage.BucketHandle
	_authBucketHandle  *storage.BucketHandle
	_buildBucketHandle *storage.BucketHandle
	clientOpts         []option.ClientOption
	mu                 *sync.Mutex
	nameSpace          string
	// bucketName        string
}

func NewBucketClient(nameSpace string, creds []byte) *BucketClient {
	return &BucketClient{
		// _client:    nil,
		// _bucketHandle: nil,
		clientOpts: []option.ClientOption{
			option.WithCredentialsJSON(creds),
		},
		mu:        &sync.Mutex{},
		nameSpace: nameSpace,
	}
}

// Gets the storage client (sets it if it's nil)
func (bc *BucketClient) _getClient() (*storage.Client, error) {
	if bc._client == nil {
		// create the client
		ctx := context.Background()
		var err error
		if bc._client, err = storage.NewClient(
			ctx,
			bc.clientOpts...,
		); err != nil {
			return nil, err
		}
	}
	return bc._client, nil
}

func (bc *BucketClient) _getBucketHandleMaybeNil(kind BucketKind) *storage.BucketHandle {
	switch kind {
	case WEB:
		return bc._webBucketHandle
	case AUTH:
		return bc._authBucketHandle
	case BUILD:
		return bc._buildBucketHandle
	default:
		panic("invalid bucket kind")
	}
}

// Initializes the bucket handles
func (bc *BucketClient) _initBucketHandles() error {
	// check if handles are created
	if bc._webBucketHandle != nil &&
		bc._authBucketHandle != nil &&
		bc._buildBucketHandle != nil {
		return nil
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()

	// check again in case handles were created while waiting
	if bc._webBucketHandle != nil &&
		bc._authBucketHandle != nil &&
		bc._buildBucketHandle != nil {
		return nil
	}

	// get the storage client and create the bucket handles
	storageClient, err := bc._getClient()
	if err != nil {
		return err
	}

	bc._webBucketHandle = storageClient.Bucket(fmt.Sprintf(
		"%s-%s",
		bc.nameSpace,
		WEB_BUCKET_NAME,
	))
	bc._authBucketHandle = storageClient.Bucket(fmt.Sprintf(
		"%s-%s",
		bc.nameSpace,
		AUTH_BUCKET_NAME,
	))
	bc._buildBucketHandle = storageClient.Bucket(fmt.Sprintf(
		"%s-%s",
		bc.nameSpace,
		BUILD_BUCKET_NAME,
	))

	return nil
}

func (bc *BucketClient) getBucketHandle(
	kind BucketKind,
) (
	handle *storage.BucketHandle,
	err error,
) {
	handle = bc._getBucketHandleMaybeNil(kind)
	if handle == nil {
		err = bc._initBucketHandles()
		if err != nil {
			return nil, err
		}
		handle = bc._getBucketHandleMaybeNil(kind)
	}
	return handle, nil
}

// Gets the object handle from the provided key
func (bc *BucketClient) objHandle(
	kind BucketKind,
	objName string,
) (*storage.ObjectHandle, error) {
	bucketHandle, err := bc.getBucketHandle(kind)
	if err != nil {
		return nil, err
	}
	return bucketHandle.Object(objName), nil
}

func (bc *BucketClient) readAttrs(
	ctx context.Context,
	kind BucketKind,
	objName string,
) (*Attrs, error) {
	objHandle, err := bc.objHandle(kind, objName)
	if err != nil {
		return nil, err
	}

	attrs, err := objHandle.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return attrs, nil
}

// Reads a val from the provided key
func read[D any](
	bc *BucketClient,
	ctx context.Context,
	kind BucketKind,
	objName string,
	valPtr *D,
) (bool, error) {
	// get the object handle
	objHandle, err := bc.objHandle(kind, objName)
	if err != nil {
		return false, err
	}

	// create the reader
	reader, err := objHandle.NewReader(ctx)
	if err != nil {
		// return empty if the object doesn't exist
		if err == storage.ErrObjectNotExist {
			return false, nil
		} else {
			return false, err
		}
	}
	defer reader.Close()

	// deserialize using Gob
	decoder := gob.NewDecoder(reader)
	if err = decoder.Decode(valPtr); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// Writes the provided val to the provided key
func write[D any](
	bc *BucketClient,
	ctx context.Context,
	kind BucketKind,
	objName string,
	val D,
) error {
	// get the object handle
	objHandle, err := bc.objHandle(kind, objName)
	if err != nil {
		return err
	}

	// create the writer
	writer := objHandle.NewWriter(ctx)
	defer writer.Close()

	// serialize using Gob
	encoder := gob.NewEncoder(writer)
	if err = encoder.Encode(val); err != nil {
		return err
	} else {
		return nil
	}
}
