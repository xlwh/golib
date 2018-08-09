/* vip.go - get vip info from rms interface */
/*
modification history
--------------------
2015/12/21, by Guang Yao, create based on bfetools/vip-monitor of zhangjiyang
2016/06/16, by Xiaowei Liu, modify: add timeout for http request
2016/07/20, by Xiaowei Liu, modify: add getVipInfo for op-bfe(department=op&product=bfe)
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
Format of message from getVipInfoByRs:
{
    "status":0,
    "message":
        {
            "220.181.111.127":
                {
                    "port":["80", "443"],
                    "port_end":["80", "443"],
                    "type":"bgw"
                },
            ...
        }
}
*/

type VipByRsItem struct {
	Port     []string
	Port_end []string
	Type     string // vip type
}

type vipByRsMsgBody map[string]VipByRsItem

type vipByRsMsg struct {
	Status  int
	Message vipByRsMsgBody
}

func getVipByRs(urlStr string) (map[string]VipByRsItem, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	// Note: in case status != 0, a decode error will be returned
	result := vipByRsMsg{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	if result.Message == nil {
		return nil, fmt.Errorf("rms err: no data return")
	}

	return result.Message, nil
}

// get vip by the given real server
// Params:
//      - rsName: name of the real server
// Returns:
//      - (vip address => VipByRsItem, err)
func GetVipByRs(rsName string) (map[string]VipByRsItem, error) {
	// generate url
	url := "http://rms.baidu.com/?r=interface/api&handler=getVipInfoByRs&hostname=" + rsName

	// request api for result
	return getVipByRs(url)
}

/*
Format of message from getVipInfoByBns:
{
    "status": 0,
    "data": [
        {
            "group.bfe-yf.bfe.yf": "220.181.111.23"
        },
        {
            "group.bfe-yf.bfe.yf": "111.13.12.143"
        },
        {
            "group.bfe-yf.bfe.yf": "220.181.163.194"
        }
}
*/

// vip Item
type vipByBnsItem map[string]string
type vipByBnsItems []vipByBnsItem

type vipByBnsMsgBody struct {
	Status int
	Data   vipByBnsItems
}

func getVipByBns(urlStr, bnsName string) ([]string, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode response body from rms
	// Note: in case status != 0, a decode error will be returned
	result := vipByBnsMsgBody{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	// rms err
	if result.Data == nil {
		return nil, fmt.Errorf("rms err: no data return")
	}

	// convert VipItem into human readable viplist
	vips := make([]string, len(result.Data))
	for _, vipItems := range result.Data {
		vip, ok := vipItems[bnsName]
		if !ok {
			return nil, fmt.Errorf("rms err: wrong response format")
		}

		vips = append(vips, vip)
	}

	return vips, nil
}

// get vip list by the given bnsName
// Params:
//      - bnsName: name of bns
// Returns:
//      - (viplist, err)
func GetVipByBns(bnsName string) ([]string, error) {
	// generate url
	url := "http://rms.baidu.com/?r=interface/rest&handler=getBnsToVip&type=get_vip&data=" + bnsName

	return getVipByBns(url, bnsName)
}

/*
Format of message from getVipInfo:
{
    "status":0,
	"msg":"",
	"data":
	    {
		    "220.181.111.191":
			    {
				    "idc":"M1",
					"isp":"\u7535\u4fe1",
					"port":
					    {
						    "80":"group.bfe-yf.bfe.yf"
						},
					"domains":
					    [
						    "tieba.com"
						]
				},
			......
		}
}
*/
type VipInfoItem struct {
	Idc     string            `json:"idc"`
	Isp     string            `json:"isp"`
	Port    map[string]string `json:"port"`
	Domains []string          `json:"domains"`
}

type VipInfoData struct {
	Msg  string                  `json:"msg"`
	Data map[string]*VipInfoItem `json:"data"` // vip => VipInfoItem, 111.111.111.111 => VipInfoItem
}

type VipByNoahPath struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   VipInfoData `json:"data"`
}

func getVipByProductNoahPath(urlStr string, token string) (map[string]*VipInfoItem, error) {
	// request api
	header := map[string]string{
		"ROP-Authorization": "RopAuth " + token,
	}
	resp, err := http_util.Read(urlStr, LONG_TIME_OUT, header)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode response body from rms
	// Note: in case status != 0, a decode error will be returned
	result := VipByNoahPath{}
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

	return result.Data.Data, nil
}

func GetVipByProductNoahPath(appKey string, secretToken string, user string,
	productNoahPath string, testUrl *string) (map[string]*VipInfoItem, error) {
	// for test
	if testUrl != nil {
		return getVipByProductNoahPath(*testUrl, "")
	}

	// get token
	token, err := GetRmsToken(appKey, secretToken, user)
	if err != nil {
		return nil, fmt.Errorf("GetRmsToken() err: %s", err.Error())
	}

	// generate url
	url := "http://api.rms.baidu.com/v1/network_ip/getVipInfoByNoahProductPath?path=" + productNoahPath

	// request api for result
	return getVipByProductNoahPath(url, *token)

}
