/* tcp_conn_table.go - tcp connection table */
/*
modification history
--------------------
2014/3/10, by Zhang Miao, create
2014/8/6, by Zhang Miao, move from waf_server
*/
/*
DESCRIPTION
*/
package net_server

import (
    "errors"
    "net"
    "sync"
)

import (
    "www.baidu.com/golang-lib/queue"
)

type TcpConn struct {
    SendQueue   queue.Queue // for send to client
    LocalAddr   string      // local address
    RemoteAddr  string      // remote address
}

/* Initialize new TcpConn */
func newTcpConn() TcpConn {
    tcpConn := TcpConn{}
    
    tcpConn.SendQueue.Init()
        
    return tcpConn
}

// send out data via TcpConn
func (conn *TcpConn) Send(data interface{}) {
    // prepare SendMsg
    msg := &SendMsg{
                cmd: SEND_CMD_SEND,
                data:data,
            }

    // add to sending queue
    conn.SendQueue.Append(msg)
}

/* tcp connection table */
type TcpConnTable struct {
    lock    sync.Mutex
    table   map[net.Conn]TcpConn
}

/* state of tcp connection  */
type TcpConnState struct {
    QueueLen    int     // length of send queue
    LocalAddr   string  // local address
    RemoteAddr  string  // remote address
}

/* Initialize table */
func (t *TcpConnTable) Init() {
    t.table = make(map[net.Conn]TcpConn)
}

/* Add connection to table */
func (t *TcpConnTable) Add(conn net.Conn) TcpConn {
    /* create structure for incoming connection */ 
    tcpConn := newTcpConn()
    tcpConn.LocalAddr = conn.LocalAddr().String()
    tcpConn.RemoteAddr = conn.RemoteAddr().String()
    
    t.lock.Lock()
    t.table[conn] = tcpConn
    t.lock.Unlock()

    return tcpConn
}

/* Remove connection from table */
func (t *TcpConnTable) Remove(conn net.Conn) (TcpConn, error) {
    t.lock.Lock()
    tcpConn, ok := t.table[conn]
    
    /* delete from table    */
    delete(t.table, conn)
    t.lock.Unlock()
    
    if ok {
        return tcpConn, nil
    } else {
        return tcpConn, errors.New("no exist")
    }
}

/* whether connection exists in table?  */
func (t *TcpConnTable) Exist(conn net.Conn) bool {
    t.lock.Lock()
    _, ok := t.table[conn]
    t.lock.Unlock()
    return ok
}

/* get state of tcp connections */
func (t *TcpConnTable) GetState() []TcpConnState {
    states := make([]TcpConnState, 0)
    
    t.lock.Lock()
    
    for _, v := range t.table {
        state := TcpConnState{}
        state.QueueLen = v.SendQueue.Len()
        state.LocalAddr = v.LocalAddr
        state.RemoteAddr = v.RemoteAddr

        states = append(states, state)
    }    
    
    t.lock.Unlock()
    
    return states
}
