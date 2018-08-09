/* register_signal.go - default register functions for signal handlers */
/*
modification history
--------------------
2014/08/13, by wangpengfei05, create
2016/1/12, by Zhang Miao, copy from go-bfereader
*/
/*
DESCRIPTION
*/
package signal_table

import (
	"syscall"
)

// register signal handlers
func RegisterSignalHandlers(signalTable *SignalTable) {
	// term handlers
	signalTable.Register(syscall.SIGTERM, TermHandler)

	// ignore handlers
	signalTable.Register(syscall.SIGHUP, IgnoreHandler)
	signalTable.Register(syscall.SIGQUIT, IgnoreHandler)
	signalTable.Register(syscall.SIGILL, IgnoreHandler)
	signalTable.Register(syscall.SIGTRAP, IgnoreHandler)
	signalTable.Register(syscall.SIGABRT, IgnoreHandler)
}
