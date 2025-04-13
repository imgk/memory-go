//go:build windows && malloc_syscall

package memory

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// DefaultAllocator is ...
var DefaultAllocator = NewAllocator()

// pointer is ...
type pointer struct {
	Pointer uintptr
}

func Alloc[T any](n int) (pointer, []T, error) {
	size := n * int(unsafe.Sizeof(*(new(T))))

	ptr, err := DefaultAllocator.Alloc(size)
	if err != nil {
		return ptr, nil, err
	}

	return ptr, unsafe.Slice((*T)(unsafe.Pointer(ptr.Pointer)), n), nil
}

func Free(p pointer) error {
	return DefaultAllocator.Free(p)
}

type Allocator struct {
	*windows.LazyDLL
	AllocProc *windows.LazyProc
	FreeProc  *windows.LazyProc
}

func NewAllocator() Allocator {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	virtualAlloc := kernel32.NewProc("VirtualAlloc")
	virtualFree := kernel32.NewProc("VirtualFree")
	return Allocator{LazyDLL: kernel32, AllocProc: virtualAlloc, FreeProc: virtualFree}
}

func (alloc Allocator) Alloc(n int) (pointer, error) {
	ptr, _, err := alloc.AllocProc.Call(
		0,
		uintptr(n),
		uintptr(0x1000), // MEM_COMMIT
		uintptr(0x04),   // PAGE_READWRITE
	)
	if ptr == 0 {
		return pointer{}, fmt.Errorf("failed to make VirtualAlloc allocator: %w", err)
	}

	return pointer{Pointer: ptr}, nil
}

func (alloc Allocator) Free(p pointer) error {
	r1, _, err := alloc.FreeProc.Call(
		p.Pointer,
		0,
		uintptr(0x8000), // MEM_RELEASE
	)
	if r1 == 0 {
		return fmt.Errorf("failed to make VirtualFree allocator: %w", err)
	}

	return nil
}
