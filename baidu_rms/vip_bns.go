/* vip_bns.go - get bns info for vip */
/*
modification history
--------------------
2016/06/13, by liuxiaowei07, create
*/
/*
DESCRIPTION
    the api need set IP white list in open.rms.baidu.com
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
Format of message from getBnsToVip:
{
    "status":0,
	"data":
	    {
		    "data":
			    {
				    "111.202.114.38":
					    [
						    {
							    "bns":"group.bfe-tc.bfe.tc:main",
								"port":"80"
							},
							{
							    "bns":"group.bfe-tc.bfe.tc:https",
								"port":"443"
							}
						]
					...
				}
		}
}
*/

// bns info for one port
type BnsByVipItem struct {
	Bns  *string // bns name
	Port *string // port
}

// bns info for all ports of vip
type BnsByVipItems []BnsByVipItem

// bns info of all vips, vip => bns info for all port of vip
type BnsByVip map[string]BnsByVipItems

type BnsByVipDataBody struct {
	Data BnsByVip // bns info body
}

type BnsByVipResult struct {
	Status int              // 0 is success
	Msg    string           // error msg when failed
	Data   BnsByVipDataBody // bns info body
}

// check BnsByVip
func checkBnsByVip(bnsInfoMap BnsByVip) error {
	for _, bnsInfo := range bnsInfoMap {
		for _, bns := range bnsInfo {
			if bns.Bns == nil || bns.Port == nil {
				return fmt.Errorf("wrong format")
			}
		}
	}
	return nil
}

/* getBnsByVips - get bns by the vips
 *
 * Params:
 *      - url: url
 *      - token: rms token
 *
 * Returns:
 *      - (bns info for all vips, err)
 */
func getBnsByVips(url string, token string) (BnsByVip, error) {
	// request api
	header := map[string]string{
		"ROP-Authorization": "RopAuth " + token,
	}
	resp, err := http_util.Read(url, TIME_OUT, header)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	result := BnsByVipResult{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	// rms err
	if result.Status != 0 {
		return nil, fmt.Errorf("rms err: %s", result.Msg)
	}
	if result.Data.Data == nil {
		return nil, fmt.Errorf("rms err: no data return")
	}
	if err = checkBnsByVip(result.Data.Data); err != nil {
		return nil, fmt.Errorf("rms err: %s", err.Error())
	}

	return result.Data.Data, nil
}

/* GetBnsByVips - get bns by the vips
 *
 * Params:
 *      - vips: vip list
 *      - app_key: app_key applied from open.rms.baidu.com
 *      - secret_token: secret_token applied from open.rms.baidu.com
 *      - user: user of special_token applied from open.rms.baidu.com
 *      - testUrl: used when unit/integration test to mock httpserver
 *
 * Returns:
 *      - (bns info for all vips, err)
 */
func GetBnsByVips(vips []string, appKey string, secretToken string, user string, testUrl *string) (BnsByVip, error) {
	// for test
	if testUrl != nil {
		return getBnsByVips(*testUrl, "")
	}

	// get token
	token, err := GetRmsToken(appKey, secretToken, user)
	if err != nil {
		return nil, fmt.Errorf("GetRmsToken() err: %s", err.Error())
	}

	// generate url
	url := "http://api.rms.baidu.com/v1/network_ip/getBnsToVip?type=get_bns_name&data=" + strings.Join(vips, ",")

	// request api for result
	return getBnsByVips(url, *token)
}
