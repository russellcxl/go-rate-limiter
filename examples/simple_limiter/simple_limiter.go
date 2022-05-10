package main

import (
	"fmt"
	"time"
)

func main() {

	//intervalLimiter()
	//burstyLimiter()
	//ticker()
	ticker2()

}

func intervalLimiter() {

	// serves req at fixed intervals

	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(2 * time.Second)

	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}
}

func burstyLimiter() {

	// serves 3 reqs at a time, then subsequent ones at intervals

	limiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		limiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(2 * time.Second) {
			limiter <- t
		}
	}()

	// creates 5 concurrent requests
	reqs := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		reqs <- i
	}
	close(reqs)

	for req := range reqs {
		l := <-limiter // will read 3, then block until there is an input from the ticker
		fmt.Println("request", req, l)
	}
}

func ticker() {
	// every 1 second, the current time will be sent into t.C channel
	t := time.NewTicker(time.Second)
	defer t.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Done!") // after 10s, this will run and stop the program
			return
		case timeFromTicker := <-t.C:
			fmt.Println("Current time: ", timeFromTicker)
		}
	}
}

func ticker2() {
	t := time.NewTicker(time.Second)

	for ; true; <-t.C {
		fmt.Println(time.Now())
	}
}