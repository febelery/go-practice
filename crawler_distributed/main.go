package main

import (
	"fmt"
	"learn/crawler/engine"
	"learn/crawler/scheduler"
	"learn/crawler/zhenai/parser"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/persist/client"
)

func main() {
	itemChan, err := client.ItemSaver(fmt.Sprintf(":%d", config.ItemSaverPort))
	if err != nil {
		panic(err)
	}

	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		ItemChan:    itemChan,
	}

	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun/shanghai",
		Parser: engine.NewFuncParser(parser.ParseCityList, "ParseCityList"),
	})
}
