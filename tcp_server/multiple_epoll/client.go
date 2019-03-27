package main

import (
	"encoding/binary"
	"flag"
	"github.com/rcrowley/go-metrics"
	"learn/tcp_server/common"
	"log"
	"net"
	"os"
	"time"
)

var (
	ip          = flag.String("ip", "127.0.0.1", "server IP")
	connections = flag.Int("conn", 10, "number of total tcp connections")
	cc          = flag.Int("c", 100, "currency count")
	startMetric = flag.String("sm", time.Now().Format("2019-01-01T00:00:00 +0800"), "start time point of all clients")
)

var (
	opsRateC = metrics.NewRegisteredTimer("ops", nil)
)

func main() {
	flag.Parse()

	common.SetLimit()

	go func() {
		startPoint, _ := time.Parse("2019-01-01T00:00:00 +0800", *startMetric)
		time.Sleep(startPoint.Sub(time.Now()))

		metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	}()

	addr := *ip + ":8972"
	log.Printf("连接到 %s", addr)

	for i := 0; i < *cc; i++ {
		go mkClient(addr, *connections/(*cc))
	}

	select {}
}

func mkClient(addr string, connections int) {
	epoller, err := common.MkEpoll()
	if err != nil {
		panic(err)
	}

	var conns []net.Conn
	for i := 0; i < connections; i++ {
		c, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			log.Println("failed to connect", i, err)
			i--
			continue
		}
		if err := epoller.Add(c); err != nil {
			log.Printf("failed to add connection %v", err)
			c.Close()
		}
		conns = append(conns, c)
	}

	log.Printf("完成初始化 %d 连接", len(conns))

	go startC(epoller)

	tts := time.Second
	if *cc > 100 {
		tts = time.Millisecond * 5
	}

	for i := 0; i < len(conns); i++ {
		time.Sleep(tts)
		conn := conns[i]

		err = binary.Write(conn, binary.BigEndian, time.Now().UnixNano())
		if err != nil {
			log.Printf("failed to write timestamp %v", err)

			if err := epoller.Remove(conn); err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
			}
		}
	}

	select {}
}

func startC(epoller *common.Epoll) {
	var nano int64
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}

			if err := binary.Read(conn, binary.BigEndian, &nano); err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
			} else {
				opsRateC.Update(time.Duration(time.Now().UnixNano() - nano))
			}

			err = binary.Write(conn, binary.BigEndian, time.Now().UnixNano())
			if err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
			}
		}
	}

}
