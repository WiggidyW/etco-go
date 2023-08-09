package staticdb

type KVReader[K any, V any] interface {
	Get(K) (V, bool)
	UnsafeGet(K) V
}
