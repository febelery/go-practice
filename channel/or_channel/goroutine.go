package main

import (
	"fmt"
	"sync"
	"time"
)

// 提供相同服务的n个节点发送请求，只要任意一个服务节点返回结果，
// 我们就可以执行下面的业务逻辑，其它n-1的节点的请求可以被取消或者忽略
func orGoroutine(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var once sync.Once
		for _, c := range chans {
			go func(c <-chan interface{}) {
				select {
				case <-c:
					once.Do(func() {
						close(out)
					})
				case <-out:
				}
			}(c)
		}
	}()

	return out
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-orGoroutine(
		sig(5*time.Second),
		sig(12*time.Second),
		sig(3*time.Second),
		sig(10*time.Second),
		sig(8*time.Second),
		sig(2*time.Second),
	)

	fmt.Printf("done after %v", time.Since(start))
}
