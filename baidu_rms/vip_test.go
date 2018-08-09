/* vip_test.go - test for vip.go */
/*
modification history
--------------------
2015/12/21, by Guang Yao, create based on bfetools/vip-monitor of zhangjiyang
2016/06/16, by Xiaowei Liu, modify: use mock for http
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
func Test_GetVipByRs_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"message":{}}`)
	}))
	defer ts.Close()

	_, err := getVipByRs(ts.URL)
	if err != nil {
		t.Errorf("err in getVipByRs(): %s", err.Error())
		return
	}
}

// server error
func Test_GetVipByRs_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getVipByRs(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetVipByRs_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getVipByRs(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// type of data body wrong
func Test_GetVipByRs_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"message":"serverid empty"}`)
	}))
	defer ts.Close()

	_, err := getVipByRs(ts.URL)
	if err == nil {
		t.Error("should decode error")
	}
}

// no data info
func Test_GetVipByRs_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0}`)
	}))
	defer ts.Close()

	_, err := getVipByRs(ts.URL)
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
	}
}

// right case
func Test_GetVipByBns_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"data":[{"test_bns":"220.181.111.23"}]}`)
	}))
	defer ts.Close()

	_, err := getVipByBns(ts.URL, "test_bns")
	if err != nil {
		t.Errorf("err in getVipByBns(): %s", err.Error())
		return
	}
}

// rms err: wrong response format
func Test_GetVipByBns_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"data":[{"wrong_bns":"220.181.111.23"}]}`)
	}))
	defer ts.Close()

	_, err := getVipByBns(ts.URL, "test_bns")
	if err == nil || err.Error() != "rms err: wrong response format" {
		t.Errorf("err should be rms err: wrong response format")
	}
}

// server error
func Test_GetVipByBns_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getVipByBns(ts.URL, "test_bns")
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetVipByBns_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getVipByBns(ts.URL, "test_bns")
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// type of data body wrong
func Test_GetVipByBns_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"data":"","msg":"没有您查询的信息！"}`)
	}))
	defer ts.Close()

	_, err := getVipByBns(ts.URL, "test_bns")
	if err == nil {
		t.Error("should error")
	}
}

// no data info
func Test_GetVipByBns_6(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0}`)
	}))
	defer ts.Close()

	_, err := getVipByBns(ts.URL, "test_bns")
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
	}
}

// test for get vip by product noah path, correct case
func TestGetVipByProductNoahPath_1(t *testing.T) {
	response := `{"status":0,"msg":"","data":{"msg":"","data":{"220.181.111.191":{"idc":"M1","isp":"\u7535\u4fe1",
				    "port":{"80":"group.bfe-yf.bfe.yf"},"domains":["tieba.com"]}}}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	vipInfos, err := getVipByProductNoahPath(ts.URL, "")
	if err != nil {
		t.Errorf("err in getVipByProductNoahPath(): %s", err.Error())
		return
	}

	if len(vipInfos) != 1 {
		t.Errorf("err in getVipByProductNoahPath(): length of vipInfos should be 1")
		return
	}

	if _, ok := vipInfos["220.181.111.191"]; !ok {
		t.Errorf("err in getVipByProductNoahPath(): vipInfos should contain 220.181.111.191")
	}
}

// test for get vip by product noah path, server error
func TestGetVipByProductNoahPath_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getVipByProductNoahPath(ts.URL, "")
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// test for get vip by product noah path, response format error
func TestGetVipByProductNoahPath_3(t *testing.T) {
	response := `{"status":0,"msg":"","data":{"msg":"","data":{"220.181.111.191":{"idc":"M1","isp":"\u7535\u4fe1",
				    "port":{"80":"group.bfe-yf.bfe.yf"},"domains":["tieba.com",]}}}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	_, err := getVipByProductNoahPath(ts.URL, "")
	errInfo := "json.Unmarshal(): invalid character ']' looking for beginning of value"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
		return
	}
}

// test for get vip by product noah path, rms return error
func TestGetVipByProductNoahPath_4(t *testing.T) {
	response := `{"status":1,"error_code":500,"msg":"get vip info failed"}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	_, err := getVipByProductNoahPath(ts.URL, "")
	if err == nil || err.Error() != "rms err: get vip info failed" {
		t.Errorf("err should be rms err: error")
		return
	}
}

// test for get vip by product noah path, rms error
func TestGetVipByProductNoahPath_5(t *testing.T) {
	response := `{"status":0,"msg":""}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	_, err := getVipByProductNoahPath(ts.URL, "")
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
		return
	}
}

// test for get vip by product noah path, testUrl not nil
func TestGetVipByProductNoahPath_6(t *testing.T) {
	response := `{"status":0,"msg":"","data":{"msg":"","data":{"220.181.111.191":{"idc":"M1","isp":"\u7535\u4fe1",
				    "port":{"80":"group.bfe-yf.bfe.yf"},"domains":["tieba.com"]}}}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()

	vipInfos, err := GetVipByProductNoahPath("", "", "", "", &ts.URL)
	if err != nil {
		t.Errorf("err in getVipByProductNoahPath(): %s", err.Error())
		return
	}

	if len(vipInfos) != 1 {
		t.Errorf("err in getVipByProductNoahPath(): length of vipInfos should be 1")
		return
	}

	if _, ok := vipInfos["220.181.111.191"]; !ok {
		t.Errorf("err in getVipByProductNoahPath(): vipInfos should contain 220.181.111.191")
	}
}
