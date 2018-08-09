/* rms_token.go - get token for api in open.rms.baidu.com */
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
Format of success message:
{
    "access_token":"d34834f1de7b47a09fcc6b31eac3a950465de237",
	"expires_in":604800,
	"uid":135907,
	"refresh_token":"7fbe16dee6c074191fe1aa69094624ba224c0c62"
}
Format of failed message:
{
    "error_code":501,
	"msg":"error message"
}
*/

type AuthTokenResult struct {
	Msg           string // the result contain msg info only when get token failed
	Access_token  string // token info will be used when access api in open.rms.baidu.com
	Expires_in    int64  // expire time, now is 7 days, the value is 7*24*60*60=604800
	Uid           int64  // uid of special token user, used by SpecialAuth for background program
	Refresh_token string // refresh_token can be used to get a new access_token when access_token expires
}

func getRmsToken(urlStr string, appKey string, secretToken string, user string) (*string, error) {
	// request rms for result
	data := fmt.Sprintf("app_key=%s&secret_token=%s&user=%s&generate_type=%s",
		appKey, secretToken, user, "special_auth")
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; param=value",
	}
	resp, err := http_util.Post(urlStr, TIME_OUT, http_util.CONTENT_FORM, []byte(data), header)
	if err != nil {
		return nil, fmt.Errorf("http_util.Post(): %s", err.Error())
	}

	// decode the result(json format)
	result := AuthTokenResult{}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(): %s", err.Error())
	}

	if result.Access_token == "" {
		return nil, fmt.Errorf("rms err: %s", result.Msg)
	}
	return &result.Access_token, nil
}

/* GetRmsToken - get token by app_key, secret_token, user
 *
 * Params:
 *      - app_key: app_key applied from open.rms.baidu.com
 *      - secret_token: secret_token applied from open.rms.baidu.com
 *      - user: user of special_token applied from open.rms.baidu.com
 *
 * Returns:
 *      - (token, err)
 */
func GetRmsToken(appKey string, secretToken string, user string) (*string, error) {
	// generate url
	url := "http://api.rms.baidu.com/auth/accessToken"

	// request api for result
	return getRmsToken(url, appKey, secretToken, user)
}
