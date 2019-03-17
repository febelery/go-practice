package main

import "fmt"

func main() {
	// make([]int,5)的含义是创建数组，并且数组初始化5个元素，5个元素的值为类型零值。
	s := make([]int, 5)
	s = append(s, 1, 2, 3)

	fmt.Println(s)
}
