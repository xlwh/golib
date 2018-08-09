/* token_bucket_limiter_test.go - unit test for token_bucket_limiter.go */
/*
modification history
--------------------
2015/5/20, by Sijie Yang, Create
*/
/*
DESCRIPTION
*/
package limit_rate

import (
	"testing"
)

func TestTokenBucketLimiter_case1(t *testing.T) {
	l := NewTokenBucketLimiter(10, 10)
	l.amount = l.capacity

	if !l.Try() {
		t.Error("should not limit")
	}
}

func TestTokenBucketLimiter_case2(t *testing.T) {
	l := NewTokenBucketLimiter(10, 10)
	l.amount = l.capacity

	i := 0
	for ; i < 10; i++ {
		if !l.Try() {
			t.Errorf("should not limit (%d)", i)
		}
	}

	l.amount = 0
	if l.Try() {
		t.Errorf("should limit (%d)", i)
	}
}
