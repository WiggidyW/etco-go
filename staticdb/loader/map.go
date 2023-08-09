package loader

type LoadOnceKVReaderGobFSMap[K comparable, V any] struct {
	*LoadOnceKVReader[
		GobFSMapLoader[K, V],
		MapWrapper[K, V],
		K,
		V,
	]
}

func NewLoadOnceKVReaderGobFSMap[K comparable, V any](
	path string,
	capacity int,
) LoadOnceKVReaderGobFSMap[K, V] {
	return LoadOnceKVReaderGobFSMap[K, V]{
		NewLoadOnceKvReader[
			GobFSMapLoader[K, V],
			MapWrapper[K, V],
			K,
			V,
		](
			NewGobFSMapLoader[K, V](path, capacity),
		),
	}
}

type GobFSMapLoader[K comparable, V any] struct {
	path     string
	capacity int
}

func NewGobFSMapLoader[K comparable, V any](
	path string,
	capacity int,
) GobFSMapLoader[K, V] {
	return GobFSMapLoader[K, V]{path, capacity}
}

func (gfml GobFSMapLoader[K, V]) Load() (MapWrapper[K, V], error) {
	m := make(map[K]V, gfml.capacity)
	err := gobFsLoad(&m, gfml.path)
	kvReader := MapWrapper[K, V]{m}
	if err != nil {
		return kvReader, err
	} else {
		return kvReader, nil
	}
}

type MapWrapper[K comparable, V any] struct {
	m map[K]V
}

func (mw MapWrapper[K, V]) Get(k K) (V, bool) {
	v, ok := mw.m[k]
	return v, ok
}

func (mw MapWrapper[K, V]) UnsafeGet(k K) V {
	return mw.m[k]
}
