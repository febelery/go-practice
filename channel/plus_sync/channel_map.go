package main

import (
	"fmt"
	"sync"
	"time"
)

type Request struct {
	op    string       // 存、取、查
	name  string       // 操作的账号
	value int          // 操作金额
	retCh chan *Result // 存放银行处理结果的通道
}

type Result struct {
	success bool // 成功
	value   int  // 查询时使用：余额
}

type Bank struct {
	saving map[string]int // 每账户的存款金额
}

// Loop 银行处理客户请求
func (b *Bank) Loop(reqCh chan *Request) {
	for req := range reqCh {
		switch req.op {
		case "deposite":
			b.Deposite(req)
		case "withdraw":
			b.Withdraw(req)
		case "query":
			b.Query(req)
		default:
			ret := &Result{
				false,
				0,
			}
			req.retCh <- ret
		}
	}
}

// 存款
func (b *Bank) Deposite(req *Request) {
	name := req.name
	amount := req.value

	if _, ok := b.saving[name]; !ok {
		b.saving[name] = 0
	}
	b.saving[name] += amount

	ret := &Result{
		true,
		0,
	}
	req.retCh <- ret
}

// 取款，不足时取款失败
func (b *Bank) Withdraw(req *Request) {
	name := req.name
	amount := req.value

	var status bool

	if balance, ok := b.saving[name]; !ok || balance < amount {
		status = false
		amount = 0
	} else {
		status = true
		b.saving[name] -= amount
	}

	ret := &Result{
		status,
		amount,
	}
	req.retCh <- ret
}

// 查询余额
func (b *Bank) Query(req *Request) {
	name := req.name

	var (
		ok      bool
		balance int
	)

	if balance, ok = b.saving[name]; !ok {
		balance = 0
	}

	ret := &Result{
		true,
		balance,
	}
	req.retCh <- ret
}

func NewBank() *Bank {
	b := &Bank{
		saving: make(map[string]int),
	}

	return b
}

func main() {
	// 创建请求的通道和银行
	reqCh := make(chan *Request, 100)
	bank := NewBank()

	// 银行处理请求
	go bank.Loop(reqCh)

	var wg sync.WaitGroup
	wg.Add(2)

	// 两个顾客同时存取钱
	go customerA(&wg, reqCh)
	go customerB(&wg, reqCh)

	wg.Wait()
	close(reqCh)

	time.Sleep(time.Second)

}

func customerA(wg *sync.WaitGroup, reqCh chan<- *Request) {
	name := "ross"
	retCh := make(chan *Result)
	defer func() {
		close(retCh)
		wg.Done()
	}()

	depReq := &Request{
		"deposite",
		name,
		100,
		retCh,
	}
	withDrawReq := &Request{
		"withdraw",
		name,
		20,
		retCh,
	}
	queryReq := &Request{
		"query",
		name,
		0,
		retCh,
	}

	// 顺序3个请求：存100，花20，剩80
	reqs := []*Request{depReq, withDrawReq, queryReq}
	// 期望Result中返回的值
	expRets := []int{0, 0, 80}
	for i, req := range reqs {
		reqCh <- req
		waitResp(req, expRets[i])
	}
}

func customerB(wg *sync.WaitGroup, reqCh chan<- *Request) {
	name := "rachel"
	retCh := make(chan *Result)
	defer func() {
		close(retCh)
		wg.Done()
	}()

	depReq := &Request{
		"deposite",
		name,
		100,
		retCh,
	}
	withDrawReq := &Request{
		"withdraw",
		name,
		200,
		retCh,
	}
	queryReq := &Request{
		"query",
		name,
		0,
		retCh,
	}

	// 顺序3个请求：存100，花200失败，剩100
	reqs := []*Request{depReq, withDrawReq, queryReq}
	// 期望Result中返回的值
	expRets := []int{0, 0, 100}
	for i, req := range reqs {
		reqCh <- req
		waitResp(req, expRets[i])
	}
}

func waitResp(req *Request, expVal int) {
	ret := <-req.retCh
	if ret.success {
		if req.op != "query" {
			fmt.Printf("%s %s %d success\n", req.name, req.op, req.value)
		} else {
			if ret.value != expVal {
				fmt.Printf("%s query result error, got %d want %d\n", req.name, ret.value, expVal)
			} else {
				fmt.Printf("%s has %d\n", req.name, ret.value)
			}
		}
		return
	}

	if req.op != "query" {
		fmt.Printf("%s %s %d failed\n", req.name, req.op, req.value)
	} else {
		fmt.Printf("%s query failed\n", req.name)
	}
}
