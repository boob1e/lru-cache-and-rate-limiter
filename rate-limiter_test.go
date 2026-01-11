package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRateLimiter_InitialState(t *testing.T) {
	fmt.Println("\n### Test: Initial State ###")
	// Create rate limiter: 10 tokens/sec, burst of 5
	limiter := NewRateLimiter(10.0, 5.0)

	fmt.Printf("Capacity: %.1f, Initial Tokens: %.1f, RefillRate: %.1f/sec\n",
		limiter.Capacity, limiter.Tokens, limiter.RefillRate)

	if limiter.Tokens != 5.0 {
		t.Errorf("Expected initial tokens to be 5.0, got %.1f", limiter.Tokens)
	}

	if limiter.Capacity != 5.0 {
		t.Errorf("Expected capacity to be 5.0, got %.1f", limiter.Capacity)
	}

	fmt.Println("✓ Rate limiter initialized with full burst capacity")
}

func TestRateLimiter_SingleRequest(t *testing.T) {
	fmt.Println("\n### Test: Single Request ###")
	limiter := NewRateLimiter(10.0, 5.0)

	allowed := limiter.Allow()
	fmt.Printf("First request allowed: %v, Remaining tokens: %.1f\n", allowed, limiter.Tokens)

	if !allowed {
		t.Error("Expected first request to be allowed")
	}

	if limiter.Tokens != 4.0 {
		t.Errorf("Expected 4 tokens remaining, got %.1f", limiter.Tokens)
	}

	fmt.Println("✓ Single request consumed 1 token")
}

func TestRateLimiter_BurstCapacity(t *testing.T) {
	fmt.Println("\n### Test: Burst Capacity ###")
	limiter := NewRateLimiter(10.0, 5.0)

	fmt.Println("Attempting 5 rapid requests (burst capacity)...")
	successCount := 0
	for i := 0; i < 5; i++ {
		if limiter.Allow() {
			successCount++
		}
		fmt.Printf("  Request %d: allowed, tokens remaining: %.1f\n", i+1, limiter.Tokens)
	}

	if successCount != 5 {
		t.Errorf("Expected all 5 burst requests to succeed, got %d", successCount)
	}

	// Use tolerance for floating point comparison
	if limiter.Tokens > 0.1 {
		t.Errorf("Expected ~0 tokens remaining after burst, got %.2f", limiter.Tokens)
	}

	fmt.Println("✓ All burst requests allowed, tokens exhausted")
}

func TestRateLimiter_RateLimitingAfterBurst(t *testing.T) {
	fmt.Println("\n### Test: Rate Limiting After Burst ###")
	limiter := NewRateLimiter(10.0, 5.0)

	// Exhaust burst capacity
	for i := 0; i < 5; i++ {
		limiter.Allow()
	}
	fmt.Println("Burst capacity exhausted (5 requests)")

	// Try one more immediate request
	allowed := limiter.Allow()
	fmt.Printf("6th immediate request allowed: %v, tokens: %.1f\n", allowed, limiter.Tokens)

	if allowed {
		t.Error("Expected 6th immediate request to be blocked after burst exhaustion")
	}

	fmt.Println("✓ Request correctly blocked when tokens exhausted")
}

func TestRateLimiter_TokenRefill(t *testing.T) {
	fmt.Println("\n### Test: Token Refill Over Time ###")
	// 10 tokens per second = 1 token every 100ms
	limiter := NewRateLimiter(10.0, 5.0)

	// Exhaust all tokens
	for i := 0; i < 5; i++ {
		limiter.Allow()
	}
	fmt.Println("Exhausted all 5 tokens")

	// Immediate request should fail
	if limiter.Allow() {
		t.Error("Request should be blocked immediately after exhaustion")
	}

	// Wait 150ms (should refill ~1.5 tokens, enough for 1 request)
	fmt.Println("Waiting 150ms for refill...")
	time.Sleep(150 * time.Millisecond)

	allowed := limiter.Allow()
	fmt.Printf("Request after 150ms: allowed=%v, tokens=%.2f\n", allowed, limiter.Tokens)

	if !allowed {
		t.Error("Expected request to be allowed after refill period")
	}

	fmt.Println("✓ Tokens refilled over time, request allowed")
}

