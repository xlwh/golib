/* rms_task_test.go - test for rms_task.go */
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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// right case
func Test_GetStatusByTask_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"msg":"ok","data":{}}`)
	}))
	defer ts.Close()

	_, err := getVipStatusByTask(ts.URL)
	if err != nil {
		t.Errorf("err in getVipStatusByTask(): %s", err.Error())
		return
	}
}

// server error
func Test_GetStatusByTask_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getVipStatusByTask(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetStatusByTask_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getVipStatusByTask(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// type of data body wrong
func Test_GetStatusByTask_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"msg":"List not exist: list_id=5","data":[]}`)
	}))
	defer ts.Close()

	_, err := getVipStatusByTask(ts.URL)
	if err == nil {
		t.Error("should decode error")
	}
}

// no data info
func Test_GetStatusByTask_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"msg":"ok"}`)
	}))
	defer ts.Close()

	_, err := getVipStatusByTask(ts.URL)
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
	}
}

// format error
func Test_GetStatusByTask_6(t *testing.T) {
	result := `{"status":0,"msg":"ok","data":{"112.80.248.52":{"isHandOver":1}}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()

	_, err := getVipStatusByTask(ts.URL)
	if err == nil || err.Error() != "rms err: wrong format" {
		t.Errorf("err should be rms err: wrong format")
	}
}
