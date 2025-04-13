//go:build malloc_cgo

package memory

// #include <stdlib.h>
import "C"

import "unsafe"

// pointer is ...
type pointer struct {
	Pointer uintptr
}

// Alloc is ...
func Alloc[T any](n int) (pointer, []T, error) {
	ptr := uintptr(C.malloc(C.size_t(n * int(unsafe.Sizeof(*(new(T)))))))
	return pointer{Pointer: ptr}, unsafe.Slice((*T)(unsafe.Pointer(ptr)), n), nil
}

// Free is ...
func Free(p pointer) error {
	C.free(unsafe.Pointer(p.Pointer))
	return nil
}
