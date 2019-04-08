package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func simplePipe() {
	c1 := exec.Command("ls")
	c2 := exec.Command("wc", "-l")

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()

	io.Copy(os.Stdout, &b2)
}

func stdoutPipe() {
	c1 := exec.Command("ls /")
	c2 := exec.Command("wc", "-l")

	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = os.Stdout

	c2.Start()
	c1.Run()
	c2.Wait()
}

func trickPipe() {
	cmd := "cat /proc/cpuinfo | egrep '^model name' | uniq | awk '{print substr($0, index($0,$4))}'"

	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("Failed to execute command: %s", cmd)
	}

	fmt.Println(string(out))
}

func main() {
	simplePipe()
	stdoutPipe()
	trickPipe()
}
