package parser

import (
	"learn/crawler/engine"
	"regexp"
)

var (
	cityProfileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
	cityUrlRe     = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/\w+/[^"]+)"`)
)

func ParserCity(contents []byte) engine.ParserResult {
	matches := cityProfileRe.FindAllSubmatch(contents, -1)

	result := engine.ParserResult{}
	for _, m := range matches {
		name := string(m[2])
		result.Requests = append(result.Requests, engine.Request{
			Url: string(m[1]),
			ParserFunc: func(contents []byte) engine.ParserResult {
				return ParserProfile(contents, name)
			},
		})
	}

	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		if len(m) > 1 {
			result.Requests = append(result.Requests, engine.Request{
				Url:        string(m[1]),
				ParserFunc: ParserCity,
			})
		}
	}

	return result
}
