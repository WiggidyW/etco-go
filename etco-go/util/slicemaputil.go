package util

func KeysToSlice[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func SliceToSet[K comparable](slice []K) map[K]struct{} {
	set := make(map[K]struct{}, len(slice))
	for _, k := range slice {
		set[k] = struct{}{}
	}
	return set
}

func KeysNotIn[K comparable, V any](slice []K, m map[K]V) []K {
	notIn := make([]K, 0, len(slice))
	for _, k := range slice {
		if _, ok := m[k]; !ok {
			notIn = append(notIn, k)
		}
	}
	return notIn
}
