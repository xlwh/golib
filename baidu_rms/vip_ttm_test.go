/* vip_ttm_test.go - test for vip_ttm.go */
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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// right case
func Test_GetTtmByVip_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"data":{"error_code":0,"error_msg":"success","ttm_vip_valid":1}}`)
	}))
	defer ts.Close()

	_, err := getTtmByVip(ts.URL)
	if err != nil {
		t.Errorf("err in getTtmByVip(): %s", err.Error())
	}
}

// wrong vip or port
func Test_GetTtmByVip_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"msg":"找不到对应的VIP"}`)
	}))
	defer ts.Close()

	_, err := getTtmByVip(ts.URL)
	errInfo := "rms err: 找不到对应的VIP"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// hk vip
func Test_GetTtmByVip_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"msg":"配置[13066 - ttm_query_url]以下不存在: network/auto_bgw/ttm_query_url/hongkong(network/auto_bgw/ttm_query_url/hongkong"}`)
	}))
	defer ts.Close()

	_, err := getTtmByVip(ts.URL)
	errInfo := "rms err: 配置[13066 - ttm_query_url]以下不存在: network/auto_bgw/ttm_query_url/hongkong(network/auto_bgw/ttm_query_url/hongkong"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// Can not found virtual server
func Test_GetTtmByVip_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"data":{"error_code":2,"error_msg":"Can not found virtual server","ttm_vip_valid":0}}`)
	}))
	defer ts.Close()

	_, err := getTtmByVip(ts.URL)
	errInfo := "rms err: Can not found virtual server"
	if err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
		return
	}
}

// server error
func Test_GetTtmByVip_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getTtmByVip(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetTtmByVip_6(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getTtmByVip(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}
