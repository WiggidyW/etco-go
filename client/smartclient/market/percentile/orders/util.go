package orders

func toPtrSlice[T any](slice []T) []*T {
	ptrSlice := make([]*T, 0, len(slice))
	for i := range slice {
		ptrSlice[i] = &slice[i]
	}
	return ptrSlice
}

// func fromPtrSlice[T any](ptrSlice []*T) []T {
// 	slice := make([]T, 0, len(ptrSlice))
// 	for i := range ptrSlice {
// 		slice[i] = *ptrSlice[i]
// 	}
// 	return slice
// }
