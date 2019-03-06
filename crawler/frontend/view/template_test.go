package view

import (
	"learn/crawler/engine"
	page "learn/crawler/frontend/model"
	"learn/crawler/model"
	"os"
	"testing"
)

func TestSearchResultView_Render(t *testing.T) {
	view := CreateSearchResultView("template.html")
	out, err := os.Create("search_result.test.html")

	page := page.SearchResult{}
	page.Hits = 123
	item := engine.Item{
		Url:  "http://www.baidu.com",
		Type: "zhenai",
		Id:   "123456",
		Payload: model.Profile{
			Name:       "Rachel",
			Gender:     "女",
			Age:        18,
			Height:     168,
			Weight:     115,
			Income:     "3001-9999",
			Marriage:   "未婚",
			Education:  "博士",
			Occupation: "科研",
			Hokou:      "汉",
			Xinzou:     "白羊座",
			House:      "已购房",
			Car:        "已购车",
		},
	}

	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}

	err = view.Render(out, page)
	if err != nil {
		panic(err)
	}
}
