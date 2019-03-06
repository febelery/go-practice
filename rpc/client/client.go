package main

import (
	"fmt"
	"learn/rpc"
	"net"
	"net/rpc/jsonrpc"
)

func main() {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	clients := jsonrpc.NewClient(conn)

	var result float64
	err = clients.Call("DemoService.Div", rpcdemo.Args{10, 3}, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	err = clients.Call("DemoService.Div", rpcdemo.Args{10, 0}, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
