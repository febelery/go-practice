package main

import (
	"flag"
	"github.com/libp2p/go-reuseport"
	"github.com/rcrowley/go-metrics"
	"io"
	"learn/tcp_server/common"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	c = flag.Int("c", 10, "concurrency")
)

var (
	opsRate = metrics.NewRegisteredMeter("ops", nil)
)

func main() {
	flag.Parse()

	common.SetLimit()

	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	for i := 0; i < *c; i++ {
		go startEpoll()
	}

	select {}
}

func startEpoll() {
	ln, err := reuseport.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}

	epoller, err := common.MkEpoll()
	if err != nil {
		panic(err)
	}

	go start(epoller)

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			log.Printf("accept err: %v", e)
			return
		}

		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}

func start(epoller *common.Epoll) {
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
			io.CopyN(conn, conn, 8)

			if err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
					conn.Close()
				}
			}

			opsRate.Mark(1)
		}
	}
}
