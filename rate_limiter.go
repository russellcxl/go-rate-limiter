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

	// Throttle Await Function
	await := func(throttle time.Duration) {
		ticker := time.NewTicker(throttle)
		go func() {
			// loops ticker channel every {throttle duration}
			for ; true; <-ticker.C {
				<-m.inChan // this will unblock when Acquire() is called
				m.tryGenerateToken() // sends a token to the outChan
			}
		}()
	}

	// Call await to start
	await(conf.Throttle)
	return m, nil
}