func TestRateLimiter_PartialRefill(t *testing.T) {
	fmt.Println("\n### Test: Partial Token Refill ###")
	// 10 tokens/sec = 1 token per 100ms
	limiter := NewRateLimiter(10.0, 5.0)

	// Use all tokens
	for i := 0; i < 5; i++ {
		limiter.Allow()
	}

	// Wait 250ms = should refill 2.5 tokens
	fmt.Println("Waiting 250ms (should refill ~2.5 tokens)...")
	time.Sleep(250 * time.Millisecond)

	// Should allow 2 requests (2.5 tokens available)
	successCount := 0
	for i := 0; i < 3; i++ {
		if limiter.Allow() {
			successCount++
		}
	}

	fmt.Printf("Allowed %d out of 3 requests\n", successCount)

	if successCount != 2 {
		t.Errorf("Expected 2 requests to succeed with 2.5 tokens, got %d", successCount)
	}

	fmt.Println("✓ Partial token refill working correctly")
}

func TestRateLimiter_CapacityLimit(t *testing.T) {
	fmt.Println("\n### Test: Token Capacity Cap ###")
	limiter := NewRateLimiter(10.0, 5.0)

	// Wait longer than needed to fill capacity
	// 5 tokens at 10/sec = 0.5 seconds to fully refill
	fmt.Println("Waiting 1 second (more than needed to refill capacity)...")
	time.Sleep(1 * time.Second)

	// Force a refill check by calling Allow
	limiter.Allow()

	// Should not exceed capacity of 5 (minus 1 consumed = 4 remaining)
	fmt.Printf("Tokens after long wait: %.1f (capacity: %.1f)\n", limiter.Tokens, limiter.Capacity)

	if limiter.Tokens > limiter.Capacity {
		t.Errorf("Tokens (%.1f) exceeded capacity (%.1f)", limiter.Tokens, limiter.Capacity)
	}

	if limiter.Tokens != 4.0 {
		t.Errorf("Expected 4 tokens after refill and 1 consumption, got %.1f", limiter.Tokens)
	}

	fmt.Println("✓ Token count capped at capacity")
}

func TestRateLimiter_ConcurrentRequests(t *testing.T) {
	fmt.Println("\n### Test: Concurrent Requests (Thread Safety) ###")
	limiter := NewRateLimiter(100.0, 10.0)

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	// Launch 20 concurrent goroutines
	numGoroutines := 20
	fmt.Printf("Launching %d concurrent requests (burst capacity: 10)...\n", numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if limiter.Allow() {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("Successful requests: %d out of %d\n", successCount, numGoroutines)
	fmt.Printf("Remaining tokens: %.1f\n", limiter.Tokens)

	// Should allow exactly 10 (the burst capacity)
	if successCount != 10 {
		t.Errorf("Expected 10 successful requests, got %d", successCount)
	}

	fmt.Println("✓ Concurrent requests handled safely with mutex")
}

func TestRateLimiter_SustainedRate(t *testing.T) {
	fmt.Println("\n### Test: Sustained Rate Limiting ###")
	// 5 tokens/sec, burst of 2
	limiter := NewRateLimiter(5.0, 2.0)

	// Exhaust burst
	limiter.Allow()
	limiter.Allow()
	fmt.Println("Exhausted burst capacity (2 tokens)")

	// Make requests at exactly the refill rate (200ms intervals = 5/sec)
	fmt.Println("Making requests every 200ms (matching refill rate)...")
	successCount := 0
	for i := 0; i < 5; i++ {
		time.Sleep(200 * time.Millisecond)
		if limiter.Allow() {
			successCount++
		}
	}

	fmt.Printf("Successful requests at sustained rate: %d/5\n", successCount)

	// Should allow all requests since we're matching the refill rate
	if successCount != 5 {
		t.Errorf("Expected all 5 sustained-rate requests to succeed, got %d", successCount)
	}

	fmt.Println("✓ Sustained rate limiting works correctly")
}
