package loader_

type LoadOnceKVReaderGobFSSlice[V any] struct {
	*LoadOnceKVReader[
		GobFSSliceLoader[V],
		SliceWrapper[V],
		int,
		V,
	]
}

func NewLoadOnceKVReaderGobFSSlice[V any](
	path string,
	capacity int,
) LoadOnceKVReaderGobFSSlice[V] {
	return LoadOnceKVReaderGobFSSlice[V]{
		NewLoadOnceKvReader[
			GobFSSliceLoader[V],
			SliceWrapper[V],
			int,
			V,
		](
			NewGobFSSliceLoader[V](path, capacity),
		),
	}
}

type GobFSSliceLoader[V any] struct {
	path     string
	capacity int
}

func NewGobFSSliceLoader[V any](
	path string,
	capacity int,
) GobFSSliceLoader[V] {
	return GobFSSliceLoader[V]{path, capacity}
}

func (gfsl GobFSSliceLoader[V]) Load() (SliceWrapper[V], error) {
	s := make([]V, 0, gfsl.capacity)
	err := gobFsLoad(&s, gfsl.path)
	kvReader := SliceWrapper[V]{s}
	if err != nil {
		return kvReader, err
	} else {
		return kvReader, nil
	}
}

type SliceWrapper[V any] struct {
	s []V
}

func (sw SliceWrapper[V]) Get(k int) (V, bool) {
	if k < 0 || k >= len(sw.s) {
		var null V
		return null, false
	}
	return sw.s[k], true
}

func (sw SliceWrapper[V]) UnsafeGet(k int) V {
	return sw.s[k]
}

func (sw SliceWrapper[V]) UnsafeGetInner() []V {
	return sw.s
}
