package go_rate_limiter

import (
	"time"
)

// RateLimiter exposes Acquire() for obtaining a Rate Limit Token
type RateLimiter interface {
	Acquire() (*Token, error)
}

type Config struct {
	// Throttle is the min time between requests for a Throttle Rate Limiter
	Throttle time.Duration
}

