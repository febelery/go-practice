// 探索Golang定时器的陷阱
// 当前一次的超时数据没有被读取，而设置了新的定时器，然后去通道读数据，
// 结果读到的是上次超时的超时事件，看似成功，实则失败，完全掉入陷阱

package main

import (
	"fmt"
	"time"
)

func test1() {
	tm := time.NewTimer(time.Second)
	quit := make(chan bool)

	go func() {
		time.Sleep(3 * time.Second)
		quit <- true
	}()

	if tm.Reset(time.Second) {
		fmt.Println("reset成功，原timer未超时，Reset返回true")
	} else {
		fmt.Println("reset失败，原timer已超时或被停止，Reset返回true")
	}

	tm.Stop()
	if tm.Reset(time.Second) {
		fmt.Println("reset成功，停止Timer后，Reset返回true")
	} else {
		fmt.Println("reset失败，停止Timer后，Reset返回false")
	}

	for {
		select {
		case <-quit:
			return
		case <-tm.C:
			if tm.Reset(time.Second) {
				fmt.Println("超时，Reset返回true")
			} else {
				fmt.Println("超时，Reset返回false")
			}
		}
	}
}

func test2() {
	tStart := time.Now()
	tm := time.NewTimer(time.Second)

	time.Sleep(time.Second * 2)
	fmt.Printf("Reset 前通道中的事件数量为：%d \n", len(tm.C))

	if tm.Reset(time.Second) {
		fmt.Println("不读通道数据，Reset返回true")
	} else {
		fmt.Println("不读通道数据，Reset返回false")
	}

	fmt.Printf("Reset 后通道中的事件数量为：%d \n", len(tm.C))

	select {
	case t := <-tm.C:
		fmt.Printf("tm 开始的时间: %v\n", tStart.Unix())
		fmt.Printf("通道中事件的时间：%v\n", t.Unix())
		if t.Sub(tStart) <= time.Second+time.Millisecond {
			fmt.Println("通道中的时间是重新设置sm 前的时间，即第一次超时的时间，所以第二次Reset失败了")
		}
	}

	fmt.Printf("读通道后，其中事件的数量:%d\n", len(tm.C))

	tm.Reset(time.Second)
	fmt.Printf("再次Reset后，通道中事件的数量:%d\n", len(tm.C))

	time.Sleep(2 * time.Second)
	fmt.Printf("超时后通道中事件的数量:%d\n", len(tm.C))
}

func test3() {
	tStart := time.Now()
	tm := time.NewTimer(time.Second)
	time.Sleep(2 * time.Second)

	// 停掉定时器再清空
	// 定时器的运行和len(Timer.C)的判断是在不同的协程中，当判断的时候通道大小可能为0，但当执行Reset()的前的这段时间，
	// 		旧的定时器超时，通道中存在超时事件，再执行Reset()也达不到预期的效果
	// 先执行Stop()，可以确保旧定时器已经停止，不会再向通道中写入超时事件
	// Stop停止Timer的执行。如果停止了t会返回真；如果t已经被停止或者过期了会返回假
	if !tm.Stop() && len(tm.C) > 0 {
		<-tm.C
	}
	tm.Reset(time.Second)

	// 超时
	t := <-tm.C
	fmt.Printf("tm开始的时间: %v\n", tStart.Unix())
	fmt.Printf("通道中事件的时间：%v\n", t.Unix())

	if t.Sub(tStart) <= time.Second+time.Millisecond {
		fmt.Println("通道中的时间是重新设置sm前的时间，即第一次超时的时间，所以第二次Reset失败了")
	} else {
		fmt.Println("通道中的时间是重新设置sm后的时间，Reset成功了")
	}
}

func main() {
	fmt.Println("<<<<<<<<<<<<<<<<<<<<第1个测试：Reset返回值和什么有关？>>>>>>>>>>>>>>>>>>>>>")
	test1()
	fmt.Println("<<<<<<<<<<<<<<<<<<<<第2个测试:超时后，不读通道中的事件，可以Reset成功吗？>>>>>")
	test2()
	fmt.Println("<<<<<<<<<<<<<<<<<<<<第3个测试：Reset前清空通道，尽可能通畅>>>>>>>>>>>>>>>>>>")
	test3()
}
