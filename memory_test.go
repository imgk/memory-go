package memory

import (
	"testing"
	"unsafe"
)

func TestAllocator(t *testing.T) {
	Pool := NewAllocator()
	for k, v := range map[int]int{
		1:    1,
		2:    2,
		3:    4,
		4:    4,
		5:    8,
		6:    8,
		7:    8,
		8:    8,
		9:    16,
		31:   32,
		63:   64,
		65:   128,
		127:  128,
		129:  256,
		257:  512,
		513:  1024,
		1025: 2048,
		2049: 4096,
	} {
		p := Pool.Get(k)
		if len(*p.Pointer) != v {
			t.Errorf("Pool.Get error, size: %v, length: %v", k, v)
		}
		Pool.Put(p)
	}
}

func TestAlloc(t *testing.T) {
	p, b, _ := Alloc[byte](1024)
	buf := *p.Pointer
	if uintptr(unsafe.Pointer(&buf[0])) != *(*uintptr)(unsafe.Pointer(&b)) {
		t.Errorf("Alloc error")
	}
}
