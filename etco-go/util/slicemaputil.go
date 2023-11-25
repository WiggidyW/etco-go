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

func KeysNotIn[
	K comparable,
	VCHECK any,
	VIN any,
](mCheck map[K]VCHECK, mIn map[K]VIN) map[K]VCHECK {
	notIn := make(map[K]VCHECK, len(mCheck))
	for k, v := range mCheck {
		if _, ok := mIn[k]; !ok {
			notIn[k] = v
		}
	}
	return notIn
}
