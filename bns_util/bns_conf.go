/* bns_conf.go - for accessing bns config	*/
/*
modification history
--------------------
2015/9/10, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package bns_util

import (
	"encoding/json"
	"fmt"
)

import (
	"www.baidu.com/golang-lib/http_util"
)

const (
	UPDATE_TIME_OUT = 5			// timeout for invoke bns api
	MAX_CONF_LEN = 1024 * 16 	// max length of bns config (in byte)
)

type BnsResponse struct {
	RetCode		int		`json:"retCode"`
	Msg			string	`json:"msg"`
}

/*
Update conf to bns

Params:
	- clusterID: ID of cluster in noah tree
	- authKey: token for clusterID
	- conf: config

Returns:
	error
*/
func BnsConfUpdate(clusterID, authKey, conf string) error {
	// check length of config
    if len(conf) > MAX_CONF_LEN {
        return fmt.Errorf("exceed max length(%d)", len(conf))
    }

	// do post
    reqUrl := "http://bns.noah.baidu.com/webfoot/index.php"
    postData := fmt.Sprintf("r=webfoot/UpdateServiceConf&serviceName=%s&authKey=%s&serviceConf=%s", 
						clusterID, authKey, conf)

    responseStr, err := http_util.PostRR(reqUrl, UPDATE_TIME_OUT, http_util.CONTENT_FORM, []byte(postData), nil)
    if err != nil {
        return fmt.Errorf("http_util.PostRR():%s", err.Error())
	}
	
	// check response
    // format {"retCode":0,"msg":"serviceConf changed successfully"}
	var response BnsResponse
	err = json.Unmarshal(responseStr, &response)
	if err != nil {
		return fmt.Errorf("http response:no json:%s", err.Error())
	}

	if int(response.RetCode) != 0 {
		return fmt.Errorf("http response:fail:%d", response.RetCode)
	}
	
	return nil
}