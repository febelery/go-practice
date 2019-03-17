package main

import "fmt"

type student struct {
	Name string
	Age  int
}

func main() {
	m := make(map[string]*student)
	stus := []student{
		{Name: "zhou", Age: 24},
		{Name: "li", Age: 23},
		{Name: "wang", Age: 22},
	}

	fmt.Printf("Before: %+v\n", m)

	// 在range循环中，变量是不会随着遍历过程发生变化的。
	// 因此在代码中stu是不会变化的，变化的是放在&stu地址上的数据，
	// 因此最后m中value都将是最后一个放在&stu中的值
	for _, stu := range stus {
		m[stu.Name] = &stu
	}

	fmt.Printf("After: %+v\n", m)
}
