package cbuffer

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

var DefaultAllocator = &Allocator{}

type _slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

type Allocator struct{}

func (a *Allocator) Malloc(size int) []byte {
	var b []byte
	ptr := &b
	s := _slice{
		array: unsafe.Pointer(C.malloc(C.size_t(size))),
		len:   size,
		cap:   size,
	}
	*((*_slice)(unsafe.Pointer(ptr))) = s
	return b
}

func (a *Allocator) Realloc(buf []byte, size int) []byte {
	if size <= cap(buf) {
		return buf[:size]
	}
	ptr := (*_slice)(unsafe.Pointer(&buf))
	ptr.array = C.realloc(ptr.array, C.size_t(size))
	ptr.len = size
	ptr.cap = size
	return buf
}

func (a *Allocator) Append(buf []byte, more ...byte) []byte {
	return a.AppendString(buf, *(*string)(unsafe.Pointer(&more)))
}

func (a *Allocator) AppendString(buf []byte, more string) []byte {
	lbuf, lmore := len(buf), len(more)
	buf = a.Realloc(buf, lbuf+lmore)
	copy(buf[lbuf:], more)
	return buf
}

func (a *Allocator) Free(buf []byte) {
	s := _slice{}
	ptr := &s
	*ptr = *((*_slice)(unsafe.Pointer(&buf)))
	C.free(ptr.array)
}

func Malloc(size int) []byte {
	return DefaultAllocator.Malloc(size)
}

func Realloc(buf []byte, size int) []byte {
	return DefaultAllocator.Realloc(buf, size)
}

func Append(buf []byte, more ...byte) []byte {
	return DefaultAllocator.Append(buf, more...)
}

func AppendString(buf []byte, more string) []byte {
	return DefaultAllocator.AppendString(buf, more)
}

func Free(buf []byte) {
	DefaultAllocator.Free(buf)
}
