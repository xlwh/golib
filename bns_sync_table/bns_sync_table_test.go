/* bns_sync_table_test.go - test for bns_sync_table.go */
/*
modification history
--------------------
2016/2/6, by Guang Yao, create
*/
/*
DESCRIPTION
*/
package bns_sync_table

import (
	"reflect"
	"strings"
	"testing"
)

import (
	"www.baidu.com/golang-lib/module_state2"
)

// test for parseBnsTags
func TestParseBnsTags(t *testing.T) {
	// normal case
	tagMap, err := parseBnsTags("interface:clientip,keepalive:0,weight:10")
	if err != nil {
		t.Errorf("unexpected error:%v", err)
		return
	}
	correctTagMap := map[string]string{
		"interface": "clientip",
		"keepalive": "0",
		"weight":    "10",
	}
	if !reflect.DeepEqual(tagMap, correctTagMap) {
		t.Errorf("tagMap is not as expected: %v", tagMap)
	}

	// bad case 1
	_, err = parseBnsTags("interface,clientip,keepalive,0,weight,10")
	if err == nil {
		t.Errorf("err is expected")
	}
	if !strings.Contains(err.Error(), "err in parse tagItem") {
		t.Errorf("err is not as expected: %v", err)
	}

	// bad case 2
	_, err = parseBnsTags(" ")
	if err == nil {
		t.Errorf("err is expected")
	}
	if !strings.Contains(err.Error(), "err in parse tagItem") {
		t.Errorf("err is not as expected: %v", err)
	}
}

// test parseBnsAllInfoLine
func TestParseBnsAllInfoLine(t *testing.T) {
	// normal case
	line := "hkg02-l-bfe03.hkg02 10.242.123.45 bfe-hk.BFE.hk 8900 0 -1 2 interface:clientip,keepalive:0,weight:10 extra container /tmp"
	instance, err := parseBnsAllInfoLine(line)
	if err != nil {
		t.Errorf("unexpected err:%v", err)
		return
	}
	correctInstance := &BnsInstance{
		"hkg02-l-bfe03.hkg02",
		"10.242.123.45",
		"bfe-hk.BFE.hk",
		8900,
		0,
		-1,
		2,
		map[string]string{
			"interface": "clientip",
			"keepalive": "0",
			"weight":    "10",
		},
		"extra",
		"container",
		"/tmp",
	}
	if !reflect.DeepEqual(instance, correctInstance) {
		t.Errorf("instance is not as expected:%v, %v", instance, correctInstance)
	}

	// good case: not tag
	line = "hkg02-l-bfe03.hkg02 10.242.123.45 bfe-hk.BFE.hk 8900 0 -1 2    /tmp"
	instance, err = parseBnsAllInfoLine(line)
	if err != nil {
		t.Errorf("unexpected err:%v", err)
		return
	}

	// bad case 1: items count != 11
	line = "hkg02-l-bfe03.hkg02 10.242.123.45 bfe-hk.BFE.hk 8900 0 -1 2"
	_, err = parseBnsAllInfoLine(line)
	if err == nil {
		t.Errorf("err is expected")
		return
	}
	if !strings.Contains(err.Error(), "unexpected bns output format") {
		t.Errorf("err is not as expected:%v", err)
	}

	// bad case 2: bad port
	line = "hkg02-l-bfe03.hkg02 10.242.123.45 bfe-hk.BFE.hk err 0 -1 2 interface:clientip,keepalive:0,weight:10 extra container /tmp"
	_, err = parseBnsAllInfoLine(line)
	if err == nil {
		t.Errorf("err is expected")
		return
	}
	if !strings.Contains(err.Error(), "err in parse port") {
		t.Errorf("err is not as expected:%v", err)
	}

	// bad case 3: bad status
	line = "hkg02-l-bfe03.hkg02 10.242.123.45 bfe-hk.BFE.hk 8900 err -1 2 interface:clientip,keepalive:0,weight:10 extra container /tmp"
	_, err = parseBnsAllInfoLine(line)
	if err == nil {
		t.Errorf("err is expected")
		return
	}
	if !strings.Contains(err.Error(), "err in parse status") {
		t.Errorf("err is not as expected:%v", err)
	}

	// bad case 4: bad load
	line = "hkg02-l-bfe03.hkg02 10.242.123.45 bfe-hk.BFE.hk 8900 0 err 2 interface:clientip,keepalive:0,weight:10 extra container /tmp"
	_, err = parseBnsAllInfoLine(line)
	if err == nil {
		t.Errorf("err is expected")
		return
	}
	if !strings.Contains(err.Error(), "err in parse load") {
		t.Errorf("err is not as expected:%v", err)
	}

	// bad case 5: bad offset
	line = "hkg02-l-bfe03.hkg02 10.242.123.45 bfe-hk.BFE.hk 8900 0 -1 err interface:clientip,keepalive:0,weight:10 extra container /tmp"
	_, err = parseBnsAllInfoLine(line)
	if err == nil {
		t.Errorf("err is expected")
		return
	}
	if !strings.Contains(err.Error(), "err in parse offset") {
		t.Errorf("err is not as expected:%v", err)
	}
}

