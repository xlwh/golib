/* vip_product_test.go - test for vip_product.go */
/*
modification history
--------------------
2016/06/13, by liuxiaowei07, create
*/
/*
DESCRIPTION
*/

package baidu_noah

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// right case
func Test_GetProductByVips_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"success":true,"message":"ok","data":[]}`)
	}))
	defer ts.Close()

	_, err := getProductByVips(ts.URL)
	if err != nil {
		t.Errorf("err in getProductByVips(): %s", err.Error())
		return
	}
}

// server error
func Test_GetProductByVips_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getProductByVips(ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetProductByVips_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getProductByVips(ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// noah error
func Test_GetProductByVips_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"success":false,"message":"没有结果!"}`)
	}))
	defer ts.Close()

	_, err := getProductByVips(ts.URL)
	if err == nil || err.Error() != "noah err: 没有结果!" {
		t.Errorf("err should be noah err: 没有结果!")
	}
}

// no data info
func Test_GetProductByVips_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"success":true,"message":"ok"}`)
	}))
	defer ts.Close()

	_, err := getProductByVips(ts.URL)
	if err == nil || err.Error() != "noah err: no data return" {
		t.Errorf("err should be noah err: no data return")
	}
}

// format error
func Test_GetProductByVips_6(t *testing.T) {
	result := `{"success":true,"message":"ok","data":[{"productId":"200001313","id":"373365",
	    "ip":"103.235.46.141","product":"BAIDU_OP_BFE","type":"bgw","username":"liuxiaowei07"}]}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()

	_, err := getProductByVips(ts.URL)
	if err == nil || err.Error() != "noah err: wrong format" {
		t.Errorf("err should be noah err: wrong format")
	}
}
