package engine

import (
	"crawl/fetcher"
	"log"
)

type SimpleEngine struct {
}

func (e *SimpleEngine) Run(seeds ...Request) {
	var requests []Request

	for _, e := range seeds {
		requests = append(requests, e)
	}

	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]
		log.Printf("Fetching url:%s\n", r.Url)
		body, err := fetcher.Fetch(r.Url)
		if err != nil {
			log.Println("Error: ", err)
			continue
		}

		parseResult := r.ParseFunc(body)

		requests = append(requests, parseResult.Requests...)

		for _, item := range parseResult.Items {
			log.Printf("[Got item] %s\n", item)
		}
	}
}
