package go_rate_limiter

// Manager implements a rate limiter interface.
type Manager struct {
	errorChan    chan error
	outChan      chan *Token
	inChan       chan struct{}
	makeToken    tokenFactory
}

// NewManager creates a manager type
func NewManager(conf *Config) *Manager {
	m := &Manager{
		errorChan:    make(chan error),
		outChan:      make(chan *Token),
		inChan:       make(chan struct{}),
		makeToken:    NewToken,
	}
	return m
}

// Acquire is called to acquire a new token
func (m *Manager) Acquire() (*Token, error) {
	go func() {
		m.inChan <- struct{}{}
	}()

	// Await rate limit token (or error)
	select {
	case t := <-m.outChan:
		return t, nil
	case err := <-m.errorChan:
		return nil, err
	}
}

// Called when a new token is needed.
func (m *Manager) tryGenerateToken() {
	// panic if token factory is not defined
	if m.makeToken == nil {
		panic(ErrTokenFactoryNotDefined)
	}

	token := m.makeToken()

	// send token to outChan
	go func() {
		m.outChan <- token
	}()
}
