package persist

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"learn/crawler/model"
	"testing"
)

func TestSave(t *testing.T) {
	expected := model.Profile{
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
	}

	id, err := save(expected)
	if err != nil {
		panic(err)
	}

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	resp, err := client.Get().Index("dating_profile").Type("zhenai").Id(id).Do(context.Background())
	if err != nil {
		panic(err)
	}

	var actual model.Profile
	err = json.Unmarshal(*resp.Source, &actual)
	if err != nil {
		panic(err)
	}

	if actual != expected {
		t.Errorf("got %v; expected %v", actual, expected)
	}

}
