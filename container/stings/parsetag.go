package main

import (
	"fmt"
	"github.com/fatih/structtag"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	parseTag()
	parseStructTag()
}

func parseTag() {
	tag := `json:"foo,omitempty,string" xml:"foo"`

	tags, err := structtag.Parse(string(tag))
	if err != nil {
		panic(err)
	}

	for _, t := range tags.Tags() {
		fmt.Printf("tag: %+v\n", t)
	}

	jsonTag, err := tags.Get("json")
	if err != nil {
		panic(err)
	}

	// change exitsting tag
	jsonTag.Name = "foo_bar"
	jsonTag.Options = nil
	tags.Set(jsonTag)

	tags.Set(&structtag.Tag{
		Key:     "hcl",
		Name:    "foo",
		Options: []string{"squash"},
	})

	fmt.Println(tags)
}

func parseStructTag() {
	src := `package main 
			type Example struct {Foo string` + "`json:\"foo\"`}"

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "demo", src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(file, func(x ast.Node) bool {
		s, ok := x.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range s.Fields.List {
			fmt.Printf("Field: %s\n", field.Names[0].Name)
			fmt.Printf("Tag:   %s\n", field.Tag.Value)
		}

		return false
	})
}
