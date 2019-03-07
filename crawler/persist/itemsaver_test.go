package persist

import (
	"context"
	"encoding/json"
	"gopkg.in/olivere/elastic.v6"
	"learn/crawler/engine"
	"learn/crawler/model"
	"testing"
)

func TestSave(t *testing.T) {
	expected := engine.Item{
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

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	const index = "dating_test"
	err = Save(client, expected, index)
	if err != nil {
		panic(err)
	}

	resp, err := client.Get().Index(index).Type(expected.Type).Id(expected.Id).Do(context.Background())
	if err != nil {
		panic(err)
	}

	var actual engine.Item
	json.Unmarshal(*resp.Source, &actual)

	actualProfile, _ := model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile

	if actual != expected {
		t.Errorf("got %v; expected %v", actual, expected)
	}

}
