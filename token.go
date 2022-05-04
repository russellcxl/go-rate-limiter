package go_rate_limiter

import (
	"github.com/segmentio/ksuid"
	"time"
)

// token factory function creates a new token
type tokenFactory func() *Token

// Token represents a Rate Limit Token
type Token struct {
	ID        string
	CreatedAt time.Time
	ExpiresAt time.Time // Defines the min amount of time the token must live before being released
}

// NewToken creates a new token
func NewToken() *Token {
	return &Token{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Time{}, // defaults to zero time
	}
}

// IsExpired returns true if current time is greater than expiration time
func (t *Token) IsExpired() bool {
	now := time.Now().UTC()
	return t.ExpiresAt.Before(now)
}

// NeedReset returns true if elapsed time since token was created
// is greater than provided reset duration
func (t *Token) NeedReset(resetAfter time.Duration) bool {
	if time.Since(t.CreatedAt) >= resetAfter {
		return true
	}
	return false
}
