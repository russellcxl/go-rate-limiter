package go_rate_limiter

import (
	"errors"
	"time"
)

// NewThrottleRateLimiter returns a throttle rate limiter
func NewThrottleRateLimiter(conf *Config) (RateLimiter, error) {
	if conf.Throttle == 0 {
		return nil, errors.New("Throttle duration must be greater than zero")
	}

	m := NewManager(conf)

	ticker := time.NewTicker(conf.Throttle)

	// starts the limiter in the background
	go func() {
		// loops the ticker channel every {throttle duration}
		for ; true; <-ticker.C {
			<-m.inChan           // this will unblock when Acquire() is called
			m.tryGenerateToken() // this will send a token to the outChan
		}
	}()

	return m, nil
}
