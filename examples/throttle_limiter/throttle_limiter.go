package main

import (
	"fmt"
	ratelimiter "github.com/russellcxl/go-rate-limiter"
	"math/rand"
	"sync"
	"time"
)

func main() {
	r, err := ratelimiter.NewThrottleRateLimiter(&ratelimiter.Config{
		Throttle: 1 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	doWork := func(idx int) {
		// Acquire a rate limit token
		token, err := r.Acquire()
		fmt.Printf("Worker %d acquired token %s at %s...\n", idx, token.ID, time.Now().UTC())
		if err != nil {
			panic(err)
		}
		// Simulate some other work; takes between 0-5 seconds
		n := rand.Intn(5)
		fmt.Printf("Worker %d working for %d seconds...\n", idx, n)
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Printf("Worker %d done\n", idx)
		defer wg.Done()
	}

	// Spin up 10 workers that need a rate limit resource
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doWork(i)
	}

	wg.Wait()
}

