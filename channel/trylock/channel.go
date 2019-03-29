package main

import "fmt"

type Mutex2 struct {
	ch chan struct{}
}

func NewMutex() *Mutex2 {
	// 利用channel边界情况下的阻塞特性实现的
	mu := &Mutex2{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

func (m *Mutex2) Lock() {
	<-m.ch
}

func (m *Mutex2) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *Mutex2) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *Mutex2) IsLocked() bool {
	return len(m.ch) == 0
}

func main() {
	m := NewMutex()
	ok := m.TryLock()
	fmt.Printf("locked v %v\n", ok)

	ok = m.TryLock()
	fmt.Printf("locked %v\n", ok)
}
