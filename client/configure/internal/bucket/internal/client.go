package bucket

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
	_storageClient    *storage.Client
	storageClientOpts []option.ClientOption
	bucketName        string
	_bucketHandle     *storage.BucketHandle
	mu                *sync.Mutex
	// bufPool *cache.BufferPool // same one used by local cache
}

// Gets the storage client (sets it if it's nil)
func (asc *BucketClient) storageClient() (*storage.Client, error) {
	if asc._storageClient == nil {
		// create the client
		ctx := context.Background()
		var err error
		if asc._storageClient, err = storage.NewClient(
			ctx,
			asc.storageClientOpts...,
		); err != nil {
			return nil, err
		}
	}
	return asc._storageClient, nil
}

// Gets the bucket handle (sets it if it's nil)
func (asc *BucketClient) bucketHandle() (*storage.BucketHandle, error) {
	if asc._bucketHandle == nil {
		// lock to prevent multiple handles/clients from being created
		asc.mu.Lock()
		defer asc.mu.Unlock()

		// check again in case another handle was created while waiting
		if asc._bucketHandle != nil {
			return asc._bucketHandle, nil
		}

		// get the storage client and create the bucket handle
		storageClient, err := asc.storageClient()
		if err != nil {
			return nil, err
		}
		asc._bucketHandle = storageClient.Bucket(asc.bucketName)
	}
	return asc._bucketHandle, nil
}

// Gets the object handle from the provided key
func (asc *BucketClient) objHandle(
	key string,
) (*storage.ObjectHandle, error) {
	bucketHandle, err := asc.bucketHandle()
	if err != nil {
		return nil, err
	}
	return bucketHandle.Object(key), nil
}

// Reads an authHashSet from the provided key
func (asc *BucketClient) Read(
	ctx context.Context,
	key string,
	valPtr any,
) (bool, error) {
	// get the object handle
	objHandle, err := asc.objHandle(key)
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

// Writes the provided auth list to the provided key
func (asc *BucketClient) Write(
	ctx context.Context,
	key string,
	val any,
) error {
	// get the object handle
	objHandle, err := asc.objHandle(key)
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
