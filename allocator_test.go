package cbuffer

import (
	"fmt"
	"testing"
)

func BenchmarkMallocFree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := Malloc(1024)
		Free(b)
	}
}

func BenchmarkRealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := Malloc(1024)
		b = Realloc(b, 2048)
		Free(b)
	}
}

func BenchmarkAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b1 := Malloc(1024)
		b2 := Malloc(1024)
		b1 = Append(b1, b2...)
		Free(b1)
		Free(b2)
	}
}

func TestAllocator(t *testing.T) {
	str := "hello world"
	buf := DefaultAllocator.Malloc(len(str))
	copy(buf, str)
	fmt.Println("s 111:", string(buf))
	buf = DefaultAllocator.Realloc(buf, len(str)+1)
	buf[len(buf)-1] = '!'
	fmt.Println("s 222:", string(buf))
	buf = DefaultAllocator.Append(buf, ' ', 'h')
	fmt.Println("s 333:", string(buf))
	buf = DefaultAllocator.AppendString(buf, "ello world!")
	fmt.Println("s 444:", string(buf))
	DefaultAllocator.Free(buf)
	fmt.Println("s 555:", len(buf), cap(buf), string(buf))
}
