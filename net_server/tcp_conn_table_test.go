/* tcp_conn_table_test.go - test for tcp_conn_table.go  */
/*
modification history
--------------------
2014/3/11, by Zhang Miao, create
2014/8/6, by Zhang Miao, move from waf_server
*/
/*
DESCRIPTION
*/
package net_server

import (
    "net"
    "testing"
)

import "www.baidu.com/golang-lib/log"

func TestTcpConnTableAdd(t *testing.T) {
    log.Init("test", "DEBUG", "./log", true, "D", 5)
    
    var tcpTable TcpConnTable

    /* initialize tcp table */
    tcpTable.Init()
        
    conn, err := net.Dial("tcp", "svn.baidu.com:http")
    
    if err != nil {
        t.Error("fail to make connection to svn.baidu.com")
    }
    
    // add to table
    tcpTable.Add(conn)
    
    // remove from table
    _, err = tcpTable.Remove(conn)
    if err != nil {
        t.Error("err in tcpTable.Remove()")
    }
    
    log.Logger.Close()
}
