/* ip_loc_table.go - test file of ipdict */
/*
modification history
--------------------
2016/7/13, by Jiang Hui, create
*/
/*
DESCRIPTION
*/

package ipdict

import (
	"strings"
	"testing"
)

import (
	"www.baidu.com/golang-lib/net_util"
)

// new case
func TestNewIpLocationTable(t *testing.T) {
	_, err := NewIpLocationTable(0, 1)
	if err == nil {
		t.Errorf("NewIpLocationTable should return err!=nil but return nil")
	}

	_, err = NewIpLocationTable(1, 0)
	if err == nil {
		t.Errorf("NewIpLocationTable should return err!=nil but return nil")
	}

	_, err = NewIpLocationTable(1000001, 10)
	if err == nil {
		t.Errorf("NewIpLocationTable should return err!=nil but not")
	}

	_, err = NewIpLocationTable(1000, 1025)
	if err == nil {
		t.Errorf("NewIpLocationTable should return err!=nil but not")
	}
}

//  Add case
func TestLocAdd(t *testing.T) {
	locTable, err := NewIpLocationTable(1, 48)
	if err != nil {
		t.Errorf("NewIpLocationTable should return err==nil but return not")
	}

	ipStrS := "223.255.192.0"
	ipStrE := "223.255.223.255"
	ipLoc := "KR:None:None"
	//223.255.192.0|223.255.223.255|KR|None|None|None|None|80|0|0|0|0
	ipS := net_util.ParseIPv4(ipStrS)
	ipE := net_util.ParseIPv4(ipStrE)
	err = locTable.Add(ipS, ipE, ipLoc)
	if err != nil {
		t.Errorf("locTable Add should return err==nil but return not")
	}

	//223.255.128.0|223.255.191.255|HK|None|XIANGGANG|XIANGGANG|None|80|0|80|80|0
	ipStrS = "223.255.128.0"
	ipStrE = "223.255.191.255"
	ipLoc = "HK:XIANGGANG:XIANGGANG"
	ipS = net_util.ParseIPv4(ipStrS)
	ipE = net_util.ParseIPv4(ipStrE)
	err = locTable.Add(ipS, ipE, ipLoc)
	if err == nil {
		t.Errorf("locTable Add should return err!=nil but return nil")
	}
}

//  Search case
func TestLocSearch(t *testing.T) {
	locTable, err := NewIpLocationTable(3, 48)
	if err != nil {
		t.Errorf("NewIpLocationTable should return err==nil but not")
	}

	//223.255.192.0|223.255.223.255|KR|None|None|None|None|80|0|0|0|0
	ipStrS := "223.255.192.0"
	ipStrE := "223.255.223.255"
	ipLoc := "KR:None:None"
	ipS := net_util.ParseIPv4(ipStrS)
	ipE := net_util.ParseIPv4(ipStrE)
	err = locTable.Add(ipS, ipE, ipLoc)
	if err != nil {
		t.Errorf("locTable Add should return err==nil but not")
	}

	ip := net_util.ParseIPv4(ipStrS)
	var loc string
	loc, err = locTable.Search(ip)
	if err != nil {
		t.Errorf("locTable Search should return err==nil but not")
	}

	if !strings.EqualFold("KR:None:None", loc) {
		t.Errorf("locTable Search should return KR:None:None but is %s", loc)
	}

	//223.255.128.0|223.255.191.255|HK|None|XIANGGANG|XIANGGANG|None|80|0|80|80|0
	ipStrS = "223.255.128.0"
	ipStrE = "223.255.191.255"
	ipLoc = "HK:XIANGGANG:XIANGGANG"
	ipS = net_util.ParseIPv4(ipStrS)
	ipE = net_util.ParseIPv4(ipStrE)
	locTable.Add(ipS, ipE, ipLoc)
	if err != nil {
		t.Errorf("locTable Add should return err==nil but not")
	}

	ip = net_util.ParseIPv4(ipStrE)
	loc, err = locTable.Search(ip)
	if err != nil {
		t.Errorf("locTable Search should return err==nil but not")
	}
	if !strings.EqualFold("HK:XIANGGANG:XIANGGANG", loc) {
		t.Errorf("locTable Search should return HK:XIANGGANG:XIANGGANG but is %s", loc)
	}
}
