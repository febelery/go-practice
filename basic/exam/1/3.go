package main

import (
	"fmt"
	"runtime"
	"sync"
)

/*
先说说WaitGroup的用途：它能够一直等到所有的goroutine执行完成，并且阻塞主线程的执行，直到所有的goroutine执行完成。
这里要注意一下，他们的执行结果是没有顺序的，调度器不能保证多个 goroutine 执行次序，且进程退出时不会等待它们结束。
WaitGroup总共有三个方法：Add(delta int),Done(),Wait()。简单的说一下这三个方法的作用。
Add:添加或者减少等待goroutine的数量
Done:相当于Add(-1)
Wait:执行阻塞，直到所有的WaitGroup数量变成0
*/
func main() {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("One i: ", i)
			wg.Done()
		}()
	}

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("Two i: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
