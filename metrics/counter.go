/* counter.go - counter */
/*
modification history
--------------------
2016/12/19, by Sijie Yang, create
*/
/*
DESCRIPTION
*/
package metrics

import (
	"sync/atomic"
)

type Counter int64

// increase counter
func (c *Counter) Inc(delta int) {
	if c == nil {
		return
	}
	atomic.AddInt64((*int64)(c), int64(delta))
}

// get counter
func (c *Counter) Get() int64 {
	if c == nil {
		return 0
	}
	return atomic.LoadInt64((*int64)(c))
}
