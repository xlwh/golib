/* b2log_write.go - write b2log record from file  */
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
	"bytes"
    "encoding/binary"
    "fmt"
	"time"
)

// generate timestamp
func timestampGen() uint64 {
	t := time.Now()
	sec := t.Unix()
	usec := t.Nanosecond() / 1000
	ts := uint64(sec * 1000 + int64(usec) / 1000)

	return ts
}

/* 
HeaderWrite - write b2log header to given buffer

Params:
	- buffer: []byte to write header to
	- payloadLen: length of payload
*/
func HeaderWrite(buffer []byte, payloadLen int) error {
	// prepare header
	header := Header{MAGIC_NUMBER, HEADER_VERSION, uint32(payloadLen), 0, 0}
	header.TimeStamp = timestampGen()

	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, header)
	if err != nil {
		return fmt.Errorf("binary.Write():%s", err.Error())
	}
	
	// write header to buffer
	copy(buffer, buff.Bytes())
	
	return nil
}
