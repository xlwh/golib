/* sync_state_test.go - test for sync_state.go */
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
	"strings"
	"testing"
)

import (
	"www.baidu.com/golang-lib/module_state2"
)

func TestNewSyncState(t *testing.T) {
	var state module_state2.State
	_, err := NewSyncState(&state, &SyncStateParams{"totalCountKey", "errCountKey",
		"sucessCountKey", "lastUpdateTimeKey"})
	if err != nil {
		t.Errorf("unexpected err:%v", err)
	}

	_, err = NewSyncState(&state, nil)
	if err == nil {
		t.Errorf("err expected")
		return
	}
	if !strings.Contains(err.Error(), "stateParams should not be nil if state is not nil") {
		t.Errorf("err is not as expected:%v", err)
	}
}

func TestIncTotal(t *testing.T) {
	var state module_state2.State
	state.Init()
	s, _ := NewSyncState(&state, &SyncStateParams{"totalCountKey", "errCountKey",
		"sucessCountKey", "lastUpdateTimeKey"})

	s.IncTotal()
	if s.state.GetCounter("totalCountKey") != 1 {
		t.Errorf("s.state[totalCountKey] should be 1: %d", s.state.GetCounter("totalCountKey"))
	}
}

func TestIncErr(t *testing.T) {
	var state module_state2.State
	state.Init()
	s, _ := NewSyncState(&state, &SyncStateParams{"totalCountKey", "errCountKey",
		"sucessCountKey", "lastUpdateTimeKey"})

	s.IncErr()
	if s.state.GetCounter("errCountKey") != 1 {
		t.Errorf("s.state[errCountKey] should be 1: %d", s.state.GetCounter("errCountKey"))
	}
}

func TestIncSucess(t *testing.T) {
	var state module_state2.State
	state.Init()
	s, _ := NewSyncState(&state, &SyncStateParams{"totalCountKey", "errCountKey",
		"sucessCountKey", "lastUpdateTimeKey"})

	s.IncSucess()
	if s.state.GetCounter("sucessCountKey") != 1 {
		t.Errorf("s.state[sucessCountKey] should be 1: %d", s.state.GetCounter("sucessCountKey"))
	}
}

func TestSetLastUpdate(t *testing.T) {
	var state module_state2.State
	state.Init()
	s, _ := NewSyncState(&state, &SyncStateParams{"totalCountKey", "errCountKey",
		"sucessCountKey", "lastUpdateTimeKey"})

	s.SetLastUpdate("2015/03/04 19:33:33")
	if s.state.GetState("lastUpdateTimeKey") != "2015/03/04 19:33:33" {
		t.Errorf("s.state[lastUpdateTimeKey] should be '2015/03/04 19:33:33': %s",
			s.state.GetState("lastUpdateTimeKey"))
	}
}
