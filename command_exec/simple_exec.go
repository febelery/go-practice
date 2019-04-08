package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func runAndExec() {
	cmd := exec.Command("echo", "hello")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	fmt.Printf("combined out: \n%s\n", string(out))
}

func captureStdoutStdErr() {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("echo", "hello")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

func main() {
	runAndExec()
	captureStdoutStdErr()
}
