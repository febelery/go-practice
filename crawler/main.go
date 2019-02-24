package main

import (
	"learn/crawler/engine"
	"learn/crawler/zhenai/parser"
)

func main() {
	engine.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParserCityList,
	})
}
