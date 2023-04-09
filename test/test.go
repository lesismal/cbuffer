package main

/*
#include <stdlib.h>

int Sum(int a, int b) {
	return a + b;
}
*/
import "C"

import (
	"fmt"
	"time"

	"github.com/lesismal/cbuffer"
)

var loop = 10000000

func main() {
	t := time.Now()
	for i := 0; i < loop; i++ {
		sum := C.Sum(C.int(i), C.int(i+1))
		sum += 1
	}
	used := time.Since(t).Nanoseconds()
	fmt.Printf("[sum test] time used: %dns, %dns/op\n", used, used/int64(loop))

	t = time.Now()
	for i := 0; i < loop; i++ {
		p := cbuffer.Malloc(1024)
		p[i%1024] = byte(i % 256)
		cbuffer.Free(p)
	}
	used = time.Since(t).Nanoseconds()
	fmt.Printf("[allocator test] time used: %dns, %dns/op\n", used, used/int64(loop))
}
