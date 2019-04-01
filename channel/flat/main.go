package main

import "fmt"

/**
平展(flat)操作，如果输入是一个channel，channel中的数据还是相同类型的channel，
那么flat将返回一个输出channel，输出channel中的数据是输入的各个channel中的数据。

它与扇入不同，扇入的输入channel在调用的时候就是固定的，并且以数组的方式提供，
而flat的输入是一个channel，可以运行时随时的加入channel。
*/

func orDone(done <-chan struct{}, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valStream
}

func flat(done <-chan struct{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)

		for {
			var stream <-chan interface{}

			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()

	return valStream
}

func main() {
	genVals := func() <-chan <-chan interface{} {
		chanSteam := make(chan (<-chan interface{}))

		go func() {
			defer close(chanSteam)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanSteam <- stream
			}
		}()

		return chanSteam
	}

	for v := range flat(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
