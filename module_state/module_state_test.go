/* module_state_test.go - test for module_state.go  */
/*
modification history
--------------------
2014/3/25, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package module_state

import (
    "fmt"
    "testing"
)

func TestModuleState(t *testing.T) {
    var state State

    state.Init()        
    state.Inc("counter", 1)
    
    // test GetAll()
    table := state.GetAll()
    fmt.Println(table)
    
    // test Get()
    value := state.Get("counter")
    if value != 1 {
        t.Error("value should be 1")
    }
    
    // test Get() for noexist key
    value = state.Get("counter1")
    if value != 0 {
        t.Error("value should be 0")
    }
    
}