/* b2log_write_test.go - test for b2log_write.go  */
/*
modification history
--------------------
2016/3/8, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package b2log

import (
	"bytes"
    "testing"
)

// test of HeaderWrite()
func Test_HeaderWrite(t *testing.T) {
	// write b2log msg to buffer
	payload := []byte("this is a test")
	buff := make([]byte, HEADER_SIZE + len(payload))
	
	// write header
	err := HeaderWrite(buff, len(payload))
	if err != nil {
		t.Errorf("HeaderWrite():%s", err.Error())
		return
	}
	
	// write payload
	copy(buff[HEADER_SIZE:], payload)

	// try to read b2log msg from buffer
	records, buff := BuffParse(buff)
	if len(records) != 1 {
		t.Errorf("len(records) should be 1, now it's %d", len(records))
		return
	}
	if bytes.Compare(records[0], payload) != 0 {
		t.Errorf("records[0] = %s", records[0])
		return
	}
	if len(buff) != 0 {
		t.Errorf("len(buff) should be 0, now it's %d", len(buff))
	}
}
