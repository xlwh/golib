/* signal_table_test.go - test of signal_table.go
/*
modification history
--------------------
2014/7/15, by Li Bingyi, create
*/
/*
DESCRIPTION
*/

package signal_table

import (
    "syscall"
    "testing"
)

func TestRegister(t *testing.T) {
    table := NewSignalTable()
    var h signalHandler

    table.Register(syscall.SIGHUP, nil)
    table.Register(syscall.SIGHUP, h)

    if table.shs[syscall.SIGHUP] != nil {
        t.Error("Register(): err in Register")
    }
}

func TestHandle(t *testing.T) {
    table := NewSignalTable()

    table.Register(syscall.SIGHUP, IgnoreHandler)

    table.handle(syscall.SIGHUP)

    if table.state.GetCounter("hangup") != 1 {
        t.Error("handler(): counter != 1")
    }
}
