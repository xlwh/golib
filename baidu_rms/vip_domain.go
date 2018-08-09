/* vip_domain.go - get domain info for vip */
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
Format of message from getBnsToVip:
[
    {
	    "dns":"cb.e.shifen.com",
		"ip":"58.217.200.81",
		"type":"bgw",
		"attr":"outside",
		"hostname":"hz01-s-bfe00.hz01",
		"ip_in1":"10.212.71.36",
		"ip_in2":null,
		"ip_department":"OP",
		"ip_product":"BFE",
		"department":"OP",
		"product":"BFE",
		"user":"yangsijie"
	},
	...
]
*/

// one domain info for vip-rs
type DnsByVipItem struct {
	Dns           *string // domain which vip has binded
	Ip            *string // ipaddr of vip
	Type          *string // type of vip, bgw/bignat/server/...
	Attr          *string // attribute of vip, outside/inside
	Hostname      *string // rs hostname of vip
	Ip_in1        *string // rs ipaddr of vip
	Ip_in2        *string // rs ipaddr of vip when rs has two ipaddr, otherwise is null
	Ip_department *string // department of vip
	Ip_product    *string // product of vip
	Department    *string // user's department of vip
	Product       *string // user's product of vip
	User          *string // person in charge of vip
}

// all domain info for vip
type DnsByVip []DnsByVipItem

// check DnsByVip
func checkDnsByVip(dnsInfoList DnsByVip) error {
	for _, dnsInfo := range dnsInfoList {
		if dnsInfo.Dns == nil || dnsInfo.Ip == nil || dnsInfo.Type == nil || dnsInfo.Attr == nil ||
			dnsInfo.Hostname == nil || dnsInfo.Ip_in1 == nil || dnsInfo.Ip_department == nil ||
			dnsInfo.Ip_product == nil || dnsInfo.Department == nil || dnsInfo.Product == nil ||
			dnsInfo.User == nil {

			return fmt.Errorf("wrong format")
		}
	}
	return nil
}

func getDnsByVip(urlStr string) (DnsByVip, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	result := DnsByVip{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	if err = checkDnsByVip(result); err != nil {
		return nil, fmt.Errorf("rms err: %s", err.Error())
	}

	return result, nil
}

/* GetDnsByVip - get dns by the vip
 *
 * Params:
 *      - vip: vip
 *
 * Returns:
 *      - (domain list for vip, err)
 */
func GetDnsByVip(vip string) (DnsByVip, error) {
	// generate url
	url := "http://rms.baidu.com/?r=interface/rest&handler=getDnsIpInfo&ipout=" + vip

	return getDnsByVip(url)
}
