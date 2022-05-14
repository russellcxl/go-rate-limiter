package go_rate_limiter

import (
	"time"
)

// RateLimiter exposes Acquire() for obtaining a Rate Limit Token
type RateLimiter interface {
	Acquire() (*Token, error)
	Release(*Token)
}

type Config struct {
	// Limit determines how many rate limit tokens can be active at a time
	Limit int

	// Throttle is the min time between requests for a Throttle Rate Limiter
	Throttle time.Duration

	// TokenResetsAfter is the maximum amount of time a token can live before being
	// forcefully released - if set to zero time then the token may live forever
	TokenResetsAfter time.Duration
}

