package main

import "testing"
import (
	"fmt"
	"runtime"
	"weak"
)

type Blob []byte

func (b Blob) String() string {
	return fmt.Sprintf("Blob(%d KB)", len(b)/1024)
}

// newBlob returns a new Blob of the given size in KB.
func newBlob(size int) *Blob {
	b := make([]byte, size*1024)
	for i := range size {
		b[i] = byte(i) % 255
	}
	return (*Blob)(&b)
}

// heapDelta returns the delta in KB between
// the current heap size and the previous heap size.
func heapDelta(prev uint64) uint64 {
	cur := getAlloc()
	if cur < prev {
		return 0
	}
	return cur - prev
}

// getAlloc returns the current heap size in KB.
func getAlloc() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc / 1024
}

func Test_WeakPointer(t *testing.T) {
	heapSize := getAlloc()
	wb := weak.Make(newBlob(1000))

	fmt.Println("value before GC =", wb.Value())
	runtime.GC()
	fmt.Println("value after GC =", wb.Value())
	fmt.Printf("heap size delta = %d KB\n", heapDelta(heapSize))
}
