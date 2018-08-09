/* pop_traffic.go - get all <idc, isp> */
/*
modification history
--------------------
2016/07/28, by Guang Yao, create
*/
/*
DESCRIPTION
TODO: use the formal interface for GTC
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
	URL_POP_TRAFFIC = "http://udi.baidu.com/q?domain=hb&handle=trident_ext_bw&uid=um"
)

/*
Format of message from POP_TRAFFIC_ISP:
NOTE: the format may diff when the interface upgrade
{
	"status": 0,
	"time": 1461053900,
	"msg":	{
		"ext_bw": [
			{
				"idc": "QD01",
				"isp": "unicom",
				"abnormal": 0,
				"external_port": {
	            	"out": {"used": 295400917182, "total": 304000000000, "predicted": 300400917182, "congested": 0},
	                "in": {"used": 295400917182, "total": 304000000000, "predicted": 300400917182, "congested": 0}
	            },
	     		"product_info": {
	                "out": {"used": 295400917182, "predicted": 300400917182},
	                "in": {"used": 295400917182, "predicted": 300400917182}
	            },
	     		"cost_info": {
	                   "physical_line": {
	                          "max_bw": 304000000000,
	                          "out": {"available_bw": 295400917182},
	                          "in": {"available_bw": 295400917182},},
	                   "idc_cost_line": {
	                          "max_bw": 304000000000,
	                          "out": {"available_bw": 295400917182},
	                          "in": {"available_bw": 295400917182},},
	                   "product_cost_line": {
	                          "max_bw": 304000000000,
	                          "out": {"available_bw": 295400917182},
	                          "in": {"available_bw": 295400917182},},
	            }
	        }
	    ]
	}
}
*/

type RawPortTrafficItem struct {
	Used      int64 `json:"used"`      // current pop traffic; in bps
	Total     int64 `json:"total"`     // current pop bandwidth; in bps
	Predicted int64 `json:"predicted"` // predict pop traffic; in bps
	Congested int   `json:"congested"` // 0: not congested; 1: congested
}

type RawPortTraffic struct {
	OutTraffic RawPortTrafficItem `json:"out"` // out total traffic
	InTraffic  RawPortTrafficItem `json:"in"`  // in total traffic
}

type RawProductTrafficItem struct {
	Used      int64   `json:"used"`      // traffic for product with the given uid; in bps
	Predicted float64 `json:"predicted"` // predict traffic for product with the given uid; in bps
}

type RawProductInfo struct {
	OutTraffic RawProductTrafficItem `json:"out"` // out traffic for product
	InTraffic  RawProductTrafficItem `json:"in"`  // in traffic for product
}

type RawBwCostItem struct {
	AvailableBw float64 `json:"available_bw"` // remain bandwidth below costline; in bps
}

type RawCostInfoItem struct {
	MaxBw      int64         `json:"max_bw"` // max bandwidth of the pop; in bps
	OutTraffic RawBwCostItem `json:"out"`    // out costline
	InfTraffic RawBwCostItem `json:"in"`     // in costline
}

type RawCostInfo struct {
	PhysicalLine    RawCostInfoItem `json:"physical_line"`     // max bandwidth of the pop
	IdcCostLine     RawCostInfoItem `json:"idc_cost_line"`     // cost line for the pop
	ProductCostLine RawCostInfoItem `json:"product_cost_line"` // product cost line for the pop
}

// traffic item
type ExtBwItem struct {
	Idc          string         `json:"idc"`           // pop idc
	Isp          string         `json:"isp"`           // pop isp
	Abnormal     int            `json:"abnormal"`      // whether the pop is normal; 0 is normal
	ExternalPort RawPortTraffic `json:"external_port"` // total traffic on the pop
	ProductInfo  RawProductInfo `json:"product_info"`  // product traffic on the pop
	CostInfo     RawCostInfo    `json:"cost_info"`     // cost info for the pop
}

type PopTrafficMsgBody struct {
	ExtBw []ExtBwItem `json:"ext_bw"`
}

type PopTrafficMsg struct {
	Status  int               `json:"status"`
	Time    int64             `json:"time"`
	Message PopTrafficMsgBody `json:"msg"`
}

// struct for return
type PopTraffic struct {
	Idc          string
	Isp          string
	Time         int64
	InTraffic    int64 // in Kbps
	InBandwidth  int64 // in Kbps
	OutTraffic   int64 // in Kbps
	OutBandwidth int64 // in Kbps
	CostLine     int64 // in Kbps
}

func getPopTraffic(urlStr string, pops []Pop) ([]PopTraffic, error) {
	if len(pops) == 0 {
		return nil, fmt.Errorf("no pops")
	}

	// prepare params of pops
	// format: [{"idc": "QD01", "isp": "unicom"}]
	popsParams := make([]map[string]string, 0)
	for _, pop := range pops {
		popParam := map[string]string{
			"idc": pop.Idc,
			"isp": pop.Isp,
		}

		popsParams = append(popsParams, popParam)
	}

	// format all params:
	//  {"cmd": "ext_bw", "params": [{"idc": "QD01", "isp": "unicom"}], "uid": "pcs/bfe"}
	params := map[string]interface{}{
		"cmd":    "ext_bw",
		"params": popsParams,
		"uid":    "pcs", // TODO: use uid for bfe
	}
	paramsStr, _ := json.Marshal(params) // not expect err

	// post to api
	resp, err := http_util.PostRR(urlStr, TIMEOUT, http_util.CONTENT_JSON, paramsStr, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.PostRR(): %s", err.Error())
	}

	// decode the result(json format)
	var result PopTrafficMsg
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s, resp=%s", err.Error(), string(resp))
	}

	if result.Status != 0 {
		return nil, fmt.Errorf("status != 0: %d", result.Status)
	}

	// convert
	popTraffics := convertTraffic(result.Time, result.Message.ExtBw)

	return popTraffics, nil
}

// convert raw traffic to final result
func convertTraffic(time int64, bwItems []ExtBwItem) []PopTraffic {
	popTraffics := make([]PopTraffic, 0)

	for _, bwItem := range bwItems {
		var popTraffic PopTraffic
		popTraffic.Idc = bwItem.Idc
		popTraffic.Isp = bwItem.Isp
		popTraffic.Time = time
		popTraffic.InTraffic = bwItem.ExternalPort.InTraffic.Used / 1024
		popTraffic.OutTraffic = bwItem.ExternalPort.OutTraffic.Used / 1024
		popTraffic.InBandwidth = bwItem.ExternalPort.InTraffic.Total / 1024
		popTraffic.OutBandwidth = bwItem.ExternalPort.OutTraffic.Total / 1024

		// determine costline
		inCostline := int64(bwItem.CostInfo.IdcCostLine.InfTraffic.AvailableBw)/1024 +
			popTraffic.InTraffic
		outCostline := int64(bwItem.CostInfo.IdcCostLine.OutTraffic.AvailableBw)/1024 +
			popTraffic.OutTraffic
		if inCostline > outCostline {
			popTraffic.CostLine = inCostline
		} else {
			popTraffic.CostLine = outCostline
		}

		popTraffics = append(popTraffics, popTraffic)
	}

	return popTraffics
}

// get traffic for the given <idc, isp>
//
// Returns:
//      - (pop traffics, err)
func GetPopTraffic(pops []Pop) ([]PopTraffic, error) {
	return getPopTraffic(URL_POP_TRAFFIC, pops)
}
