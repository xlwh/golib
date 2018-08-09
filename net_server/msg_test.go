/* msg_test.go - test for msg.go */
/*
modification history
--------------------
2014/3/12, by Zhang Miao, create
2014/8/6, by Zhang Miao, move from waf_server
*/
/*
DESCRIPTION
*/
package net_server

import (
    "testing"
)

import (
    "www.baidu.com/golang-lib/log"
)    

func TestHeaderEncode(t *testing.T) {
    log.Init("test", "DEBUG", "./log", true, "D", 5)

    header := MsgHeader{[4]byte{1, 2, 3, 4}, 5, 6, 7, 8}
        
    /* try to encode    */
    buf, err := MsgHeaderEncode(header)
        
    if err != nil {
        t.Error("fail in MsgHeaderEncode")
    }

    /* try to decode    */
    var headerNew MsgHeader
    headerNew, err = MsgHeaderDecode(buf)
    
    if err != nil {
        t.Error("fail in MsgHeaderDecode:", err.Error())
    }
    
    if header != headerNew {
        t.Error("two headers should be equal")
    }

    /* test MsgHeaderLen()   */
    len := MsgHeaderLen()

    if len != 16 {
        t.Error("length of MsgHeader should be 16")
    }
        
    log.Logger.Close()
}