// test resolveBns
func TestResolveBns(t *testing.T) {
	// normal case
	_, err := resolveBns("gslb-scheduler.bfe.all")
	if err != nil {
		t.Errorf("unexpected err:%v", err)
		return
	}
}

// test update, has and get
func TestTableOps(t *testing.T) {
	table := new(BnsSyncTable)
	instances := []*BnsInstance{
		&BnsInstance{
			"hkg02-l-bfe03.hkg02",
			"10.242.123.45",
			"bfe-hk.BFE.hk",
			8900,
			0,
			-1,
			2,
			map[string]string{
				"interface": "clientip",
				"keepalive": "0",
				"weight":    "10",
			},
			"extra",
			"container",
			"/tmp",
		},
		&BnsInstance{
			"hkg02-l-bfe02.hkg02",
			"10.242.123.43",
			"bfe-hk.BFE.hk",
			8900,
			0,
			-1,
			2,
			map[string]string{
				"interface": "clientip",
				"keepalive": "0",
				"weight":    "10",
			},
			"extra",
			"container",
			"/tmp",
		},
	}

	// update
	table.update(instances)

	// has host and has ip
	if !table.HasHost("hkg02-l-bfe03.hkg02") || !table.HasHost("hkg02-l-bfe02.hkg02") {
		t.Errorf("table is not as expected: %+v", table)
		return
	}
	if table.HasHost("hkg02-l-bfe01.hkg02") || table.HasHost("hkg02-l-bfe00.hkg02") {
		t.Errorf("table is not as expected: %+v", table)
		return
	}
	if !table.HasIP("10.242.123.43") || !table.HasIP("10.242.123.45") {
		t.Errorf("table is not as expected: %+v", table)
		return
	}
	if table.HasIP("10.242.123.41") || table.HasIP("10.242.123.42") {
		t.Errorf("table is not as expected: %+v", table)
		return
	}

	// get all
	tableInstances := table.GetAll()
	alterInstances := []*BnsInstance{instances[1], instances[0]}
	if !reflect.DeepEqual(tableInstances, instances) &&
		!reflect.DeepEqual(tableInstances, alterInstances) {
		t.Errorf("table is not as expected: %+v", table)
		return
	}
}

// test NewBnsSyncTable
func TestNewBnsSyncTable(t *testing.T) {
	var state module_state2.State
	state.Init()

	_, err := NewBnsSyncTable("gslb-scheduler.bfe.all", 5, "./test_data/", &state,
		&SyncStateParams{"totalCountKey", "errCountKey", "sucessCountKey", "lastUpdateTimeKey"})
	if err != nil {
		t.Errorf("unexpected err:%v", err)
		return
	}

}
