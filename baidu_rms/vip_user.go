/* vip_user.go - get user info for vip */
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
Format of message from getVipUser:
{
    "status":0,
	"msg":"yangsijie"
}
*/

type UserByVip struct {
	Status int    // status, 0 is success
	Msg    string // error msg when failed, user name when success
}

func getUserByVip(urlStr string) (*string, error) {
	// request api
	resp, err := http_util.Read(urlStr, TIME_OUT, nil)
	if err != nil {
		return nil, fmt.Errorf("http_util.Read(): %s", err.Error())
	}

	// decode the result(json format)
	result := UserByVip{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	if result.Status != 0 {
		return nil, fmt.Errorf("rms err: %s", result.Msg)
	}

	return &result.Msg, nil
}

/* GetUserByVip - get user by the vip
 *
 * Params:
 *      - vip: vip
 *
 * Returns:
 *      - (user of vip, err)
 */
func GetUserByVip(vip string) (*string, error) {
	// generate url
	url := "http://rms.baidu.com/index.php?r=interface/api&handler=getVipUser&vip=" + vip

	// request api for result
	return getUserByVip(url)
}
