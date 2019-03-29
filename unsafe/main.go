package main

import (
	"fmt"
	"unsafe"
)

type user struct {
	name string
	age  int
}

func pointer() {
	i := 10
	ip := &i

	// unsafe.Pointer是一种特殊意义的指针，它可以包含任意类型的地址
	var fp *float64 = (*float64)(unsafe.Pointer(ip))

	*fp = *fp * 3

	fmt.Println(i)
}

func uintptrT() {
	u := new(user)
	fmt.Println(*u)

	// 我们一般使用*T作为一个指针类型，表示一个指向类型T变量的指针
	pName := (*string)(unsafe.Pointer(u))
	*pName = "ross"

	// 因为age不是第一个字段，所以我们需要内存偏移，内存偏移牵涉到的计算只能通过uintptr，
	// 所我们要先把user的指针地址转为uintptr，然后我们再通过unsafe.Offsetof(u.age)获取需要偏移的值，进行地址运算(+)偏移即可
	pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + unsafe.Offsetof(u.age)))
	*pAge = 20

	fmt.Println(*u)
}

func sizeOf() {
	fmt.Println(unsafe.Sizeof(true))
	fmt.Println(unsafe.Sizeof(int8(0)))
	fmt.Println(unsafe.Sizeof(int16(10)))
	fmt.Println(unsafe.Sizeof(int32(10000000)))
	fmt.Println(unsafe.Sizeof(int64(10000000000000)))
	fmt.Println(unsafe.Sizeof(int(10000000000000000)))
}

func main() {
	fmt.Println("指针转换为unsafe.Pointer")
	pointer()

	fmt.Println("uintptr转换为unsafe.Pointer")
	uintptrT()

	fmt.Println("sunsafe sizeof 函数")
	sizeOf()

}
