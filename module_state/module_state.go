/* module_state.go - for collecting state info of a module  */
/*
modification history
--------------------
2014/3/25, by Zhang Miao, create
*/
/*
DESCRIPTION

Usage:
    import "www.baidu.com/golang-lib/module_state"
    
    var state module_state.State
    
    state.Init()
    
    state.Inc("counter", 1)
    
    table := state.Get()
*/
package module_state

import (
    "sync"    
)

/* state    */
type StateTable map[string]int64

type State struct {
    lock    sync.Mutex
    table   StateTable
}

/* Initialize the state */
func (s *State) Init() {
    s.table = make(StateTable)
}

/* Add value to key */
func (s *State) Inc(key string, value int) {
    s.lock.Lock()
    defer s.lock.Unlock()
    
    _, ok := s.table[key]
    
    if !ok {
        s.table[key] = int64(value)
    } else {
        s.table[key] += int64(value)
    }    
}

/* Get all states    */
func (s *State) GetAll() StateTable {
    tCopy := make(StateTable)
    
    s.lock.Lock()
    defer s.lock.Unlock()
    
    for key, value := range s.table {
        tCopy[key] = value
    }
    
    return tCopy
}

/* Get value of given key    */
func (s *State) Get(key string) int64 {    
    s.lock.Lock()    
    value, ok := s.table[key]    
    s.lock.Unlock()
    
    if !ok {
        value = 0
    }
    
    return value
}