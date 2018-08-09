/* sync_state.go - for record state in sync */
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
	"fmt"
)

import (
	"www.baidu.com/golang-lib/module_state2"
)

// for record state in sync
type SyncState struct {
	state       *module_state2.State // the link to state
	stateParams *SyncStateParams     // params in state
}

func NewSyncState(state *module_state2.State, stateParams *SyncStateParams) (*SyncState, error) {
	// stateParams should not be nil if state is not nil
	if state != nil && stateParams == nil {
		return nil, fmt.Errorf("stateParams should not be nil if state is not nil")
	}

	s := new(SyncState)
	s.state = state
	s.stateParams = stateParams

	return s, nil
}

func (s *SyncState) IncTotal() {
	if s.state != nil {
		s.state.Inc(s.stateParams.TotalCountKey, 1)
	}
}

func (s *SyncState) IncErr() {
	if s.state != nil {
		s.state.Inc(s.stateParams.ErrCountKey, 1)
	}
}

func (s *SyncState) IncSucess() {
	if s.state != nil {
		s.state.Inc(s.stateParams.SucessCountKey, 1)
	}
}

func (s *SyncState) SetLastUpdate(lastUpdateTime string) {
	if s.state != nil {
		s.state.Set(s.stateParams.LastUpdateTimeKey, lastUpdateTime)
	}
}
