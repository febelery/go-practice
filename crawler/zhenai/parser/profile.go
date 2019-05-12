package parser

import (
	"learn/crawler/engine"
	"learn/crawler/model"
	"learn/crawler_distributed/config"
	"regexp"
	"strconv"
)

var profileRe = regexp.MustCompile(`<div class="des interceptor-cl" data-v-[^>]+>([^\|]{1,10}) \| ([\d]{1,3})岁 \| ([^\|]{1,10}) \| ([^\|]{1,10}) \| ([\d]{2,3})cm \| ([^<]{1,14})</div>`)
var infoRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>([^<]{1,20})</div>`)
var idUrlRe = regexp.MustCompile(`http://album.zhenai.com/u/[\d]+`)
var hokouRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>([^<]{1,20}族)</div>`)
var carRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>([^<]{1,20}车)</div>`)
var houseRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>([^<]{1,20}房)</div>`)
var birthplaceRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>籍贯:([^<]+)</div>`)
var shapeRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>体型:([^<]+)</div>`)
var maleRe = regexp.MustCompile(`分享他`)
var femaleRe = regexp.MustCompile(`分享她`)

func parseProfile(contents []byte, url string, name string) engine.ParserResult {
	profile := model.Profile{}
	profile.Name = name

	if male := maleRe.FindSubmatch(contents); len(male) >= 1 {
		profile.Gender = "男"
	}

	if female := femaleRe.FindSubmatch(contents); len(female) >= 1 {
		profile.Gender = "女"
	}

	profile.Hokou = extractString(contents, hokouRe)
	profile.Car = extractString(contents, carRe)
	profile.House = extractString(contents, houseRe)
	profile.Birthplace = extractString(contents, birthplaceRe)
	profile.Shape = extractString(contents, shapeRe)

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

	//matches := guessRe.FindAllSubmatch(contents, -1)
	//for _, m := range matches {
	//	result.Requests = append(result.Requests, engine.Request{
	//		Url:    string(m[1]),
	//		Parser: NewProfileParser(string(m[2])),
	//	})
	//}

	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)

	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

type ProfileParser struct {
	userName string
}

func (p *ProfileParser) Parse(contents []byte, url string) engine.ParserResult {
	return parseProfile(contents, url, p.userName)
}

func (p *ProfileParser) Serialize() (name string, args interface{}) {
	return config.ParseProfile, p.userName
}

func NewProfileParser(name string) *ProfileParser {
	return &ProfileParser{
		userName: name,
	}
}
