package reader

import "github.com/WiggidyW/weve-esi/client/cachekeys"

type BucketReaderParams struct {
	ObjectName string
} // object name

func (p BucketReaderParams) CacheKey() string {
	return cachekeys.BucketReaderCacheKey(p.ObjectName)
}
