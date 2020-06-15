package engine

import (
	"crawl/fetcher"
	"fmt"
	"log"
)

type Scheduler interface {
	Submit(request Request)
	configWorkChan(chan Request)
}

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}

type SimpleScheduler struct {
	WorkChan chan Request
}

func (s *SimpleScheduler) Submit(request Request) {
	go func() { s.WorkChan <- request }()
}

func (s *SimpleScheduler) configWorkChan(c chan Request) {
	s.WorkChan = c
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	in := make(chan Request)
	out := make(chan ParseResult)

	for i := 1; i <= e.WorkerCount; i++ {
		createWorker(in, out, i)
	}
	e.Scheduler.configWorkChan(in)

	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	itemCount := 1
	for {
		result := <-out
		for _, item := range result.Items {
			log.Printf("[Got item] %d. %s\n", itemCount, item)
			itemCount++
		}

		for _, url := range result.Requests {
			e.Scheduler.Submit(url)
		}
	}
}

func createWorker(in chan Request, out chan ParseResult, num int) {
	go func() {
		for {
			request := <-in
			result, err := work(request, num)

			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

func work(request Request, num int) (ParseResult, error) {
	fmt.Printf("No.%d Fetch URL: %s\n", num, request.Url)
	body, err := fetcher.Fetch(request.Url)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return ParseResult{}, err
	}
	return request.ParseFunc(body), nil
}
