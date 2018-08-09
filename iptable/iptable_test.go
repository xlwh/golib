/* iptable_test.go - unit test file for iptable.go */
/*
modification history
--------------------
2016/12/21, by Zhang Jiyang, create
*/
/*
 */
package iptable

import (
	"net"
	"testing"
)

func TestIPTableCase1(t *testing.T) {
	table := NewIPTable()

	capacity := 10
	maxItemSize := 10
	ipdict, err := NewIPDict(capacity, maxItemSize)
	if err != nil {
		t.Errorf("NewIPDict should success")
	}

	version := "testVersion"
	table.Update(ipdict, version)

	// update
	table.Update(ipdict, version)

	// check
	if table.Version() != version {
		t.Error("version should be same")
	}

	// search
	ip := net.ParseIP("192.168.0.1")
	_, ok := table.Search(ip)
	if ok {
		t.Error("should search failed")
	}
}
