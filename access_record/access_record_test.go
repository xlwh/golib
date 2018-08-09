/* access_record_test.go - test for access_record.go */
/*
modification history
--------------------
2015/9/2, by Guang Yao, create
*/
package access_record

import (
	"testing"
	"time"
)

// normal case
func Test_AccessRecord_case1(t *testing.T) {
	record := NewAccessRecord()

	record.Inc("127.0.0.1", time.Now())

	// get
	r := record.Get("127.0.0.1")
	if r == nil {
		t.Errorf("err in record.Get: nil returned")
		return
	}

	// inc again
	record.Inc("127.0.0.1", time.Now())
	r = record.Get("127.0.0.1")
	if r == nil {
		t.Errorf("err in record.Get: nil returned")
		return
	}
	if r.Count != 2 {
		t.Errorf("err in record.Get: r.Count[%d]!=2", r.Count)
		return
	}

	// inc another and get all
	record.Inc("127.0.0.2", time.Now())
	records := record.GetAll()
	if len(records) != 2 {
		t.Errorf("err in record.GetAll: wrong record number: %d", len(records))
		return
	}
}
