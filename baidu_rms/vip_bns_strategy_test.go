/* vip_bns_strategy_test.go - test for vip_bns_strtegy.go */
/*
modification history
--------------------
2016/06/24, by liuxiaowei07, create
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
func Test_GetBnsStrategyByVip_1(t *testing.T) {
	bnsStrategy := `{"status":0,"msg":"success","data":{"bns":"group.bfe-yf.bfe.yf:main",
	    "vip":"111.13.105.50","port":80,"status":"不在流程中","bns_strategy":"0","strategy_threshold":"100"}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, bnsStrategy)
	}))
	defer ts.Close()

	_, err := getBnsStrategyByVip(ts.URL)
	if err != nil {
		t.Errorf("err in getBnsStrategyByVip(): %s", err.Error())
	}
}

// server error
func Test_GetBnsStrategyByVip_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getBnsStrategyByVip(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetBnsStrategyByVip_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getBnsStrategyByVip(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// status is 1
func Test_GetBnsStrategyByVip_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"msg":"error","data":{"data":{}}}`)
	}))
	defer ts.Close()

	_, err := getBnsStrategyByVip(ts.URL)
	if err == nil || err.Error() != "rms err: error" {
		t.Errorf("err should be rms err: error")
	}
}

// data is nil
func Test_GetBnsStrategyByVip_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"msg":"success"}`)
	}))
	defer ts.Close()

	_, err := getBnsStrategyByVip(ts.URL)
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
	}
}

// format error
func Test_GetBnsStrategyByVip_6(t *testing.T) {
	bnsStrategy := `{"status":0,"msg":"success","data":{"bns":"group.bfe-yf.bfe.yf:main",
	    "vip":"111.13.105.50","port":80,"status":"不在流程中","bns_strategy":"0"}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, bnsStrategy)
	}))
	defer ts.Close()

	_, err := getBnsStrategyByVip(ts.URL)
	if err == nil || err.Error() != "rms err: wrong format" {
		t.Errorf("err should be rms err: wrong format")
	}
}
