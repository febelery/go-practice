package main

import (
	"fmt"
	"time"
)

type Mutex3 struct {
	ch chan struct{}
}

func NewMutex3() *Mutex3 {
	mu := &Mutex3{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

func (m *Mutex3) Lock() {
	<-m.ch
}

func (m *Mutex3) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *Mutex3) TryLock(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-time.After(timeout):
	}
	return false
}

func (m *Mutex3) IsLocked() bool {
	return len(m.ch) == 0
}

func main() {
	m := NewMutex3()
	ok := m.TryLock(time.Second)
	fmt.Printf("locked v %v\n", ok)

	ok = m.TryLock(time.Second)
	fmt.Printf("locked %v\n", ok)
}
