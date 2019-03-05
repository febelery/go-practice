package parser

import (
	"learn/crawler/engine"
	"learn/crawler/model"
	"regexp"
	"strconv"
)

var profileRe = regexp.MustCompile(`<div class="des f-cl" data-v-3c42fade>([^\|]{1,10}) \| ([\d]{1,3})岁 \| ([^\|]{1,10}) \| ([^\|]{1,10}) \| ([\d]{2,3})cm \| ([^<]{1,14})</div>`)
var idUrlRe = regexp.MustCompile(`http://album.zhenai.com/u/[\d]+`)

func ParserProfile(contents []byte, url string, name string) engine.ParserResult {
	profile := model.Profile{}

	profile.Name = name

	match := profileRe.FindSubmatch(contents)
	if len(match) >= 2 {
		profile.Age, _ = strconv.Atoi(string(match[2]))
		profile.Education = string(match[3])
		profile.Marriage = string(match[4])
		profile.Height, _ = strconv.Atoi(string(match[5]))
		profile.Income = string(match[6])
	}

	id := ""
	match = profileRe.FindSubmatch(contents)
	if len(match) >= 2 {
		id = string(match[1])
	}

	result := engine.ParserResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "zhenai",
				Id:      id,
				Payload: profile,
			},
		},
	}

	return result
}
