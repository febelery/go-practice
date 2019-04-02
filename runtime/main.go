package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
)

func Foo() {
	fmt.Printf("我是 %s, %s 在调用我!\n", printMyName(), printCallerName())
	Bar()
	trace()
	fmt.Println("***************************************************")
	trace2()
	fmt.Println("***************************************************")
	DumpStacks()
}

func Bar() {
	fmt.Printf("我是 %s, %s 又在调用我!\n", printMyName(), printCallerName())
}

func printMyName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func printCallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func trace() {
	pc := make([]uintptr, 10)
	n := runtime.Callers(0, pc)

	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		fmt.Printf("%s:%d %s\n", file, line, f.Name())
	}
}

func trace2() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	n := runtime.Callers(0, pc)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		fmt.Printf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
}

func DumpStacks() {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, true)]
	fmt.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===", buf)
}

func PrintStack() {
	debug.PrintStack()
}

func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func main() {
	Foo()
}
