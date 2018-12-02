package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Join(str []string, sep string) string {
	if len(str) == 0 {
		return ""
	}
	if len(str) == 1 {
		return str[0]
	}

	buffer := bytes.NewBufferString(str[0])
	for _, s := range str[1:] {
		buffer.WriteString(sep)
		buffer.WriteString(s)
	}

	return buffer.String()
}

func main() {
	s := "Hello世界真大"
	fmt.Println(s)

	for _, b := range s {
		fmt.Printf("%X ", b)
	}
	fmt.Println()

	for i, ch := range s { // ch is a rune
		fmt.Printf("(%d %X) ", i, ch)
	}
	fmt.Println()

	fmt.Println("Rune count: ", utf8.RuneCountInString(s))

	bytes := []byte(s)
	for len(bytes) > 0 {
		ch, size := utf8.DecodeRune(bytes)
		bytes = bytes[size:]
		fmt.Printf("%c ", ch)
	}
	fmt.Println()

	for i, ch := range []rune(s) {
		fmt.Printf("(%d %c) ", i, ch)
	}
	fmt.Println()

	fmt.Println(strings.ContainsAny("failure", "u & i"))
	fmt.Println(strings.Count("five", ""))
	fmt.Println(strings.Count("fivevev", "ve"))
	fmt.Printf("Fields: %q\n", strings.Fields("  foo bar  baz   "))
	fmt.Printf("FieldsFunc: %q\n", strings.FieldsFunc("  foo bar  baz   ", unicode.IsSpace))
	fmt.Printf("Split: %q\n", strings.Split("foo,bar,baz", ","))
	fmt.Printf("SplitAfter: %q\n", strings.SplitAfter("foo,bar,baz", ",")) //保留sep(,)
	fmt.Printf("SplitN: %q\n", strings.SplitN("foo,bar,baz", ",", 2))
	fmt.Println(strings.HasPrefix("ross", "r"))
	fmt.Printf("%d\n", strings.IndexFunc("studygolang", func(r rune) bool {
		if r > 'u' { //Unicode码点
			return true
		}
		return false
	}))

	fmt.Println(strings.Join([]string{"hello", "world", "yes"}, "&"))
	fmt.Println(Join([]string{"hello", "world", "yes"}, "&"))
	fmt.Println(strings.Repeat("ross ", 3))

	//替换
	fmt.Println(strings.Replace("oink oink oink", "oink", "moo", -1))
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
	fmt.Println(r.Replace("This is <b>HTML</b>!"))

}
