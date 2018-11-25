package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func readerAt() {
	reader := strings.NewReader("学习123456")
	p := make([]byte, 9)
	n, err := reader.ReadAt(p, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s, %d\n", p, n)
}

func writeAt() {
	file, err := os.Create("io/writeAt.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("Golang中文社区——这里是多余")
	n, err := file.WriteAt([]byte("Go语言中文网"), 24)
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
}

func readerFrom() {
	file, err := os.Open("io/writeAt.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(os.Stdout)
	writer.ReadFrom(file)
	writer.Flush()
}

func seeker() {
	reader := strings.NewReader("Golang学习")
	reader.Seek(-6, io.SeekEnd)
	r, _, _ := reader.ReadRune()
	fmt.Printf("%c\n", r)
}

func utf8Index(str, substr string) int {
	index := strings.Index(str, substr)
	if index < 0 {
		return -1
	}
	return utf8.RuneCountInString(str[:index])
}

func limitedReader() {
	content := "This Is LimitReader Example"
	reader := strings.NewReader(content)
	limitReader := &io.LimitedReader{R: reader, N: 7}
	for limitReader.N > 0 {
		tmp := make([]byte, 2)
		limitReader.Read(tmp)
		fmt.Printf("%s", tmp)
	}
	fmt.Println("")
}

func multiReader() {
	readers := []io.Reader{
		strings.NewReader("from string reader"),
		bytes.NewBufferString("from bytes buffer"),
	}
	reader := io.MultiReader(readers...)
	data := make([]byte, 0, 128)
	buf := make([]byte, 10)

	for n, err := reader.Read(buf); err != io.EOF; n, err = reader.Read(buf) {
		if err != nil {
			panic(err)
		}
		data = append(data, buf[:n]...)
	}
	fmt.Printf("%s\n", data)
}

func multiWriter() {
	file, err := os.Create("io/tmp.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writers := []io.Writer{
		file,
		os.Stdout,
	}
	writer := io.MultiWriter(writers...)
	writer.Write([]byte("学习golang\n"))
}

func teeReader() {
	reader := io.TeeReader(strings.NewReader("How to learn golang.\n"), os.Stdout)
	reader.Read(make([]byte, 20))
}

func main() {

	readerAt()
	writeAt()
	readerFrom()
	seeker()
	fmt.Println(utf8Index("行业交流群", "交流"))
	limitedReader()
	multiReader()
	multiWriter()
	teeReader()

}
