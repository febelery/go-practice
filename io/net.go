package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
)

func read() {
	conn, err := net.Dial("tcp", "rpcx.site:80")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	var sb strings.Builder
	buf := make([]byte, 256)
	for {
		// io.Reader接口定义了Read(p []byte) (n int, err error)方法，我们可以使用它从Reader中读取一批数据。
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		sb.Write(buf[:n])
	}

	fmt.Println("response: ", sb.String())
	fmt.Println("total response size: ", sb.Len())
}

func readAll() {
	conn, err := net.Dial("tcp", "rpcx.site:80")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		if err != io.EOF {
			fmt.Println("read error:", err)
		}
		panic(err)
	}

	fmt.Println("response:", string(data))
	fmt.Println("total response size:", len(data))
}

func readAtLeast() {
	conn, err := net.Dial("tcp", "rpcx.site:80")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	// 发送请求, http 1.0 协议
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	// 读取response
	var sb strings.Builder
	buf := make([]byte, 256)
	for {
		n, err := io.ReadAtLeast(conn, buf, 256)
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				fmt.Println("read error:", err)
			}
			break
		}
		sb.Write(buf[:n])
	}
	// 显示结果
	fmt.Println("response:", sb.String())
	fmt.Println("total response size:", sb.Len())
}

func readLimit() {
	conn, err := net.Dial("tcp", "rpcx.site:80")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	var sb strings.Builder
	buf := make([]byte, 256)
	rr := io.LimitReader(conn, 102400)
	for {
		n, err := io.ReadAtLeast(rr, buf, 256)
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				fmt.Println("read error:", err)
			}
			break
		}
		sb.Write(buf[:n])
	}
	// 显示结果
	fmt.Println("response:", sb.String())
	fmt.Println("total response size:", sb.Len())
}

func copy() {
	conn, err := net.Dial("tcp", "rpcx.site:80")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	var sb strings.Builder
	_, err = io.Copy(&sb, conn)
	if err != nil {
		fmt.Println("read error:", err)
	}

	fmt.Println("response: ", sb.String())
	fmt.Println("total response size: ", sb.Len())
}

func main() {
	read()
	copy()
	readAll()
}
