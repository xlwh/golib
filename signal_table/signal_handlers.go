/* signal_handlers.go - common signal handlers  */
/*
modification history
--------------------
2014/7/15, by Li Bingyi, create
*/
/*
DESCRIPTION
    This file provides two common handlers to deal with signals.
    TermHandler handler deal with the signal that you want to terminate the process
    IgnoreHandler handler deal with the signal that you want to ignore.
*/

package signal_table

import (
    "os"
)

import (
    "www.baidu.com/golang-lib/log"
)

/* TermHandler deal with the signal that should terminate the process */
func TermHandler(s os.Signal) {
    log.Logger.Info("termHandler(): receive signal[%v], terminate.", s)
    log.Logger.Close()
    os.Exit(0)
}

/* IgnoreHandler deal with the signal that should be ignored */
func IgnoreHandler(s os.Signal) {
    log.Logger.Info("ignoreHandler(): receive signal[%v], ignore.", s)
}
