package main

import (
	"fmt"
	"github.com/psilva261/timsort"
	"math/rand"
	"sort"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	testData1 := make([]int, 0, 1000000)
	testData2 := make([]int, 0, 1000000)
	testData3 := make([]int, 0, 1000000)
	testData4 := make([]int, 0, 1000000)

	times := 1000000
	for i := 0; i < times; i++ {
		val := rand.Intn(200000000)
		testData1 = append(testData1, val)
		testData2 = append(testData2, val)
		testData3 = append(testData3, val)
		testData4 = append(testData4, val)
	}

	start := time.Now()
	quickSort(testData1, 0, len(testData1)-1)
	fmt.Println("single goroutine: ", time.Now().Sub(start))
	if !sort.IntsAreSorted(testData1) {
		fmt.Println("wrong quick_sort implementation")
	}

	done := make(chan struct{})
	start = time.Now()
	go quickSortGo(testData2, 0, len(testData2)-1, done, 5)
	<-done
	fmt.Println("multiple goroutine: ", time.Now().Sub(start))

	start = time.Now()
	sort.Ints(testData3)
	fmt.Println("std lib: ", time.Now().Sub(start))

	start = time.Now()
	timsort.Ints(testData4, func(a, b int) bool { return a <= b })
	fmt.Println("timsort: ", time.Now().Sub(start))
}

func quickSortGo(a []int, lo, hi int, done chan struct{}, depth int) {
	if lo >= hi {
		done <- struct{}{}
		return
	}

	depth--
	p := partition(a, lo, hi)

	if depth > 0 {
		childDone := make(chan struct{}, 2)

		go quickSortGo(a, lo, p-1, childDone, depth)
		go quickSortGo(a, p+1, hi, childDone, depth)

		<-childDone
		<-childDone
	} else {
		quickSort(a, lo, p-1)
		quickSort(a, p+1, hi)
	}

	done <- struct{}{}
}

func quickSort(a []int, lo, hi int) {
	if lo >= hi {
		return
	}

	p := partition(a, lo, hi)

	quickSort(a, lo, p-1)
	quickSort(a, p+1, hi)
}

func partition(a []int, lo, hi int) int {
	pivot := a[hi]

	i := lo - 1
	for j := lo; j < hi; j++ {
		if a[j] < pivot {
			i++
			a[j], a[i] = a[i], a[j]
		}
	}

	a[i+1], a[hi] = a[hi], a[i+1]
	return i + 1
}
