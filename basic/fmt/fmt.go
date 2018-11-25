package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

type Website struct {
	Name string
}

type Person struct {
	Name string
	Age  int
	Sex  int
}

var site = Website{Name: "learn golang"}

var (
	name string
	age  int
)

func fmtPrint() {
	//普通占位符号
	fmt.Printf("%v\n", site)
	fmt.Printf("%+v\n", site)
	fmt.Printf("%#v\n", site)
	fmt.Printf("%T\n", site)
	fmt.Printf("%%\n")

	//布尔占位符
	fmt.Printf("%t\n", true)

	//整数占位符
	fmt.Printf("二进制:%b\n", 5)
	fmt.Printf("十进制:%d\n", 5)
	fmt.Printf("十六进制:%x\n", 12)
	fmt.Printf("相应Unicode码点所表示的字符:%c\n", 0x4E2D)
	fmt.Printf("%q\n", 0x4E2D)

	//浮点数和复数
	fmt.Printf("科学计数法:%e\n", 10.2)
	fmt.Printf("有小数点而无指数:%f\n", 10.2)
	fmt.Printf("根据情况选择 %%e 或 %%f 以产生更紧凑的（无末尾的0）输出:%g\n", 10.2)
	fmt.Printf("%.3f\n", 8413.1272)

	//字符串与字节切片
	fmt.Printf("%s\n", []byte("学习Golang"))
	fmt.Printf("%s\n", "学习Golang")

	//指针
	fmt.Printf("%p\n", &site)
}

func (this *Person) String() string {
	buffer := bytes.NewBufferString("This is ")
	buffer.WriteString(this.Name + ", ")
	if this.Sex == 0 {
		buffer.WriteString("He ")
	} else {
		buffer.WriteString("She ")
	}
	buffer.WriteString("is ")
	buffer.WriteString(strconv.Itoa(this.Age))
	buffer.WriteString(" years old.")
	return buffer.String()
}

func (this *Person) Format(f fmt.State, c rune) {
	if c == 'L' {
		f.Write([]byte(this.String()))
		f.Write([]byte(" Person has three fields.\n"))
	} else {
		f.Write([]byte(fmt.Sprintln(this.String())))
	}
}

func main() {
	fmtPrint()

	p := &Person{"Ross", 12, 0}
	fmt.Println(p)
	fmt.Printf("%L", p)
	fmt.Println(reflect.TypeOf(p),reflect.ValueOf(p).Kind())

	n, _ := fmt.Sscan("Ross 18", &name, &age)
	fmt.Println(n, name, age)

}
