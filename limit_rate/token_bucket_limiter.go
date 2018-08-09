/* token_bucket_limiter.go - a rate limiter using token bucket algorithm */
/*
modification history
--------------------
2016/12/01, by Sijie Yang, Create
*/
/*
DESCRIPTION
*/
package limit_rate

import (
	"sync"
	"time"
)

type TokenBucketLimiter struct {
	rate     int64 // tokens added per second
	capacity int64 // capacity of bucket
	amount   int64 // current amount of tokens in bucket
	last     int64 // timestamp of last try
	lock     sync.Mutex
}

/* NewTokenBucketLimiter - create a rate limiter
 *
 * Params:
 *     - ops  : maximum operation per second
 *     - burst: maximum burst number of operation
 *
 * Return:
 *     - rate limter
 */
func NewTokenBucketLimiter(ops int64, burst int64) *TokenBucketLimiter {
	l := new(TokenBucketLimiter)
	if ops <= 0 {
		ops = 1000 // default maximum operation per second
	}
	if burst <= 0 {
		burst = 1000 // default maximum burst number of operation
	}

	// Note: each operation will take 1000 tokens to bucket
	l.rate = ops * 1000
	l.capacity = burst * 1000
	l.last = time.Now().UnixNano() / int64(time.Millisecond)
	return l
}

/* Try - check whether operation is allowable or not
 *
 * Return:
 *     - ret: true if allowable, false if not
 */
func (l *TokenBucketLimiter) Try() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano() / int64(time.Millisecond)

	// number of tokens added since last check
	added := l.rate * (now - l.last) / 1000

	// update current number of tokens in bucket
	l.amount = l.amount + added
	if l.amount > l.capacity {
		l.amount = l.capacity
	}
	l.last = now

	// each operation takes 1000 tokens from token bucket
	if l.amount >= 1000 {
		l.amount = l.amount - 1000
		return true
	} else {
		return false
	}
}
