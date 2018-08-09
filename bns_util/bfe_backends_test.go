/* bfe_backends_test.go - unit test for bfe_backends */
/*
modification history
--------------------
2017/11/30, by yuxiaofei, create
*/
/*
DESCRIPTION
*/

package bns_util

import (
	"testing"
)

// test parseBfeBackends, without default weight
func TestParseBfeBackends_1(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 0 -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0 weight:10  `)
	withDefaultWeight := false

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 2 {
		t.Errorf("length of bks should be 2, not %d", len(bks))
	}
}

// test parseBfeBackends, with default weight
func TestParseBfeBackends_2(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 0 -1 1
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0`)
	withDefaultWeight := true

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 2 {
		t.Errorf("length of bks should be 2, not %d", len(bks))
	}

	for _, bk := range bks {
		if *bk.Weight != 1 {
			t.Errorf("weight should be 1, not %d", *bk.Weight)
		}
	}
}

// test parseBfeBackends, with default weight, has status not 0
func TestParseBfeBackends_3(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 999 -1 1 
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0`)
	withDefaultWeight := true

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 1 {
		t.Errorf("length of bks should be 1, not %d", len(bks))
	}

	for _, bk := range bks {
		if *bk.Weight != 1 {
			t.Errorf("weight should be 1, not %d", *bk.Weight)
		}
	}
}

// test parseBfeBackends, without default weight, has status not 0
func TestParseBfeBackends_4(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 999 -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0 weight:10  `)
	withDefaultWeight := false

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 1 {
		t.Errorf("length of bks should be 1, not %d", len(bks))
	}
}

// test parseBfeBackends, without default weight, has field lenght < 8
func TestParseBfeBackends_5(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 0 -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0`)
	withDefaultWeight := false

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 1 {
		t.Errorf("length of bks should be 1, not %d", len(bks))
	}
}

// test parseBfeBackends, with default weight, has field lenght < 7
func TestParseBfeBackends_6(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 0 -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1`)
	withDefaultWeight := true

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 1 {
		t.Errorf("length of bks should be 1, not %d", len(bks))
	}
}

// test parseBfeBackends, without default weight, port field not int
func TestParseBfeBackends_7(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf abc 0 -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0 weight:10   `)
	withDefaultWeight := false

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 1 {
		t.Errorf("length of bks should be 1, not %d", len(bks))
	}
}

// test parseBfeBackends, with default weight, port field not int
func TestParseBfeBackends_8(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf abc 0 -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0 weight:10   `)
	withDefaultWeight := true

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 1 {
		t.Errorf("length of bks should be 1, not %d", len(bks))
	}
}

// test parseBfeBackends, without default weight, status field not int
func TestParseBfeBackends_9(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf abc 0 -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0 weight:10   `)
	withDefaultWeight := false

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 1 {
		t.Errorf("length of bks should be 1, not %d", len(bks))
	}
}

// test parseBfeBackends, with default weight, status field not int
func TestParseBfeBackends_10(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 a -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 0 -1 0 weight:10   `)
	withDefaultWeight := true

	if err := parseBfeBackends(out, &bks, withDefaultWeight); err != nil {
		t.Errorf("parseBfeBackends(): %v", err)
	}

	if len(bks) != 1 {
		t.Errorf("length of bks should be 1, not %d", len(bks))
	}
}

// test parseBfeBackends, without default weight, 0 avail instance
func TestParseBfeBackends_11(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 999 -1 1 weight:10   
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 999 -1 0 weight:10   `)
	withDefaultWeight := false

	err := parseBfeBackends(out, &bks, withDefaultWeight)
	if err == nil && err.Error() != "get 0 avail instance from output" {
		t.Errorf("parseBfeBackends() wrong error: %v", err)
	}
}

// test parseBfeBackends, with default weight, 0 avail instance
func TestParseBfeBackends_12(t *testing.T) {
	var bks BfeBackendConfList
	out := []byte(`yf-s-bfe01.yf01 10.36.22.24 small.BFE.yf 8900 999 -1 1
yf-s-bfe00.yf01 10.38.159.43 small.BFE.yf 8900 999 -1 0`)
	withDefaultWeight := true

	err := parseBfeBackends(out, &bks, withDefaultWeight)
	if err == nil && err.Error() != "get 0 avail instance from output" {
		t.Errorf("parseBfeBackends() wrong error: %v", err)
	}
}
