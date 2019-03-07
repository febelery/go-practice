package parser

import (
	"learn/crawler/engine"
	"learn/crawler/model"
	"regexp"
	"strconv"
)

var profileRe = regexp.MustCompile(`<div class="des f-cl" data-v-[^>]+>([^\|]{1,10}) \| ([\d]{1,3})岁 \| ([^\|]{1,10}) \| ([^\|]{1,10}) \| ([\d]{2,3})cm \| ([^<]{1,14})</div>`)
var infoRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>([^<]{1,20})</div>`)
var idUrlRe = regexp.MustCompile(`http://album.zhenai.com/u/[\d]+`)
var hokouRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>([^<]{1,20}族)</div>`)
var carRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>([^<]{1,20}车)</div>`)
var houseRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>([^<]{1,20}房)</div>`)
var birthplaceRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>籍贯:([^<]+)</div>`)
var shapeRe = regexp.MustCompile(`<div class="m-btn pink" data-v-[^>]+>体型:([^<]+)</div>`)
var maleRe = regexp.MustCompile(`分享他`)
var femaleRe = regexp.MustCompile(`分享她`)

func ParserProfile(contents []byte, url string, name string) engine.ParserResult {
	profile := model.Profile{}
	profile.Name = name

	if male := maleRe.FindSubmatch(contents); len(male) >= 1 {
		profile.Gender = "男"
	}

	if female := femaleRe.FindSubmatch(contents); len(female) >= 1 {
		profile.Gender = "女"
	}

	if hokou := hokouRe.FindSubmatch(contents); len(hokou) >= 2 {
		profile.Hokou = string(hokou[1])
	}

	if car := carRe.FindSubmatch(contents); len(car) >= 2 {
		profile.Car = string(car[1])
	}

	if house := houseRe.FindSubmatch(contents); len(house) >= 2 {
		profile.House = string(house[1])
	}

	if birthplace := birthplaceRe.FindSubmatch(contents); len(birthplace) >= 2 {
		profile.Birthplace = string(birthplace[1])
	}

	if shape := shapeRe.FindSubmatch(contents); len(shape) >= 2 {
		profile.Shape = string(shape[1])
	}

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
