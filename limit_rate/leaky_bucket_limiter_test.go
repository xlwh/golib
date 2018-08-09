/* leaky_bucket_limiter_test.go - unit test for leaky_bucket_limiter.go */
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

func TestLeakyBucketLimiter_case1(t *testing.T) {
	l := NewLeakyBucketLimiter(120, 120)
	if !l.Try() {
		t.Error("should not limit")
	}
}

func TestLeakyBucketLimiter_case2(t *testing.T) {
	l := NewLeakyBucketLimiter(120, 120)
	i := 0
	for ; i < 120; i++ {
		if !l.Try() {
			t.Errorf("should not limit (%d)", i)
		}
	}
	for ; i < 150; i++ {
		if !l.Try() {
			t.Errorf("should not limit (%d)", i)
		}
	}
}

func TestLeakyBucketLimiter_case3(t *testing.T) {
	l := NewLeakyBucketLimiter(10, 10)
	i := 0
	for ; i < 10; i++ {
		if _, ok := l.check(); !ok {
			t.Errorf("should not limit (%d)", i)
		}
	}
	for ; i < 20; i++ {
		if _, ok := l.check(); ok {
			t.Errorf("should limit (%d)", i)
		}
	}
}
