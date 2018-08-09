/* signal_table.go - signal table  */
/*
modification history
--------------------
2018/5/8, by Lin Ming, create
*/
/*
DESCRIPTION

Usage:

    import (
        "syscall"
    )

    import (
        "www.baidu.com/golang-lib/signal_table"
    )

    var signalTable signal_table.SignalTable
    var s os.Signal
    var handler signalHandler

    // new signal table
    signalTable = signal_table.NewSignalTable()

    // register signal handlers
    signalTable.Register(s, handler)

    // start signal handle go-routin
    signalTable.StartSignalHandle()
*/

package signal_table

import (
"encoding/json"
"os"
"os/signal"
)

import (
"www.baidu.com/golang-lib/module_state2"
)

type signalHandler func(s os.Signal)

type SignalTable struct {
	shs   map[os.Signal]signalHandler // signal handle table
	state module_state2.State         // signal handle state
}

/* new and init signal table */
func NewSignalTable() (*SignalTable){
	table := new(SignalTable)
	table.shs = make(map[os.Signal]signalHandler)
	table.state.Init()
	return table
}

/* register signal handle to the table */
func (t *SignalTable) Register(s os.Signal, handler signalHandler) {
	if _, ok := t.shs[s]; !ok {
		t.shs[s] = handler
	}
}

/* handle for the related signal */
func (t *SignalTable) handle(sig os.Signal) {
	t.state.Inc(sig.String(), 1)

	if handler, ok := t.shs[sig]; ok {
		handler(sig)
	}
}

// signal handle go-routine
func (table *SignalTable)signalHandle() {

	var sigs []os.Signal
	for sig := range table.shs {
		sigs = append(sigs, sig)
	}

	c := make(chan os.Signal, len(sigs))
	signal.Notify(c, sigs...)

	for {
		sig := <-c
		table.handle(sig)
	}
}

/*  start go-routine for signal handle */
func (t *SignalTable)StartSignalHandle() {
	go t.signalHandle()
}

/* get state counter of signal handle */
func (t *SignalTable) SignalStateGet() ([]byte, error) {

	buff, err := json.Marshal(t.state.GetAll())

	return buff, err
}

/* set noah prefix key */
func (t *SignalTable) SetNoahKeyPrefix(key string) {
	t.state.SetNoahKeyPrefix(key)
}

/* get noah prefix key */
func (t *SignalTable) GetNoahKeyPrefix() string {
	return t.state.GetNoahKeyPrefix()
}
