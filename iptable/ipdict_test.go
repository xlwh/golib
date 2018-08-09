/* ipdict_test.go - unit test for ipdict.go */
/*
modification history
--------------------
2016/12/21, by Zhang Jiyang, create
*/
/*
 */
package iptable

import (
	"fmt"
	"net"
	"testing"
)

func TestIPDictNormalCase1(t *testing.T) {
	dictSize := 100
	maxItemSize := 100
	d, err := NewIPDict(dictSize, maxItemSize)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 10; i++ {
		startIP := fmt.Sprintf("%d.%d.%d.%d", i, i, i, i+1)
		j := i + 1
		endIP := fmt.Sprintf("%d.%d.%d.%d", j, j, j, j)

		sIP := net.ParseIP(startIP).To4()
		eIP := net.ParseIP(endIP).To4()

		if err = d.Add(sIP, eIP, startIP); err != nil {
			t.Error(err)
			return
		}
	}

	if err = d.SortAndCheck(); err != nil {
		t.Error(err)
		return
	}

	// search exist ip
	ip := net.ParseIP("1.1.1.4").To4()
	val, ok := d.Search(ip)
	if !ok || val != "1.1.1.2" {
		t.Errorf("val should be 1.1.1.2, while %s", val)
	}

	// search no exist ip
	ip = net.ParseIP("0.0.0.0").To4()
	_, ok = d.Search(ip)
	if ok {
		t.Errorf("val should not exist")
	}
}

func TestIPDictWrongCase1(t *testing.T) {
	var sIP, eIP net.IP

	dictSize := 10
	maxItemSize := 10
	d, err := NewIPDict(dictSize, maxItemSize)
	if err != nil {
		t.Error(err)
		return
	}

	// add one more item
	if err = d.Add(sIP, eIP, "largerThanNormalVal"); err == nil {
		t.Error("should add failed")
		return
	}

	for i := 0; i < 10; i++ {
		startIP := fmt.Sprintf("%d.%d.%d.%d", i, i, i, i+1)
		j := i + 1
		endIP := fmt.Sprintf("%d.%d.%d.%d", j, j, j, j)

		sIP = net.ParseIP(startIP).To4()
		eIP = net.ParseIP(endIP).To4()

		if err = d.Add(sIP, eIP, startIP); err != nil {
			t.Error(err)
			return
		}
	}

	// add one more item
	if err = d.Add(sIP, eIP, "test"); err == nil {
		t.Error("should add failed")
		return
	}
}
