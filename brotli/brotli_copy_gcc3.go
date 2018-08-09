// +build !gcc4

/* brotli_copy_gcc3.go - memcopy for gcc3 */
/*
modification history
--------------------
2017/04/11, by Sijie Yang, create
*/
/*
DESCRIPTION
  Note: build constraint must be preceded only by blank lines and
other *line* comments comments, and be followed by a blank line
*/

package brotli

/*
#include <string.h>
*/
import "C"

import (
	"unsafe"
)

func MemCopy(dst unsafe.Pointer, src unsafe.Pointer, size C.size_t) {
	C.memcpy(dst, src, C.ulong(size))
}
