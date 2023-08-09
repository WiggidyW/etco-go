package cache

import (
	// "bytes"
	"sync"
)

type BufferPool struct {
	pool     sync.Pool
	capacity int
}

func NewBufferPool(capacity int) *BufferPool {
	return &BufferPool{
		pool:     sync.Pool{},
		capacity: capacity,
	}
}

func (dp *BufferPool) Get() *[]byte {
	if buf := dp.pool.Get(); buf != nil {
		return buf.(*[]byte)
	}
	newBuf := make([]byte, 0, dp.capacity)
	return &newBuf
}

func (dp *BufferPool) Put(buf *[]byte) {
	if len(*buf) > dp.capacity {
		dp.capacity = len(*buf) // increase capacity
	}
	*buf = (*buf)[:0]
	dp.pool.Put(buf)
}

// type BufferPool struct {
// 	pool     sync.Pool
// 	capacity int
// }

// func NewBufferPool(capacity int) *BufferPool {
// 	return &BufferPool{
// 		pool:     sync.Pool{},
// 		capacity: capacity,
// 	}
// }

// func (bp *BufferPool) get() *bytes.Buffer {
// 	if buf := bp.pool.Get(); buf != nil {
// 		return buf.(*bytes.Buffer)
// 	}
// 	return bytes.NewBuffer(make([]byte, 0, bp.capacity))
// }

// func (bp *BufferPool) put(buf *bytes.Buffer) {
// 	buf.Reset()
// 	bp.pool.Put(buf)
// }

// type DstPool struct {
// 	pool     sync.Pool
// 	capacity int
// }

// func NewDstPool(capacity int) *DstPool {
// 	return &DstPool{
// 		pool:     sync.Pool{},
// 		capacity: capacity,
// 	}
// }

// func (dp *DstPool) get() *[]byte {
// 	if dst := dp.pool.Get(); dst != nil {
// 		return dst.(*[]byte)
// 	}
// 	newDst := make([]byte, 0, dp.capacity)
// 	return &newDst
// }

// func (dp *DstPool) put(dst *[]byte) {
// 	*dst = (*dst)[:0]
// 	dp.pool.Put(dst)
// }
