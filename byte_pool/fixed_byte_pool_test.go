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

func TestFixedBytePool(t *testing.T) {
	key := []byte("hello world")
	keyLen := len(key)

	pool := NewFixedBytePool(2, keyLen)

	if pool.MaxElemSize() != keyLen {
		t.Error("t.elemeSize error")
	}

	if err := pool.Set(1, key); err != nil {
		t.Error("set should be success")
	}

	resuItem := pool.Get(1)

	if len(key) != len(resuItem) {
		t.Error("testItem and resuItem not same len")
	}
	if bytes.Compare(key, resuItem) != 0 {
		t.Error("testItem, and resuItem not equal")
	}

	if err := pool.Set(2, key); err == nil {
		t.Error("set should failed")
	}

	key = []byte("large than max ele size")
	if err := pool.Set(1, key); err == nil {
		t.Error("set should failed")
	}

}
