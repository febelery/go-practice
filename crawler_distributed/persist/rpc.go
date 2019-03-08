package persist

import (
	"gopkg.in/olivere/elastic.v6"
	"learn/crawler/engine"
	"learn/crawler/persist"
	"log"
)

type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

func (s *ItemSaverService) Save(item engine.Item, result *string) error {
	err := persist.Save(s.Client, item, s.Index)
	log.Printf("Item %v saved.", item)
	if err == nil {
		*result = "ok"
	} else {
		log.Printf("Error saving item %v: %v.", item, err)
	}
	return err
}
