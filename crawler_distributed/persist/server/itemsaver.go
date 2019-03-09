package main

import (
	"flag"
	"fmt"
	"gopkg.in/olivere/elastic.v6"
	"learn/crawler_distributed/config"
	"learn/crawler_distributed/persist"
	"learn/crawler_distributed/rpcsupport"
	"log"
)

var port = flag.Int("port", 1234, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	err := serveRpc(fmt.Sprintf(":%d", *port), config.ElasticIndex)
	log.Fatal(err)
}

func serveRpc(host, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{
		Client: client,
		Index:  index,
	})
}
