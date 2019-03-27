package main

import (
	"flag"
	"fmt"
	"learn/tcp_server/common"
	"log"
	"net"
	"time"
)

var (
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 100, "number of tcp connections")
)

func main() {
	flag.Parse()

	common.SetLimit()

	addr := *ip + ":8972"
	log.Printf("连接到 %s", addr)
	var conns []net.Conn

	for i := 0; i < *connections; i++ {
		c, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			fmt.Println("failed to connect", i, err)
			i--
			time.Sleep(time.Second)
			continue
		}
		conns = append(conns, c)
		time.Sleep(time.Millisecond)
	}

	defer func() {
		for _, c := range conns {
			c.Close()
		}
	}()

	log.Printf("完成初始化 %d 连接", len(conns))

	tts := time.Second
	if *connections > 100 {
		tts = time.Millisecond * 5
	}

	for {
		for i := 0; i < len(conns); i++ {
			time.Sleep(tts)
			conn := conns[i]
			conn.Write([]byte(fmt.Sprintf("Client: %d. Message: Hello World\r\n", i)))
		}
	}
}
