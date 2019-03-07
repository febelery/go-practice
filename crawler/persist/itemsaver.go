package persist

import (
	"context"
	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v6"
	"learn/crawler/engine"
	"log"
)

func ItemSaver(index string) (chan engine.Item, error) {
	out := make(chan engine.Item)
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	go func() {
		itemCount := 0
		for {
			item := <-out
			itemCount++
			log.Printf("Item Saver: got item #%d: %v", itemCount, item)

			err := Save(client, item, index)
			if err != nil {
				log.Printf("Item saver: error saving item %v: %v", item, err)
			}
		}
	}()
	return out, nil
}

func Save(client *elastic.Client, item engine.Item, index string) error {
	if item.Type == "" {
		return errors.New("must supply Type")
	}

	indexService := client.Index().
		Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}

	_, err := indexService.Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}
