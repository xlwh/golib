/* rs.go - get rs info from rms interface */
/*
modification history
--------------------
2016/07/20, by Xiaowei Liu, create
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
Format of message from getRsInfoByHosts:
{
    "status":0,
	"message":
	    [
		    {
		        "host":"yf-s-bfe00.yf01",  // hostname of rs
				"id":"363407",             // id of rs in rms
				"ip":"10.38.159.43",       // ip addr of rs
				"type":"server_info"       // type of rs
			},
			...
		]
}
*/

type RsInfoByHostItem struct {
	Host *string // rs hostname
	Ip   *string // rs ip
	Id   *string // rs id
	Type *string // rs type
}

type RsInfoByHost struct {
	Status  int
	Message []RsInfoByHostItem
}

func checkRsInfoByHost(rsInfoList []RsInfoByHostItem) error {
	for _, rsInfo := range rsInfoList {
		if rsInfo.Host == nil || rsInfo.Ip == nil || rsInfo.Id == nil || rsInfo.Type == nil {
			return fmt.Errorf("wrong format: %v", rsInfo)
		}
	}
	return nil
}

func getRsInfoByHost(urlStr string) ([]RsInfoByHostItem, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	// Note: in case status != 0, a decode error will be returned
	result := RsInfoByHost{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	if result.Message == nil {
		return nil, fmt.Errorf("rms err: no data return")
	}

	if err = checkRsInfoByHost(result.Message); err != nil {
		return nil, fmt.Errorf("checkRsInfoByHost(): %s", err.Error())
	}

	return result.Message, nil
}

// get rs info (include id and type) by hostname of rs
// Params:
//      - rsNames: hostname list of the real servers
// Returns:
//      - (rs info list, err)
func GetRsInfoByHost(rsNames []string) ([]RsInfoByHostItem, error) {
	// generate url
	url := `http://rms.baidu.com/?r=interface/api&handler=getRsInfoByHosts&hosts=["` +
		strings.Join(rsNames, `","`) + `"]`

	// request api for result
	return getRsInfoByHost(url)
}

/*
more information of api in: http://rms.baidu.com/?r=interface/rest&handler=searchServers
Format of message from searchServers:
[
    {
	    "id":"363407",                  // id of rs in rms
		"idc":"YF",                     // idc of rs
		"idc_id":"147",                 // idc id of rs
		"hostname":"yf-s-bfe00.yf01",   // hostname of rs
		"ip_in1":"10.38.159.43"         // ip addr of rs
	},
	...
]
*/

type RsInfoByIpItem struct {
	Hostname *string // rs hostname
	Ip_in1   *string // rs ip
	Id       *string // rs id
	Idc      *string // rs idc
	Idc_id   *string // rs idc id
}

func checkRsInfoByIp(rsInfoList []RsInfoByIpItem) error {
	for _, rsInfo := range rsInfoList {
		if rsInfo.Hostname == nil || rsInfo.Ip_in1 == nil || rsInfo.Id == nil ||
			rsInfo.Idc == nil || rsInfo.Idc_id == nil {
			return fmt.Errorf("wrong format: %v", rsInfo)
		}
	}
	return nil
}

func getRsInfoByIp(urlStr string) ([]RsInfoByIpItem, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode response body from rms
	result := []RsInfoByIpItem{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	if err = checkRsInfoByIp(result); err != nil {
		return nil, fmt.Errorf("checkRsInfoByIp(): %s", err.Error())
	}

	return result, nil
}

// get rs info (include id, idc and idc id) by ip of rs
// Params:
//      - ips: ip list of the real servers
// Returns:
//      - (rs info list, err)
func GetRsInfoByIp(rsIps []string) ([]RsInfoByIpItem, error) {
	// generate url
	url := "http://rms.baidu.com/?r=interface/rest&handler=searchServers&" +
		"return_type=json&show_fields=id,idc,idc_id,hostname,ip_in1&ip_in1=" +
		strings.Join(rsIps, ",")

	return getRsInfoByIp(url)
}
