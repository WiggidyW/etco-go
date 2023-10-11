package etcogobucket

import (
	"context"
	"encoding/gob"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Basal Bucket client
// singleton
// has methods to read and write auth lists
type BucketClient struct { // single one for the program
	_client       *storage.Client
	_bucketHandle *storage.BucketHandle
	clientOpts    []option.ClientOption
	mu            *sync.Mutex
	// bucketName        string
}

func NewBucketClient(creds []byte) *BucketClient {
	return &BucketClient{
		// _client:    nil,
		// _bucketHandle: nil,
		clientOpts: []option.ClientOption{
			option.WithCredentialsJSON(creds),
		},
		mu: &sync.Mutex{},
	}
}

// Gets the storage client (sets it if it's nil)
func (bc *BucketClient) innerClient() (*storage.Client, error) {
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

// Gets the bucket handle (sets it if it's nil)
func (bc *BucketClient) bucketHandle() (*storage.BucketHandle, error) {
	if bc._bucketHandle == nil {
		// lock to prevent multiple handles/clients from being created
		bc.mu.Lock()
		defer bc.mu.Unlock()

		// check again in case another handle was created while waiting
		if bc._bucketHandle != nil {
			return bc._bucketHandle, nil
		}

		// get the storage client and create the bucket handle
		storageClient, err := bc.innerClient()
		if err != nil {
			return nil, err
		}
		bc._bucketHandle = storageClient.Bucket(
			BUCKET_NAME, // used to be bc.bucketName
		)
	}
	return bc._bucketHandle, nil
}

// Gets the object handle from the provided key
func (bc *BucketClient) objHandle(
	objName string,
) (*storage.ObjectHandle, error) {
	bucketHandle, err := bc.bucketHandle()
	if err != nil {
		return nil, err
	}
	return bucketHandle.Object(objName), nil
}

func (bc *BucketClient) readAttrs(
	ctx context.Context,
	objName string,
) (*Attrs, error) {
	objHandle, err := bc.objHandle(objName)
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
	objName string,
	valPtr *D,
) (bool, error) {
	// get the object handle
	objHandle, err := bc.objHandle(objName)
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
	objName string,
	val D,
) error {
	// get the object handle
	objHandle, err := bc.objHandle(objName)
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
