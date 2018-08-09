/* vip_isp_test.go - test for vip_isp.go */
/*
modification history
--------------------
2016/07/20, by liuxiaowei07, create
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
func Test_GetIspByVip_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"msg":"Ok","data":{"id":"1","vip":"101.0.0.101","name":"电信"}}`)
	}))
	defer ts.Close()

	isp, err := getIspByVip(ts.URL)
	if err != nil {
		t.Errorf("err in getIspByVip(): %s", err.Error())
		return
	}

	if *isp != "电信" {
		t.Errorf("err in getIspByVip(): isp should be 电信")
	}
}

// server error
func Test_GetIspByVip_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getIspByVip(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetIspByVip_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getIspByVip(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// data nil
func Test_GetIspByVip_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"msg":"Ok"}`)
	}))
	defer ts.Close()

	_, err := getIspByVip(ts.URL)
	errInfo := "rms error: no isp info return"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// isp name nil
func Test_GetIspByVip_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"msg":"Ok","data":{}}`)
	}))
	defer ts.Close()

	_, err := getIspByVip(ts.URL)
	errInfo := "rms error: no isp info return"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}
