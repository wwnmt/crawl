package main

import (
	"crawl/engine"
	"crawl/parse"
)

func main() {
	//var simpleEngine engine.SimpleEngine
	//simpleEngine.Run(engine.Request{
	//	Url:       "https://book.douban.com/tag",
	//	ParseFunc: parse.TagParser,
	//})

	e := engine.ConcurrentEngine{
		Scheduler:   &engine.SimpleScheduler{},
		WorkerCount: 10,
	}

	e.Run(engine.Request{
		Url:       "https://book.douban.com/tag",
		ParseFunc: parse.TagParser,
	})
}
