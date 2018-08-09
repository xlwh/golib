/* vip_rs_test.go - test for vip_rs.go */
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
func Test_GetRsByVip_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"message":{"vip":"180.149.131.103","status":"1","type":"bgw","idc":"YF","port_info_list":[]}}`)
	}))
	defer ts.Close()

	_, err := getRsByVip(ts.URL)
	if err != nil {
		t.Errorf("err in getRsByVip(): %s", err.Error())
		return
	}
}

// server error
func Test_GetRsByVip_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getRsByVip(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetRsByVip_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getRsByVip(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// type of message body wrong
func Test_GetRsByVip_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"message":"List not exist: list_id=5"}`)
	}))
	defer ts.Close()

	_, err := getRsByVip(ts.URL)
	if err == nil {
		t.Error("should meet error")
	}
}

// no data info
func Test_GetRsByVip_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0}`)
	}))
	defer ts.Close()

	_, err := getRsByVip(ts.URL)
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
	}
}
