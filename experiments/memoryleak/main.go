package main

import (
	"fmt"
	"runtime"
)

var b []int

func f(a []int) {
	b = a[1000:]
}

func main() {
	m := &runtime.MemStats{}
	for {
		a := []int{}
		for i := 0; i < 1001; i++ {
			a = append(a, i)
		}
		runtime.ReadMemStats(m)
		f(a)
		fmt.Println(b, m.Alloc)
	}
}
