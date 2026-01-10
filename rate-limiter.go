package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu             sync.Mutex
	Capacity       float64
	Tokens         float64
	RefillRate     float64 // how many tokens per n seconds
	LastRefillTime time.Time
}

func NewRateLimiter(rate float64, burst float64) *RateLimiter {
	return &RateLimiter{
		Capacity:       burst,
		Tokens:         burst,
		RefillRate:     rate,
		LastRefillTime: time.Now(),
	}
}

func (l *RateLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	elapsedTime := time.Since(l.LastRefillTime)
	tokensToAdd := elapsedTime.Seconds() * l.RefillRate
	totalTokens := min(tokensToAdd+l.Tokens, l.Capacity)
	l.Tokens = totalTokens
	l.LastRefillTime = time.Now()

	if l.Tokens < 1 {
		return false
	}
	l.Tokens--
	return true
}
