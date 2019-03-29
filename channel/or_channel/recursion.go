package main

import (
	"fmt"
	"time"
)

// 递归

func orRecursion(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			m := len(channels) / 2
			select {
			case <-orRecursion(channels[:m]...):
			case <-orRecursion(channels[m:]...):
			}
		}
	}()

	return orDone
}

func sig2(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {

	start := time.Now()

	<-orRecursion(
		sig2(10*time.Second),
		sig2(20*time.Second),
		sig2(30*time.Second),
		sig2(40*time.Second),
		sig2(50*time.Second),
		sig2(01*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
