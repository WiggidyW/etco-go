package util

import (
	"reflect"
	"unsafe"
)

func ToPtrSlice[T any](slice []T) []*T {
	ptrSlice := make([]*T, 0, len(slice))
	for i := range slice {
		ptrSlice[i] = &slice[i]
	}
	return ptrSlice
}

// the 'B' type must either be an alias of the 'A' type or a struct with the same fields
func UnsafeSliceToSlice[A any, B any](in []A) []B {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&in))
	return *(*[]B)(unsafe.Pointer(&header))
}
