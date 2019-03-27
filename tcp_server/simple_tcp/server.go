package main

import (
	"bufio"
	"learn/tcp_server/common"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	common.SetLimit()

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
