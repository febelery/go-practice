package common

import (
	"github.com/rcrowley/go-metrics"
	"io"
	"log"
	"net"
	"sync"
)

type Pool struct {
	workers   int
	maxTasks  int
	taskQueue chan net.Conn
	mu        sync.Mutex
	closed    bool
	done      chan struct{}
}

func NewPool(w int, t int) *Pool {
	return &Pool{
		workers:   w,
		maxTasks:  t,
		taskQueue: make(chan net.Conn, t),
		done:      make(chan struct{}),
	}
}

func (p *Pool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *Pool) AddTask(conn net.Conn) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.taskQueue <- conn
}

func (p *Pool) Start(epoller Epoll, opsRate metrics.Meter) {
	for i := 0; i < p.workers; i++ {
		go p.startWorker(epoller, opsRate)
	}
}

func (p *Pool) startWorker(epoller Epoll, opsRate metrics.Meter) {
	for {
		select {
		case <-p.done:
			return
		case conn := <-p.taskQueue:
			if conn != nil {
				handleClose(conn, epoller, opsRate)
			}
		}
	}
}

func handleClose(conn net.Conn, epoller Epoll, opsRate metrics.Meter) {
	_, err := io.CopyN(conn, conn, 8)
	if err != nil {
		if err := epoller.Remove(conn); err != nil {
			log.Printf("failed to remove %v", err)
		}
		conn.Close()
	}
	opsRate.Mark(1)
}
