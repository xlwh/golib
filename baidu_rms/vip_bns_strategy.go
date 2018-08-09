/* vip_bns_strategy.go - get bns strategy and strategy threshold for vip */
/*
modification history
--------------------
2016/06/24, by liuxiaowei07, create
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
Format of message from getBnsVipStrategyThreshold:
{
    "status":0,
	"msg":"success",
	"data":
	    {
		    "bns":"group.bfe-yf.bfe.yf:main",
			"vip":"111.13.105.50",
			"port":80,
			"status":"不在流程中 | 不在流程中",
			"bns_strategy":"0", //为bns感知策略：0 为操作单，1 为半自动， 2为全自动
			"strategy_threshold":"100" // bns感知策略为1时，当rms根据bns自动调整VIP的rs的变动不超过该阈值，则自动触发单无需审核直接通过
		}
}
*/

// bns strategy info
type BnsStrategyData struct {
	Bns                *string // bns for vip-port
	Vip                *string // ipaddr for vip
	Port               *int    // port for vip
	Status             *string // whether the vip is in proccesing
	Bns_strategy       *string // whether the rms task of add/delete rs by bns for vip should auto pass
	Strategy_threshold *string // when bns_strategy is 1, if change rate is less than threshold, rms task could auto pass
}

// result struct of bns strategy api
type BnsStrategy struct {
	Status int    // 0 is success, 1 is failed
	Msg    string // success | error info
	Data   *BnsStrategyData
}

func checkBnsStrategyData(result *BnsStrategyData) error {
	if result.Bns == nil || result.Vip == nil || result.Port == nil || result.Status == nil ||
		result.Bns_strategy == nil || result.Strategy_threshold == nil {

		return fmt.Errorf("wrong format")
	}

	return nil
}

func getBnsStrategyByVip(url string) (*BnsStrategyData, error) {
	// request api
	resp, err := http_util.Read(url, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	result := BnsStrategy{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	// rms err
	if result.Status != 0 {
		return nil, fmt.Errorf("rms err: %s", result.Msg)
	}
	if result.Data == nil {
		return nil, fmt.Errorf("rms err: no data return")
	}

	if err = checkBnsStrategyData(result.Data); err != nil {
		return nil, fmt.Errorf("rms err: %s", err.Error())
	}

	return result.Data, nil
}

/* GetBnsStrategyByVip - get bns strategy info by the vip
 *
 * Params:
 *      - vip: ipaddr for vip
 *      - port: port for vip
 *      - bns: bns for vip-port
 *
 * Returns:
 *      - (bns strategy info for vip, err)
 */
func GetBnsStrategyByVip(vip, port, bns string) (*BnsStrategyData, error) {
	// generate url
	url := fmt.Sprintf("http://rms.baidu.com/index.php?r=interface/api&handler=getBnsVipStrategyThreshold&bns=%s&vip=%s&port=%s", bns, vip, port)

	// request api for result
	return getBnsStrategyByVip(url)
}
