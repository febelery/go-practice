package common

import (
	"bufio"
	"log"
	"net"
	"strings"
	"syscall"
)

func SetLimit() {
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

func ReadData(conn net.Conn) {
	reader := bufio.NewReader(conn)

	line, _, _ := reader.ReadLine()
	if len(line) > 0 {
		log.Printf("server recieve data: %s \n", strings.TrimSpace(string(line)))
	}
}
