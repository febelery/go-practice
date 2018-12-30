package main

import (
	"bufio"
	"fmt"
	"io"
	"learn/functional/fib"
	"strings"
)

type intGen func() int

func (g intGen) Read(p []byte) (n int, err error) {
	next := g()
	if next > 10000 {
		return 0, io.EOF
	}

	s := fmt.Sprintf("%d\n", next)

	return strings.NewReader(s).Read(p)
}

func main() {
	var f intGen = fib.Fibonacci()

	printFileContents(f)
}

func printFileContents(reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
