/* byte_pool_test.go - unit test file for byte_pool.go*/ /*
modification history
--------------------
2014/8/26, by zhangjiyang01@baidu.com, create
*/

package byte_pool

import (
	"bytes"
	"testing"
)

func TestBytePool(t *testing.T) {
	eleNum := 2
	maxElemSize := 13

	pool := NewBytePool(eleNum, maxElemSize)
	if pool.MaxElemSize() != maxElemSize {
		t.Error("t.elemeSize error")
	}

	key := []byte("hello world")
	if err := pool.Set(1, key); err != nil {
		t.Error("set should be success")
	}

	result := pool.Get(1)
	if len(key) != len(result) {
		t.Error("result should keep length")
	}

	if bytes.Compare(key, result) != 0 {
		t.Error("item should keep unchanged")
	}

	if err := pool.Set(2, key); err == nil {
		t.Error("set should failed")
	}

	key = []byte("large than max ele size")
	if err := pool.Set(1, key); err == nil {
		t.Error("set should failed")
	}
}
