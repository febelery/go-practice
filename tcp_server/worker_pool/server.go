package main

import (
	"flag"
	"github.com/rcrowley/go-metrics"
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

var epoller *common.Epoll
var workerPool *common.Pool

func main() {
	flag.Parse()

	common.SetLimit()

	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	epoller, err = common.MkEpoll()

	workerPool = common.NewPool(*c, 1000000)
	workerPool.Start(*epoller, opsRate)

	if err != nil {
		panic(err)
	}

	go start()

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

	workerPool.Close()
}

func start() {
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

			workerPool.AddTask(conn)
		}
	}
}
