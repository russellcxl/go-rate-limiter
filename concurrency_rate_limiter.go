package go_rate_limiter

// NewMaxConcurrencyRateLimiter returns a max concurrency rate limiter
func NewMaxConcurrencyRateLimiter(conf *Config) (RateLimiter, error) {
	if conf.Limit <= 0 {
		return nil, ErrInvalidLimit
	}

	m := NewManager(conf)

	// cronjob to reset tokens every X seconds in case users forget to release them
	m.runResetTokenTask(conf.TokenResetsAfter)

	go func() {
		for {
			select {
			case <-m.inChan:
				m.tryGenerateToken()
			case t := <-m.releaseChan:
				m.releaseToken(t)
			}
		}
	}()

	return m, nil
}
