/* rs_test.go - test for rs.go */
/*
modification history
--------------------
2016/07/20, by Xiaowei Liu, create
*/
/*
DESCRIPTION
*/

package baidu_rms

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// right case
func Test_GetRsInfoByHost_1(t *testing.T) {
	result := `{"status":0,"message":[{"host":"yf-s-bfe00.yf01","id":"363407","ip":"10.38.159.43","type":"server_info"}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()

	_, err := getRsInfoByHost(ts.URL)
	if err != nil {
		t.Errorf("err in getRsInfoByHost(): %s", err.Error())
		return
	}
}

// server error
func Test_GetRsInfoByHost_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getRsInfoByHost(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetRsInfoByHost_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getRsInfoByHost(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// type of data body wrong
func Test_GetRsInfoByHost_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":1,"message":"serverid empty"}`)
	}))
	defer ts.Close()

	_, err := getRsInfoByHost(ts.URL)
	if err == nil {
		t.Error("should decode error")
	}
}

// no data info
func Test_GetRsInfoByHost_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0}`)
	}))
	defer ts.Close()

	_, err := getRsInfoByHost(ts.URL)
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
	}
}

// wrong format
func Test_GetRsInfoByHost_6(t *testing.T) {
	result := `{"status":0,"message":[{"host":"yf-s-bfe00.yf01","id":"363407","ip":"10.38.159.43"}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()

	_, err := getRsInfoByHost(ts.URL)
	errInfo := "checkRsInfoByHost(): wrong format"
	if err == nil || !strings.HasPrefix(err.Error(), errInfo) {
		t.Errorf("err should start with %s", errInfo)
	}
}

// right case
func Test_GetRsInfoByIp_1(t *testing.T) {
	result := `[{"id":"363407","idc":"YF","idc_id":"147","hostname":"yf-s-bfe00.yf01","ip_in1":"10.38.159.43"}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()

	_, err := getRsInfoByIp(ts.URL)
	if err != nil {
		t.Errorf("err in getRsInfoByIp(): %s", err.Error())
		return
	}
}

// server error
func Test_GetRsInfoByIp_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getRsInfoByIp(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetRsInfoByIp_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getRsInfoByIp(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// wrong format
func Test_GetRsInfoByIp_4(t *testing.T) {
	result := `[{"id":"363407","idc":"YF","idc_id":"147","hostname":"yf-s-bfe00.yf01"}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()

	_, err := getRsInfoByIp(ts.URL)
	errInfo := "checkRsInfoByIp(): wrong format"
	if err == nil || !strings.HasPrefix(err.Error(), errInfo) {
		t.Errorf("err should start with %s", errInfo)
	}
}
