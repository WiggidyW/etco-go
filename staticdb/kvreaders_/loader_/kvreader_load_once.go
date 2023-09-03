package loader_

import (
	"github.com/WiggidyW/etco-go/staticdb/kvreaders_/loader_/loadonceflag_"
)

type LoadOnceKVReader[
	L Loader[R],
	R KVReader[K, V],
	K any,
	V any,
] struct {
	flag     *loadonceflag_.LoadOnceFlag
	kvReader *Container[R] // nil until loaded
	kvLoader L
}

func NewLoadOnceKvReader[
	L Loader[R],
	R KVReader[K, V],
	K any,
	V any,
](kvLoader L) *LoadOnceKVReader[L, R, K, V] {
	return &LoadOnceKVReader[L, R, K, V]{
		flag:     loadonceflag_.NewLoadOnceFlag(),
		kvLoader: kvLoader,
	}
}

func (lor *LoadOnceKVReader[L, R, K, V]) Load() error {
	reader, err := lor.kvLoader.Load()
	if err != nil {
		return err
	}
	lor.kvReader = &Container[R]{Inner: reader}
	lor.flag.LoadFinish()
	return nil
}

func (lor *LoadOnceKVReader[L, R, K, V]) LoadSendErr(chn chan<- error) {
	if err := lor.Load(); err != nil {
		chn <- err
	}
}

// Blocks until Load() has been called and completed
func (lor *LoadOnceKVReader[L, R, K, V]) Get(key K) (V, bool) {
	lor.flag.Check() // ensure reader data is loaded and safe to use
	return lor.kvReader.Inner.Get(key)
}

// Blocks until Load() has been called and completed
func (lor *LoadOnceKVReader[L, R, K, V]) UnsafeGet(key K) V {
	lor.flag.Check() // ensure reader data is loaded and safe to use
	return lor.kvReader.Inner.UnsafeGet(key)
}

func (lor *LoadOnceKVReader[L, R, K, V]) UnsafeGetInner() R {
	lor.flag.Check() // ensure reader data is loaded and safe to use
	return lor.kvReader.Inner
}
