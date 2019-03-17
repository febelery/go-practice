package main

import "fmt"

func main() {
	deferCall()
}

func deferCall() {
	// defer是先进后出，逆序执行。
	defer func() {
		fmt.Println("one")
	}()

	defer func() {
		fmt.Println("two")
	}()

	defer func() {
		fmt.Println("three")
	}()

	panic("error")
}
