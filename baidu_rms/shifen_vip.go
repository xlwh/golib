/* shifen_vip.go - get vip info for shifen */
/*
modification history
--------------------
2016/08/03, by wuzhenxing, create
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
  "status": 0,
  "data": {
    "name": "jpaas-movie.e.shifen.com",
    "service_name": "movie",
    "rms_product": "nuomi",
    "noah_path": "BAIDU_LBS_nuomi",
    "rms_department": "LBS",
    "ipInfo": [
      {
        "ip": "123.125.115.198",
        "status": "reserve",
        "idc": "TC",
        "isp": "联通"
      }
	]
  },
  "msg": "ok",
  "cache": "ok"
}
*/

// response for getVipIsp API
type shifenVipBody struct {
	Status int           `json:"status"`
	Data   shifenVipData `json:"data"`
	Msg    string        `json:"msg"`
	Cache  string        `json:"cache"`
}

type shifenVipData struct {
	Name          string   `json:"name"`
	ServiceName   string   `json:"service_name"`
	RmsProduct    string   `json:"rms_product"`
	NoahPath      string   `json:"noah_path"`
	RmsDepartment string   `json:"rms_department"`
	IpInfo        []ipInfo `json:"ipInfo"`
}

type ipInfo struct {
	Ip     string `json:"ip"`
	Status string `json:"status"`
	Idc    string `json:"idc"`
	Isp    string `json:"isp"`
}

type ShifenVip struct {
	Name    string
	IpInfos []IpInfo
}

type IpInfo struct {
	Ip     string
	Status string
}

func getVipsByShifen(shifen string, url string) (*ShifenVip, error) {
	dataBuff, err := http_util.Read(url, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %v", err)
	}

	shifenVipRaw := shifenVipBody{}
	err = json.Unmarshal(dataBuff, &shifenVipRaw)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	if shifenVipRaw.Status != 0 {
		return nil, fmt.Errorf("shifenVipRaw status %d", shifenVipRaw.Status)
	}

	if shifenVipRaw.Data.Name != shifen {
		return nil, fmt.Errorf("shifen name in shifenVip Raw [%s] conflict with shifen [%s]",
			shifenVipRaw.Data.Name, shifen)
	}

	shifenVip := ShifenVip{}
	shifenVip.Name = shifenVipRaw.Data.Name

	ipInfos := []IpInfo{}
	repeatedCheckMap := map[string]bool{}
	for _, ipInfoRaw := range shifenVipRaw.Data.IpInfo {
		if _, ok := repeatedCheckMap[ipInfoRaw.Ip]; ok {
			return nil, fmt.Errorf("ip info repeated [%s]", ipInfoRaw.Ip)
		}
		repeatedCheckMap[ipInfoRaw.Ip] = true

		ipInfo := IpInfo{
			Ip:     ipInfoRaw.Ip,
			Status: ipInfoRaw.Status,
		}
		ipInfos = append(ipInfos, ipInfo)
	}
	shifenVip.IpInfos = ipInfos

	return &shifenVip, nil
}

/* GetVipsByShifen - get vip info by shifen
 *
 * Params:
 *      - shifen: shifen
 *
 * Returns:
 *      - (vip info of shifen, err)
 */
func GetVipsByShifen(shifen string) (*ShifenVip, error) {
	// generate url
	url := "http://rms.baidu.com/?r=interface/rest&handler=getDomainInfoByName&domain=" + shifen

	// request api for result
	return getVipsByShifen(shifen, url)
}
