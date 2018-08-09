/* pop.go - get all pops */
/*
modification history
--------------------
2016/07/28, by Guang Yao, create
*/
/*
DESCRIPTION
*/

package trident

import (
	"encoding/json"
	"fmt"
)

import (
	"www.baidu.com/golang-lib/http_util"
)

const (
	// TODO: use uid for gtc
	URL_POP = "http://udi.baidu.com/q?domain=hb&handle=trident_idc_isp&uid=um"
)

/*
Format of message from URL_POP:
{
	"status": 0,
	"msg":
		{
			"idc_isp":
				[
					[u'HZ01', u'CHINANET'],
					[u'NJ01', u'CHINANET']
				]
		}
}
*/

type RawIdcIspItems [][]string // [ [pop_idc, pop_isp], [pop_idc, pop_isp], ...]

type IdcIspMsgBody struct {
	IdcIsps RawIdcIspItems `json:"idc_isp"`
}

type IdcIspMsg struct {
	Status  int           `json:"status"`
	Message IdcIspMsgBody `json:"msg"`
}

// type for return
type Pop struct {
	Idc string
	Isp string
}

func getPop(urlStr string) ([]Pop, error) {
	// prepare params
	params := map[string]string{
		"cmd": "idc_isp",
	}
	paramsStr, _ := json.Marshal(params) // not expect err

	// post to api
	resp, err := http_util.PostRR(urlStr, TIMEOUT, http_util.CONTENT_JSON, paramsStr, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.PostRR(): %s", err.Error())
	}

	// decode the result(json format)
	var result IdcIspMsg
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s, resp=%s", err.Error(), string(resp))
	}

	if result.Status != 0 {
		return nil, fmt.Errorf("status != 0: %d", result.Status)
	}

	// convert
	ret := make([]Pop, 0)
	for _, rawIdcIsp := range result.Message.IdcIsps {
		if len(rawIdcIsp) != 2 {
			return nil, fmt.Errorf("invalid format: %v", rawIdcIsp)
		}

		pop := Pop{rawIdcIsp[0], rawIdcIsp[1]}

		ret = append(ret, pop)
	}

	return ret, nil
}

// get all pops
//
// Returns:
//      - (all pops, err)
func GetPop() ([]Pop, error) {
	// request api for result
	return getPop(URL_POP)
}
