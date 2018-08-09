/* ret_msg_test.go - test of ret_msg.go */
/*
modification history
--------------------
2015/9/8, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package http_util

import (
	"encoding/json"
    "testing"
)

// test for requestGen(), without given (ipaddr, port)
func Test_RetMsgGenJson_Case1(t *testing.T) {
    // invoke requestGen()
    msgBuf := RetMsgGenJson("OK", "nothing happens")

	// decode the msg
	var msg interface{}
	err := json.Unmarshal(msgBuf, &msg)
	if err != nil {
		t.Errorf("json.Unmarshal():%s", err.Error())
		return
	}

	strMap, ok := msg.(map[string]interface{})
	if !ok {
		t.Errorf("fail to convert to map")
		return
	}
	value, ok := strMap["retCode"]
	if !ok {
		t.Errorf("no retCode")
		return
	}
	if value != "OK" {
		t.Errorf("retCode(%s) != OK", value)
		return
	}

	value, ok = strMap["msg"]
	if !ok {
		t.Errorf("no msg")
		return
	}
	if value != "nothing happens" {
		t.Errorf("msg(%s) != nothing happens", value)
		return
	}
}
