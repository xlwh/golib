/* leaky_bucket_limiter.go - a rate limiter using leaky bucket algorithm */
/*
modification history
--------------------
2015/5/20, by Sijie Yang, Create
*/
/*
DESCRIPTION
    move from go-bfe/src/bfe_util
*/
package limit_rate

import (
	"sync"
	"time"
)

type LeakyBucketLimiter struct {
	rate     int64 // units dripped per second
	capacity int64 // capacity of bucket
	amount   int64 // current amount of units in bucket
	last     int64 // timestamp of last check
	lock     sync.Mutex
}

/* NewLeakyBucketLimiter - create a rate limiter
*
* Params:
*     - ops  : maximum operation per second
*     - burst: maximum burst number of operation. If the operations rate exceed 'ops',
*              their processing is delayed such that operations are processed at a defined rate;
*              After the number of excessive operations exceeds 'burst', new incoming
               operations will be dropped.
* Return:
*     - rate limter
*/
func NewLeakyBucketLimiter(ops int64, burst int64) *LeakyBucketLimiter {
	l := new(LeakyBucketLimiter)
	if ops <= 0 {
		ops = 1000 // default maximum operation per second
	}
	if burst <= 0 {
		burst = 1000 // default maximum burst number of operation
	}

	// Note: each operation will add 1000 units to bucket
	l.rate = ops * 1000
	l.capacity = burst * 1000
	l.last = time.Now().UnixNano() / int64(time.Millisecond)
	return l
}

/* Try - check whether operation is allowable or should be dropped
 *
 * Return:
 *     - ret: true if allowable, false if not
 */
func (l *LeakyBucketLimiter) Try() bool {
	delay, allowable := l.check()
	// should drop operation
	if !allowable {
		return false
	}

	// delay if nessary such that operations are processed at a defined rate
	time.Sleep(delay)
	return true
}

// check bucket and return status for new operation
func (l *LeakyBucketLimiter) check() (time.Duration, bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano() / int64(time.Millisecond)

	// number of units dripped since last check
	leak := l.rate * (now - l.last) / 1000

	// update current number of units in bucket
	l.amount = l.amount - leak
	if l.amount < 0 {
		l.amount = 0
	}
	l.last = now

	// time to wait for incoming operation
	delay := time.Duration(l.amount*1000/l.rate) * time.Millisecond

	// each operation add 1000 units to leaky bucket
	l.amount = l.amount + 1000
	if l.amount > l.capacity {
		l.amount = l.capacity
		return 0, false
	}

	return delay, true
}
