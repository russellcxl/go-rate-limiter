package go_rate_limiter

import (
	"log"
	"sync/atomic"
	"time"
)

// token factory function creates a new token
type tokenFactory func() *Token

// Manager implements a rate limiter interface.
type Manager struct {
	errorChan    chan error
	outChan      chan *Token
	inChan       chan bool
	releaseChan  chan *Token
	needToken    int64
	activeTokens map[string]*Token
	limit        int
	makeToken    tokenFactory
}

// NewManager creates a manager type
func NewManager(conf *Config) *Manager {
	m := &Manager{
		errorChan:    make(chan error),
		outChan:      make(chan *Token),
		inChan:       make(chan bool),
		activeTokens: make(map[string]*Token),
		releaseChan:  make(chan *Token),
		needToken:    0,
		limit:        conf.Limit,
		makeToken:    NewToken,
	}
	return m
}

// Acquire is called to acquire a new token
func (m *Manager) Acquire() (*Token, error) {
	go func() {
		m.inChan <- true
	}()

	// Blocks until token (or error) is received from
	select {
	case t := <-m.outChan:
		return t, nil
	case err := <-m.errorChan:
		return nil, err
	}
}

// Release is called to release an active token
func (m *Manager) Release(t *Token) {
	go func() {
		m.releaseChan <- t
	}()

}

func (m *Manager) incNeedToken() {
	atomic.AddInt64(&m.needToken, 1)
}

func (m *Manager) decNeedToken() {
	atomic.AddInt64(&m.needToken, -1)
}

func (m *Manager) awaitingToken() bool {
	return atomic.LoadInt64(&m.needToken) > 0
}

func (m *Manager) isLimitReached() bool {
	if len(m.activeTokens) >= m.limit {
		return true
	}
	return false
}

func (m *Manager) releaseToken(token *Token) {
	if token == nil {
		log.Print("unable to release nil token")
		return
	}

	if _, ok := m.activeTokens[token.ID]; !ok {
		log.Printf("unable to relase token %s - not in use\n", token)
		return
	}

	// Delete from map
	delete(m.activeTokens, token.ID)

	// process anything waiting for a rate limit
	if m.awaitingToken() {
		m.decNeedToken()
		go m.tryGenerateToken()
	}
}

func (m *Manager) tryGenerateToken() {
	// panic if token factory is not defined
	if m.makeToken == nil {
		panic(ErrTokenFactoryNotDefined)
	}

	// cannot continue if limit has been reached
	if m.isLimitReached() {
		m.incNeedToken()
		return
	}

	token := m.makeToken()

	// Add token to active map
	m.activeTokens[token.ID] = token

	// send token to outChan
	go func() {
		m.outChan <- token
	}()
}

// in case workers forget to release their token; cronjob checks for expired ones
func (m *Manager) runResetTokenTask(resetAfter time.Duration) {
	go func() {
		ticker := time.NewTicker(resetAfter)
		for; true; <-ticker.C {
			for _, token := range m.activeTokens {
				if token.NeedReset(resetAfter) {
					go func(t *Token) {
						m.releaseChan <- t
					}(token)
				}
			}
		}
	}()
}
