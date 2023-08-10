package staticdb

type HashSet[T comparable] map[T]struct{}

func (hs HashSet[T]) Has(k T) bool {
	_, ok := hs[k]
	return ok
}
