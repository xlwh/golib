/* vip_isp.go - get isp info for vip */
/*
modification history
--------------------
2016/07/20, by liuxiaowei07, create
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
Format of message from getVipIsp:
{
    "status":0,
	"msg":"Ok",
	"data":
	    {
		    "id":"1",
			"vip":"220.181.111.111",
			"name":"电信"
		}
}
*/

// response for getVipIsp API
type VipIspData struct {
	Id   *string // value will be 1
	Vip  *string // ip addr
	Name *string // isp name
}
type IspByVip struct {
	Status int         // status, 0 is success, 1 is failed
	Msg    string      // Ok if success, err info if failed
	Data   *VipIspData // isp info
}

func getIspByVip(urlStr string) (*string, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	// Note: in case status != 0, a decode error will be returned
	result := IspByVip{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	if result.Data == nil || result.Data.Name == nil {
		return nil, fmt.Errorf("rms error: no isp info return")
	}

	return result.Data.Name, nil
}

/* GetIspByVip - get isp by the vip
 *
 * Params:
 *      - vip: vip
 *
 * Returns:
 *      - (isp of vip, err)
 */
func GetIspByVip(vip string) (*string, error) {
	// generate url
	url := "http://rms.baidu.com/index.php?r=interface/api&handler=getVipIsp&vip=" + vip

	// request api for result
	return getIspByVip(url)
}
