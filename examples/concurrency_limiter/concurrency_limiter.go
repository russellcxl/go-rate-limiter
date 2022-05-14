package main

import (
	"fmt"
	ratelimiter "github.com/russellcxl/go-rate-limiter"
	"math/rand"
	"sync"
	"time"
)

func main() {
	r, err := ratelimiter.NewMaxConcurrencyRateLimiter(&ratelimiter.Config{
		Limit: 3,
		TokenResetsAfter: 10 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	doWork := func(id int) {
		// Acquire a rate limit token
		token, err := r.Acquire()
		fmt.Printf("Worker %d acquired token %s at %s...\n", id, token.ID, time.Now().UTC())
		if err != nil {
			panic(err)
		}
		// Simulate some work
		n := rand.Intn(5)
		fmt.Printf("Worker %d working for %d seconds...\n", id, n)
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Printf("Worker %d done\n", id)

		// release the token once the worker is done; can simulate case where worker forgets to release token
		//r.Release(token)

		wg.Done()
	}

	// Spin up a 10 workers that need a rate limit resource
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doWork(i)
	}

	wg.Wait()
}