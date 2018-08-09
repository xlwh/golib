/* vip_bns_test.go - test for vip_bns.go */
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
func Test_GetBnsByVips_1(t *testing.T) {
	token := "test_token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("ROP-Authorization") == "RopAuth test_token" {
			fmt.Fprintln(w, `{"status":0,"data":{"data":{}}}`)
		}
	}))
	defer ts.Close()

	_, err := getBnsByVips(ts.URL, token)
	if err != nil {
		t.Errorf("err in getBnsByVips(): %s", err.Error())
	}
}

// wrong token
func Test_GetBnsByVips_2(t *testing.T) {
	token := "wrong_token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("ROP-Authorization") == "RopAuth wrong_token" {
			fmt.Fprintln(w, `{"status":1,"error_code":0,"msg":"access token is invalid"}`)
		}
	}))
	defer ts.Close()

	_, err := getBnsByVips(ts.URL, token)
	errInfo := "rms err: access token is invalid"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// wrong vips
// if one wrong vip in params, the return value is wrong
// if vip is legal but vip is not exist, the return value is right without the non-existent vip
func Test_GetBnsByVips_3(t *testing.T) {
	token := "input_illegal_ip"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("ROP-Authorization") == "RopAuth input_illegal_ip" {
			fmt.Fprintln(w, `{"status":1,"error_code":500,"msg":"Illegal IP addresses: 111.111.111.111111"}`)
		}
	}))
	defer ts.Close()

	_, err := getBnsByVips(ts.URL, token)
	errInfo := "rms err: Illegal IP addresses: 111.111.111.111111"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// server error
func Test_GetBnsByVips_4(t *testing.T) {
	token := "test_token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getBnsByVips(ts.URL, token)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetBnsByVips_5(t *testing.T) {
	token := "test_token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getBnsByVips(ts.URL, token)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// no data info
func Test_GetBnsByVips_6(t *testing.T) {
	token := "test_token"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0}`)
	}))
	defer ts.Close()

	_, err := getBnsByVips(ts.URL, token)
	if err == nil || err.Error() != "rms err: no data return" {
		t.Errorf("err should be rms err: no data return")
	}
}

// format error
func Test_GetBnsByVips_7(t *testing.T) {
	token := "test_token"

	result := `{"status":0,"data":{"data":{"111.202.114.38":[{"bns":"group.bfe-tc.bfe.tc:main"}]}}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()

	_, err := getBnsByVips(ts.URL, token)
	if err == nil || err.Error() != "rms err: wrong format" {
		t.Errorf("err should be rms err: wrong format")
	}
}

// right case
func Test_GetBnsByVips_8(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"data":{"data":{}}}`)
	}))
	defer ts.Close()

	_, err := GetBnsByVips(nil, "", "", "", &ts.URL)
	if err != nil {
		t.Errorf("err in getBnsByVips(): %s", err.Error())
	}
}
