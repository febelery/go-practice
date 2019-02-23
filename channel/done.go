package main

import (
	"fmt"
	"sync"
)

type worker3 struct {
	in   chan int
	done func()
}

func doWork(id int, w worker3) {
	for n := range w.in {
		fmt.Printf("Worker %d received %c\n", id, n)
		w.done()
	}
}

func createWorker3(id int, wg *sync.WaitGroup) worker3 {
	w := worker3{
		in: make(chan int),
		done: func() {
			wg.Done()
		},
	}
	go doWork(id, w)
	return w
}

func main() {
	var wg sync.WaitGroup

	var workers [10]worker3
	for i := 0; i < 10; i++ {
		workers[i] = createWorker3(i, &wg)
	}

	wg.Add(20)
	for i, worker := range workers {
		worker.in <- 'a' + i
	}
	for i, worker := range workers {
		worker.in <- 'A' + i
	}

	wg.Wait()
}
