/* string_slice_test.go - test for string_slice.go    */
/*
modification history
--------------------
2015/8/13, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package string_slice

import (
	"testing"
)

func Test_InSlice_case1(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !InSlice("a", slice) {
		t.Errorf("'a' should in slice")
	}

	if InSlice("d", slice) {
		t.Errorf("'d' should not in slice")
	}
}

func Test_EqualIgnoreOrder_case1(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"a", "c", "b"}

	if !EqualIgnoreOrder(s1, s2) {
		t.Errorf("s1 and s2 should be equal without order")
	}

	s3 := []string{"a", "b", "b"}
	if EqualIgnoreOrder(s1, s3) {
		t.Errorf("s1 and s3 should be unequal without order")
	}

	s4 := []string{"a", "b"}
	if EqualIgnoreOrder(s1, s4) {
		t.Errorf("s1 and s4 should be unequal without order")
	}
}
