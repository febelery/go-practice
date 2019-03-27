package main

import (
	"flag"
	"io"
	"learn/tcp_server/common"
	"log"
	"net"
	"os"
	"os/exec"
)

var (
	c       = flag.Int("c", 1, "concurrency")
	prefork = flag.Bool("prefork", false, "use prefork")
	child   = flag.Bool("child", false, "is child proc")
)

func main() {
	flag.Parse()

	common.SetLimit()

	var ln net.Listener
	var err error

	if *prefork {
		ln = doPrefork(*c)
	} else {
		ln, err = net.Listen("tcp", ":8972")
		if err != nil {
			panic(err)
		}
	}

	startEpoll(ln)

	select {}

}

func startEpoll(ln net.Listener) {
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
				}
				conn.Close()
			}
		}
	}
}

func doPrefork(c int) net.Listener {
	var listener net.Listener
	if !*child {
		addr, err := net.ResolveTCPAddr("tcp", ":8972")
		if err != nil {
			log.Fatal(err)
		}

		tcplistner, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}

		fl, err := tcplistner.File()
		if err != nil {
			log.Fatal(err)
		}

		children := make([]*exec.Cmd, c)
		for i := range children {
			children[i] = exec.Command(os.Args[0], "-prefork", "-child")
			children[i].Stdout = os.Stdout
			children[i].Stderr = os.Stderr
			children[i].ExtraFiles = []*os.File{fl}
			err = children[i].Start()
			if err != nil {
				log.Fatalf("failed to start child: %v", err)
			}
		}

		for _, ch := range children {
			if err := ch.Wait(); err != nil {
				log.Fatalf("failed to wait child's starting: %v", err)
			}
		}
	} else {
		var err error
		listener, err = net.FileListener(os.NewFile(3, ""))
		if err != nil {
			log.Fatal(err)
		}
	}

	return listener
}
