package cbuffer

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

var DefaultAllocator = &Allocator{}

var cgoMux sync.Mutex

type _slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}

type Allocator struct{}

func (a *Allocator) Malloc(size int) []byte {
	var buf []byte
	ptr := &buf
	cgoMux.Lock()
	pc := C.calloc(1, C.size_t(size))
	cgoMux.Unlock()
	s := _slice{
		array: unsafe.Pointer(pc),
		len:   size,
		cap:   size,
	}
	*((*_slice)(unsafe.Pointer(ptr))) = s
	BufferMap.Add(&buf[0])
	return buf
}

func (a *Allocator) Realloc(buf []byte, size int) []byte {
	if size <= cap(buf) {
		return buf[:size]
	}
	p1 := &buf[0]
	ptr := (*_slice)(unsafe.Pointer(&buf))
	cgoMux.Lock()
	pc := C.realloc(ptr.array, C.size_t(size))
	cgoMux.Unlock()
	ptr.array = unsafe.Pointer(pc)
	ptr.len = size
	ptr.cap = size
	p2 := &buf[0]
	if p1 != p2 {
		BufferMap.Delete(p1)
		BufferMap.Add(p2)
	}
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
	if BufferMap.Delete(&buf[0]) {
		a.free(buf)
	}
}

func (a *Allocator) free(buf []byte) {
	s := _slice{}
	ptr := &s
	*ptr = *((*_slice)(unsafe.Pointer(&buf)))
	cgoMux.Lock()
	C.free(ptr.array)
	cgoMux.Unlock()
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

var (
	BufferMap = newBufferMap(64)
)

func newBufferMap(bucketNum int) *bufferMap {
	bm := &bufferMap{
		buckets: make([]*bucket, 0, bucketNum),
	}
	for i := 0; i < bucketNum; i++ {
		bm.buckets = append(bm.buckets, &bucket{values: map[uintptr]struct{}{}})
	}
	return bm
}

type bufferMap struct {
	buckets []*bucket
}

func (bm *bufferMap) Buckets() []*bucket {
	return bm.buckets
}

func (bm *bufferMap) Add(pb *byte) {
	ptr := uintptr(unsafe.Pointer(pb))
	idx := hash(ptr) % uint64(len(bm.buckets))
	bucket := bm.buckets[idx]
	bucket.Add(ptr)
}

func (bm *bufferMap) Delete(pb *byte) bool {
	ptr := uintptr(unsafe.Pointer(pb))
	idx := hash(ptr) % uint64(len(bm.buckets))
	bucket := bm.buckets[idx]
	return bucket.Delete(ptr)
}

func (bm *bufferMap) Length() int {
	length := 0
	for _, bucket := range bm.buckets {
		length += bucket.Length()
	}
	return length
}

func (bm *bufferMap) Distribution() []int {
	ret := make([]int, len(bm.buckets))
	for i, bucket := range bm.buckets {
		ret[i] = bucket.Length()
	}
	return ret
}

type bucket struct {
	sync.Mutex
	values map[uintptr]struct{}
}

func (b *bucket) Add(ptr uintptr) {
	b.Lock()
	_, ok := b.values[ptr]
	if ok {
		panic(fmt.Errorf("buffer conflict: %v", ptr))
	}
	b.values[ptr] = struct{}{}
	b.Unlock()
}

func (b *bucket) Delete(ptr uintptr) bool {
	b.Lock()
	_, ok := b.values[ptr]
	if ok {
		delete(b.values, ptr)
	}
	b.Unlock()
	return ok
}

func (b *bucket) Length() int {
	b.Lock()
	length := len(b.values)
	b.Unlock()
	return length
}

func hash(n uintptr) uint64 {
	var h = uint64(0)
	for ; n > 0; n >>= 8 {
		h = 31*h + uint64(n&0xFF)
	}
	return h
}
