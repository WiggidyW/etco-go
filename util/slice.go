package util

func ToPtrSlice[T any](slice []T) []*T {
	ptrSlice := make([]*T, 0, len(slice))
	for i := range slice {
		ptrSlice[i] = &slice[i]
	}
	return ptrSlice
}
