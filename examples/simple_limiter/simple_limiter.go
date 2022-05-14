package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	intervalLimiter()
	//intervalTicker2()
	//burstyLimiter()
	//ticker()

}

func intervalLimiter() {

	// serves req at fixed intervals

	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	t := time.NewTicker(1 * time.Second)

	for req := range requests {
		<-t.C // will block until the next tick from the timer
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

func intervalTicker2() {

	inChannel := make(chan bool)

	// limiter that limits the system to 1 task per second
	go func() {
		t := time.NewTicker(time.Second)
		for ; true; <-t.C {
			<-inChannel
			fmt.Println("TASK RECEIVED AT:", time.Now())
		}
	}()

	var wg sync.WaitGroup

	// start 5 concurrent tasks
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			inChannel <- true
			defer wg.Done()
		}()
	}

	wg.Wait()

}
