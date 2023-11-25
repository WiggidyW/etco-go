package keys

import (
	"fmt"
	"hash"
	"hash/fnv"
	"sync"
)

var (
	hasherPool sync.Pool = sync.Pool{New: func() any { return fnv.New64() }}
)

func censor(s string) string {
	hasher := hasherPool.Get().(hash.Hash64)
	hasher.Write([]byte(s))
	censored := fmt.Sprintf("%016x", hasher.Sum64())
	hasher.Reset()
	hasherPool.Put(hasher)
	return censored
}
