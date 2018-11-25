package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readSlice() {
	reader := bufio.NewReader(strings.NewReader("good good laugh. \nday day run"))
	line, _ := reader.ReadSlice('\n')
	fmt.Printf("the line:%s\n", line)
	n, _ := reader.ReadSlice('\n')
	fmt.Printf("the line:%s\n", line)
	fmt.Println(string(n))
}

func readBytes() {
	reader := bufio.NewReader(strings.NewReader("good good laugh. \nday day run"))
	line, _ := reader.ReadBytes('\n')

	//line = bytes.TrimRight(line, "\r\n")

	fmt.Printf("the line:%s\n", line)
	n, _ := reader.ReadBytes('\n')
	fmt.Printf("the line:%s\n", line)
	fmt.Println(string(n))
}

func scanner() {
	file, err := os.Create("io/scanner.txt")
	if err != nil {
		panic(file)
	}
	defer file.Close()
	file.WriteString("hello world.\ntomarrow will be good than today.\nsee you again.\nnice to meet you.")
	file.Seek(0, os.SEEK_SET)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func scannerStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func scannerCountWord() {
	const input = "This is The Golang Standard Library.\nWelcome you!"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	fmt.Println(count)
}

func main() {

	readSlice()
	readBytes()
	scanner()
	scannerCountWord()

}
