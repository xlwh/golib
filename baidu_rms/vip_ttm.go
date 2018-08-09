/* vip_ttm.go - get ttm info of vip from rms */
/*
modification history
--------------------
2016/06/13, by liuxiaowei07, create
*/
/*
DESCRIPTION
*/
package baidu_rms

import (
	"encoding/json"
	"fmt"
)

import (
	"www.baidu.com/golang-lib/http_util"
)

/*
Format of message from getTtmValid:
{
    "status":0,
	"msg":"找不到对应的VIP",  // if status != 0
	"data":
	    {
		    "error_code":0,
			"error_msg":"success",
			"ttm_vip_valid":0
		}
}
*/

// data struct for decode json result
type vipTtmData struct {
	Code int    `json:"error_code"`    // error_code, 0 is ok
	Msg  string `json:"error_msg"`     // error_msg
	Ttm  int    `json:"ttm_vip_valid"` // 0/1, 0 is not transmit transparently
}

type vipTtmInfo struct {
	Data   *vipTtmData `json:"data"`   // ttm info body when success
	Msg    string      `json:"msg"`    // error message when failed
	Status int         `json:"status"` // 0 is success
}

func getTtmByVip(urlStr string) (*bool, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	result := &vipTtmInfo{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	// err in rms
	if result.Status != 0 {
		return nil, fmt.Errorf("rms err: %s", result.Msg)
	}
	if result.Data.Code != 0 {
		return nil, fmt.Errorf("rms err: %s", result.Data.Msg)
	}

	// result is ok
	flag := false
	if result.Data.Ttm == 1 {
		flag = true
	}
	return &flag, nil
}

/* GetTtmByVip - get ttm info by the vip
 *
 * Params:
 *      - vip: vip
 *
 * Returns:
 *      - (ttm status for vip, err)
 */
func GetTtmByVip(vip, port string) (*bool, error) {
	// request rms for result
	url := fmt.Sprintf("http://rms.baidu.com/index.php?r=interface/api&handler=getTtmValid&vip=%s&port=%s", vip, port)

	// request api for result
	return getTtmByVip(url)
}
