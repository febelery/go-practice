package parser

import (
	"learn/crawler/engine"
	"regexp"
)

var (
	cityProfileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	cityUrlRe     = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/\w+/[^"]+)"`)
)

func ParseCity(contents []byte, _ string) engine.ParserResult {
	matches := cityProfileRe.FindAllSubmatch(contents, -1)

	result := engine.ParserResult{}
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: NewProfileParser(string(m[2])),
		})
	}

	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		if len(m) > 1 {
			result.Requests = append(result.Requests, engine.Request{
				Url:    string(m[1]),
				Parser: engine.NewFuncParser(ParseCity, "ParseCity"),
			})
		}
	}

	return result
}
