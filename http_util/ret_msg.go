/* ret_msg.go - generate return msg in json */
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
)

// generate the json string to return
func RetMsgGenJson(retCode string, msg string) []byte {
	m := make(map[string]string)
	m["retCode"] = retCode
	m["msg"] = msg

	msgJson, err := json.Marshal(m)
	if err != nil {
		// should never get here; alway return a result
		return []byte{0}
	}

	return msgJson
}