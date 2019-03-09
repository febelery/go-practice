package parser

import (
	"learn/crawler/engine"
	"learn/crawler_distributed/config"
	"regexp"
)

const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParseCityList(contents []byte, _ string) engine.ParserResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1)

	result := engine.ParserResult{}

	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: engine.NewFuncParser(ParseCity, config.ParseCity),
		})
	}

	return result
}
