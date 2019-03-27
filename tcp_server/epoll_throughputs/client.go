package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

var (
	ip          = flag.String("ip", "localhost", "server IP")
	connections = flag.Int("conn", 1, "number of tcp connections")
	startMetric = flag.String("sm", time.Now().Format("2019-03-25T15:04:05 +0800"), "start time point of all clients")
)

func main() {
	flag.Parse()

	setLimit()

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

func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	rLimit.Cur = rLimit.Max

	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Printf("set cur limit: %d", rLimit.Cur)
}
