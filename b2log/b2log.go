/* b2log.go - read and write b2log record   */
/*
modification history
--------------------
2014/11/4, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package b2log

import (
    "unsafe"
)

// magic number
// This can only work in little-endian machine (e.g., x86)
// Remember to make MAGIC_NUMBER and MAGIC_NUMBER_STR consistent
const (
	MAGIC_NUMBER	= 0xB0AEBEA7
	HEADER_VERSION 	= 1
)

var   MAGIC_NUMBER_STR  = []byte{0xA7, 0xBE, 0xAE, 0xB0}

// size of Header
var HEADER_SIZE = int(unsafe.Sizeof(demoHeader))
var demoHeader Header       // this var is only for getting size of Header

const MAX_RECORD_LEN    = 100 * 1024    // max length of single b2log record

// header for b2log record
type Header struct {
    MagicNumber   uint32    // magic number
    Version       uint32    // version 
    UnCompressLen uint32    // length of upcompress log
    CompressLen   uint32    // length of compress log
    TimeStamp     uint64    // timestamp the log generated
}

// binary format of b2log record
type Record []byte

