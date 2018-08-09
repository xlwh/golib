/* rms_task.go - check status of rms task about vip */
/*
modification history
--------------------
2016/06/15, by liuxiaowei07, create
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
Format of message from getBgwListCurrentStep:
{
    "status":0,
	"msg":"ok",
	"data":
	    {
		    "112.80.248.52":
			    {
				    "isHandOver":1,
					"currentStep":"2016-06-14 16:23:05 => 资产交付"
				}
		}
}
*/

// status of echo vip in task
type VipStatusByTaskItem struct {
	IsHandOver  *int    // 1 is handed over, 0 is not handed over
	CurrentStep *string // current step of task
}

// status of all vips in task, vip => status info for vip
type VipStatusByTaskItems map[string]VipStatusByTaskItem

type VipStatusByTask struct {
	Status int                  // 0 is success, 1 is failed
	Msg    string               // ok when success, error msg when failed
	Data   VipStatusByTaskItems // status of task
}

// check vip status info format in task for result of api
func checkVipStatusByTaskItems(vipStatusInfoList VipStatusByTaskItems) error {
	for _, vipStatusInfo := range vipStatusInfoList {
		if vipStatusInfo.IsHandOver == nil || vipStatusInfo.CurrentStep == nil {
			return fmt.Errorf("wrong format")
		}
	}

	return nil
}

func getVipStatusByTask(urlStr string) (VipStatusByTaskItems, error) {
	// requset api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	// Note: in case status != 0, a decode error will be returned
	result := VipStatusByTask{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	// check result
	if result.Data == nil {
		return nil, fmt.Errorf("rms err: no data return")
	}
	if err = checkVipStatusByTaskItems(result.Data); err != nil {
		return nil, fmt.Errorf("rms err: %s", err.Error())
	}

	return result.Data, nil
}

/* GetVipStatusByTask - get status by rms task id
 *
 * Params:
 *      - taskId: rms task id
 *
 * Returns:
 *      - (status for all vips in task, err)
 */
func GetVipStatusByTask(taskId int64) (VipStatusByTaskItems, error) {
	// generate url
	url := fmt.Sprintf("http://rms.baidu.com/?r=interface/api&handler=getBgwListCurrentStep&list_id=%d", taskId)

	// request api for result
	return getVipStatusByTask(url)
}
