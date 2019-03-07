package engine

import (
	"learn/crawler/fetcher"
	"log"
)

func worker(r Request) (ParserResult, error) {
	log.Printf("Fetching %s", r.Url)
	body, err := fetcher.Fetcher(r.Url)
	if err != nil {
		log.Printf("Fetcher: error fetching url %s: %v", r.Url, err)
		return ParserResult{}, err
	}

	return r.Parser.Parse(body, r.Url), nil
}
