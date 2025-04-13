//go:build (darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris) && malloc_syscall

package memory

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

// pointer is ...
type pointer struct {
	Pointer uintptr
	Cap     uintptr
}

func Alloc[T any](n int) (pointer, []T, error) {
	size := uintptr(n * int(unsafe.Sizeof(*(new(T)))))
	ptr, _, errno := unix.Syscall6(
		unix.SYS_MMAP,
		0,
		size,
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_ANONYMOUS|unix.MAP_PRIVATE,
		0,
		0,
	)
	if errno != 0 {
		return pointer{}, nil, fmt.Errorf("failed to make MMAP allocator: %w", errno)
	}

	return pointer{Pointer: ptr, Cap: size}, unsafe.Slice((*T)(unsafe.Pointer(ptr)), n), nil
}

func Free(p pointer) error {
	_, _, errno := unix.Syscall(
		unix.SYS_MUNMAP,
		p.Pointer,
		p.Cap,
		0,
	)
	if errno != 0 {
		return fmt.Errorf("failed to make MUNMAP allocator: %w", errno)
	}

	return nil
}
