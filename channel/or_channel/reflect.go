package main

import (
	"fmt"
	"reflect"
	"time"
)

// 利用反射“随机”地从一组可选的channel中接收数据，并关闭输出channel
func orReflect(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		var cases []reflect.SelectCase

		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 取出最先到达的，其他忽略
		reflect.Select(cases)
	}()

	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()

	return c
}

func main() {
	start := time.Now()

	<-orReflect(
		sig(10*time.Second),
		sig(30*time.Second),
		sig(40*time.Second),
		sig(50*time.Second),
		sig(01*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
