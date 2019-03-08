package main

import (
	"learn/crawler/engine"
	"learn/crawler/model"
	"learn/crawler_distributed/rpcsupport"
	"testing"
	"time"
)

func TestItemSaver(t *testing.T) {
	const host = ":1234"

	go serveRpc(host, "ross")
	time.Sleep(time.Second)

	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(client)
	}

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
			Birthplace: "北京",
			Hokou:      "汉",
			Shape:      "苗条",
			House:      "已购房",
			Car:        "已购车",
		},
	}

	result := ""
	client.Call("ItemSaverService.Save", item, &result)

	if err != nil || result != "ok" {
		t.Errorf("result: %s; err: %s", result, err)
	}
}
