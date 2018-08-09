/* vip_info_test.go - test for vip_info.go */
/*
modification history
--------------------
2016/07/08, by liuxiaowei07, create
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
func Test_GetExpireByVips_1(t *testing.T) {
	expireResult := `{"status":0,"data":[{"ip":"111.13.100.249","port":"80","load_balance":"rr",
	    "bns":"group.bfe-yf.bfe.yf:main","duty_id":"0","bns_strategy":"操作单",
		"expire":"2024-07-31"}],"msg":"ok"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, expireResult)
	}))
	defer ts.Close()

	_, err := getExpireByVips(ts.URL)
	if err != nil {
		t.Errorf("err in getBnsStrategyByVip(): %s", err.Error())
	}
}

// server error
func Test_GetExpireByVips_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getExpireByVips(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetExpireByVips_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getExpireByVips(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// status is 1
func Test_GetExpireByVips_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"msg":"1231.125.115.20 不合法！","data":""}`)
	}))
	defer ts.Close()

	_, err := getExpireByVips(ts.URL)
	if err == nil {
		t.Error("should decode error")
	}
}

// data is nil
func Test_GetExpireByVips_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"msg":"success"}`)
	}))
	defer ts.Close()

	_, err := getExpireByVips(ts.URL)
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
	}
}

// format error
func Test_GetExpireByVips_6(t *testing.T) {
	expireResult := `{"status":0,"data":[{"ip":"111.13.100.249","port":"80",
	    "bns":"group.bfe-yf.bfe.yf:main","duty_id":"0","bns_strategy":"操作单",
		"expire":"2024-07-31"}],"msg":"ok"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, expireResult)
	}))
	defer ts.Close()

	_, err := getExpireByVips(ts.URL)
	if err == nil || err.Error() != "rms err: wrong format" {
		t.Errorf("err should be rms err: wrong format")
	}
}
