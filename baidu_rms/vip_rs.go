/* vip_rs.go - get rs info of vip from rms */
/*
modification history
--------------------
2016/06/13, by liuxiaowei07, create based on bfetools/vip-monitor of zhangjiyang
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
Format of message from getRsInfo:
{
    "status":0,
	"message":
	    {
		    "vip":"180.149.131.103",
			"status":"1",
			"type":"bgw",
			"idc":"YF",
			"port_info_list":
			    [
				    {
					    "vs_port":"90",
						"vs_port_end":"90",
						"port_type":"TCP",
						"rs_info_list":
						    [
							    {
								    "rs_ip":"10.38.171.51",
									"rs_port":"8900",
									"server_id":"117636",
									"weight":"1"
								},
								...
							]
					},
					...
				]
		}
}
*/

// data struct for decode json result
type RsInfo struct {
	Rs_ip     string `json:"rs_ip"`     // real server ip address
	Rs_port   string `json:"rs_port"`   // real server Port
	Server_id string `json:"server_id"` // real Server id in rms
	Weight    string `json:"weight"`    // weight of load balance
}
type RsInfoList []RsInfo

type PortInfo struct {
	Vs_port      string     `json:"vs_port"`      // vip port
	Vs_port_end  string     `json:"vs_port_end"`  // vip port end
	Port_type    string     `json:"port_type"`    // port type tcp/udp
	Rs_info_list RsInfoList `json:"rs_info_list"` // realServer list
}

type PortInfoList []PortInfo

type VipRsMessage struct {
	Vip            string       `json:"vip"`    // vip address
	Status         string       `json:"status"` // vip's status
	Type           string       `json:"type"`   // vip's type
	Idc            string       `json:"idc"`    // vip's idc
	Port_info_list PortInfoList `json:"port_info_list"`
}

type VipRsInfo struct {
	Message *VipRsMessage // rs info body when success
	Status  int           // 0 is success
}

func getRsByVip(urlStr string) (*VipRsMessage, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	// Note: in case status != 0, a decode error will be returned
	result := &VipRsInfo{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	if result.Message == nil {
		return nil, fmt.Errorf("rms err: no data return")
	}

	return result.Message, nil
}

/* GetRsByVip - get rs by the vip
 *
 * Params:
 *      - vip: vip
 *
 * Returns:
 *      - (rs info for vip, err)
 */
func GetRsByVip(vip string) (*VipRsMessage, error) {
	// generate url
	url := "http://rms.baidu.com/?r=interface/api&handler=getRsInfo&vip=" + vip

	// request api for result
	return getRsByVip(url)
}
