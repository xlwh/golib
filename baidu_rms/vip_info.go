/* vip_info.go - get expire/load balance/duty id/bns/bns strategy info for vips */
/*
modification history
--------------------
2016/07/08, by liuxiaowei07, create
*/
/*
DESCRIPTION
*/

package baidu_rms

import (
	"encoding/json"
	"fmt"
	"strings"
)

import (
	"www.baidu.com/golang-lib/http_util"
)

/*
Format of message from getVipExpire:
{
    "status":0,
	"data":
	    [
		    {
			    "ip":"111.13.100.249",
				"port":"80",
				"load_balance":"rr", // 负载均衡策略
				"bns":"group.bfe-yf.bfe.yf:main",  // bns
				"duty_id":"0",  // 值班表 ID
				"bns_strategy":"操作单",  // bns感知策略：操作单，半自动，全自动
				"expire":"2024-07-31"
			},
			...
		],
	"msg":"ok"
}
*/

// result info body for rms expire api
type ExpireInfoItem struct {
	Ip           *string // ipaddr for vip
	Port         *string // port for vip
	Load_balance *string // load balance, e.g., rr/wrr/...
	Expire       *string // expire date for vip, e.g., 2027-09-10 or 2020-08-19 16:03:00
	Bns          *string // bns for vip-port, if no binding bns, bns, bns_strategy and duty_id will be null
	Bns_strategy *string // whether the rms task of add/delete rs by bns for vip should auto pass
	Duty_id      *string // rms task of add/delete rs will send to the person on duty based on duty id
}

type ExpireData []ExpireInfoItem

// result struct of rms expire api
type ExpireResult struct {
	Status int    // 0 is success, 1 is failed
	Msg    string // ok | error info
	Data   ExpireData
}

func checkExpireData(result ExpireData) error {
	for _, expireInfo := range result {
		if expireInfo.Ip == nil || expireInfo.Port == nil || expireInfo.Load_balance == nil ||
			expireInfo.Expire == nil {

			return fmt.Errorf("wrong format")
		}
	}

	return nil
}

func getExpireByVips(url string) (ExpireData, error) {
	// request api
	resp, err := http_util.Read(url, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	// Note: in case status != 0, data will "", a decode error will be returned
	result := ExpireResult{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	// rms err
	if result.Data == nil {
		return nil, fmt.Errorf("rms err: no data return")
	}

	if err = checkExpireData(result.Data); err != nil {
		return nil, fmt.Errorf("rms err: %s", err.Error())
	}

	return result.Data, nil
}

/* GetExpireByVips - get expire/load balance/duty id/bns/bns strategy info by the vips
 *
 * Params:
 *      - vips: vip list
 *
 * Returns:
 *      - (expire/load balance/duty id/bns/bns strategy info for vips, err)
 */
func GetExpireByVips(vips []string) (ExpireData, error) {
	// generate url
	url := "http://rms.baidu.com/index.php?r=interface/api&handler=getVipExpire&ips=" + strings.Join(vips, ",")

	// request api for result
	return getExpireByVips(url)
}
