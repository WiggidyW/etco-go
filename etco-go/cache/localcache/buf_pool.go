package localcache

import (
	"sync"
)

type BufferPool struct {
	pool     sync.Pool
	capacity int
}

func newBufferPool(capacity int) *BufferPool {
	return &BufferPool{
		pool:     sync.Pool{},
		capacity: capacity,
	}
}

func (dp *BufferPool) get() *[]byte {
	if buf := dp.pool.Get(); buf != nil {
		return buf.(*[]byte)
	}
	newBuf := make([]byte, 0, dp.capacity)
	return &newBuf
}

func (dp *BufferPool) put(buf *[]byte) {
	if len(*buf) > dp.capacity {
		dp.capacity = len(*buf) // increase capacity
	}
	*buf = (*buf)[:0]
	dp.pool.Put(buf)
}
