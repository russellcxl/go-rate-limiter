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
	})

	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	doWork := func(id int) {
		// Acquire a rate limit token
		token, err := r.Acquire()
		fmt.Printf("Rate Limit Token %s acquired at %s...\n", token.ID, time.Now().UTC())
		if err != nil {
			panic(err)
		}
		// Simulate some work
		n := rand.Intn(5)
		fmt.Printf("Worker %d Sleeping %d seconds...\n", id, n)
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Printf("Worker %d Done\n", id)
		r.Release(token)
		wg.Done()
	}

	// Spin up a 10 workers that need a rate limit resource
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go doWork(i)
	}

	wg.Wait()
}