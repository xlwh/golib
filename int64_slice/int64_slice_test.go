/* int64_slice_test.go - test for int64_slice.go  */
/*
modification history
--------------------
2016/2/2, by Guang Yao, create
*/

package int64_slice

import (
	"testing"
)

func TestSort(t *testing.T) {
	s := Int64Slice{1454400523, 1454400521, 1454400522}
	s_sorted := Int64Slice{1454400521, 1454400522, 1454400523}
	s.Sort()

	for i, _ := range s {
		if s[i] != s_sorted[i] {
			t.Errorf("Invalid Sort result: [%v]=>[%v]",
				Int64Slice{1454400523, 1454400521, 1454400522}, s)
		}
	}
}
