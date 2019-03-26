package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"strings"
	"syscall"
)

func main() {
	setLimit()

	ln, e := net.Listen("tcp", ":8972")
	if e != nil {
		panic(e)
	}

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	var connections []net.Conn
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}
		}

		go handleConn(conn)
		connections = append(connections, conn)
		if len(connections)%100 == 0 {
			log.Printf("total number of connections: %d", len(connections))
		}
	}
}

func handleConn(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		line, _ := reader.ReadBytes('\n')
		if len(line) > 0 {
			log.Printf("server recieve data: %s \n", strings.TrimSpace(string(line)))
		}
	}
	//io.Copy(ioutil.Discard, conn)
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
