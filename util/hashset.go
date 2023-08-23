package util

type HashSet[K comparable] interface {
	Has(k K) bool
}

type MapHashSet[K comparable, V any] map[K]V

func (mhs MapHashSet[K, V]) Has(k K) bool {
	_, ok := mhs[k]
	return ok
}
