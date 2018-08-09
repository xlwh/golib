/* vip_domain_test.go - test for vip_domain.go */
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
func Test_GetDnsByVip_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[]`)
	}))
	defer ts.Close()

	_, err := getDnsByVip(ts.URL)
	if err != nil {
		t.Errorf("err in getDnsByVip(): %s", err.Error())
		return
	}
}

// server error
func Test_GetDnsByVip_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getDnsByVip(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetDnsByVip_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getDnsByVip(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// format error
func Test_GetDnsByVip_4(t *testing.T) {
	result := `[{"dns":"cb.e.shifen.com","ip":"58.217.200.81","type":"bgw","attr":"outside",
	    "hostname":"hz01-s-bfe00.hz01","ip_in1":"10.212.71.36","ip_in2":null,"ip_department":"OP",
		"ip_product":"BFE","department":"OP","product":"BFE"}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()

	_, err := getDnsByVip(ts.URL)
	if err == nil || err.Error() != "rms err: wrong format" {
		t.Errorf("err should be rms err: wrong format")
	}
}
