/* stack.go - utilities for runtime stack */
/*
modification history
--------------------
2017/06/13, by Sijie Yang, create
*/
/*
DESCRIPTION
*/
package gotrack

import (
	"runtime"
)

const StackSize = 4096

func CurrentStackTrace(size int) []byte {
	// use default size for stack trace
	if size <= 0 {
		size = StackSize
	}

	// get stack trace
	buf := make([]byte, size)
	n := runtime.Stack(buf, false)
	return buf[:n]
}
