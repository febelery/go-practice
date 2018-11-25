package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func byteRWerExample() {
FOREND:
	for {
		fmt.Println("请输入要通过WriteByte写入的一个ASCII字符（b：返回上级；q：退出）：")
		var ch byte
		fmt.Scanf("%c\n", &ch)
		switch ch {
		case 'b':
			fmt.Println("返回上级菜单！")
			break FOREND
		case 'q':
			fmt.Println("程序退出！")
			os.Exit(0)
		default:
			buffer := new(bytes.Buffer)
			err := buffer.WriteByte(ch)
			if err == nil {
				fmt.Println("写入一个字节成功！准备读取该字节……")
				newCh, _ := buffer.ReadByte()
				fmt.Printf("读取的字节：%c\n", newCh)
			} else {
				fmt.Println("写入错误")
			}
		}
	}
}

func readerExample() {
FOREND:
	for {
		fmt.Println("")
		fmt.Println("*******从不同来源读取数据*********")
		fmt.Println("*******请选择数据源，请输入：*********")
		fmt.Println("1 表示 标准输入")
		fmt.Println("2 表示 普通文件")
		fmt.Println("3 表示 从字符串")
		fmt.Println("4 表示 从网络")
		fmt.Println("b 返回上级菜单")
		fmt.Println("q 退出")
		fmt.Println("***********************************")

		var ch string
		fmt.Scanln(&ch)
		var (
			data []byte
			err  error
		)
		switch strings.ToLower(ch) {
		case "1":
			fmt.Println("请输入不多于9个字符，以回车结束：")
			data, err = readFrom(os.Stdin, 11)
		case "2":
			file, err := os.Open("io/01.txt")
			if err != nil {
				fmt.Println("打开文件 01.txt 错误:", err)
				continue
			}
			data, err = readFrom(file, 9)
			file.Close()
		case "3":
			data, err = readFrom(strings.NewReader("from string"), 12)
		case "4":
			fmt.Println("暂未实现！")
		case "b":
			fmt.Println("返回上级菜单！")
			break FOREND
		case "q":
			fmt.Println("程序退出！")
			os.Exit(0)
		default:
			fmt.Println("输入错误！")
			continue
		}

		if err != nil {
			fmt.Println("数据读取失败，可以试试从其他输入源读取！")
		} else {
			fmt.Printf("读取到的数据是：%s\n", data)
		}

	}
}

func readFrom(reader io.Reader, num int) ([]byte, error) {
	p := make([]byte, num)
	n, err := reader.Read(p)
	if n > 0 {
		return p[:n], nil
	}
	return p, err
}

func main() {
MAINFOR:
	for {
		fmt.Println("")
		fmt.Println("*******请选择示例：*********")
		fmt.Println("1 表示 io.Reader 示例")
		fmt.Println("2 表示 io.ByteReader/ByteWriter 示例")
		fmt.Println("q 退出")
		fmt.Println("***********************************")

		var ch string
		fmt.Scanln(&ch)

		switch ch {
		case "1":
			readerExample()
		case "2":
			byteRWerExample()
		case "q":
			fmt.Println("程序退出！")
			break MAINFOR
		default:
			fmt.Println("输入错误！")
			continue
		}
	}
}
