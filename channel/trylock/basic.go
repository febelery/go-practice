package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

const mutexLocked = 1 << iota

type Mutex struct {
	mu sync.Mutex
}

func (m *Mutex) Lock() {
	m.mu.Lock()
}

func (m *Mutex) Unlock() {
	m.mu.Unlock()
}

// 可以使用这个方法避免在获取锁的时候当前goroutine被阻塞住。
func (m *Mutex) TryLock() bool {
	// mutex实现锁主要利用CAS对它的一个int32字段做操作
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.mu)), 0, mutexLocked)
}

func (m *Mutex) IsLocked() bool {
	return atomic.LoadInt32((*int32)(unsafe.Pointer(&m.mu))) == mutexLocked
}

func main() {
	var m Mutex
	if m.TryLock() {
		fmt.Println("locked: ", m.IsLocked())
	} else {
		fmt.Println("failed to lock")
	}
}
