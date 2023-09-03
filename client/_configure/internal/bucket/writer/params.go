package writer

import "github.com/WiggidyW/etco-go/client/cachekeys"

type BucketWriterParams[D any] struct {
	ObjectName string // object name (domain key + access type)
	Val        D
}

func (p BucketWriterParams[D]) AntiCacheKey() string {
	return cachekeys.BucketReaderCacheKey(p.ObjectName)
}
