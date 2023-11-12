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

func (bp *BufferPool) Cap() int { return bp.capacity }

func (bp *BufferPool) Get() *[]byte {
	if buf := bp.pool.Get(); buf != nil {
		return buf.(*[]byte)
	}
	newBuf := make([]byte, 0, bp.capacity)
	return &newBuf
}

func (bp *BufferPool) Put(buf *[]byte) {
	if len(*buf) > bp.capacity {
		bp.capacity = len(*buf) // increase capacity
	}
	*buf = (*buf)[:0]
	bp.pool.Put(buf)
}
