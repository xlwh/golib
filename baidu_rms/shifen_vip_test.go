/* shifen_vip_test.go - test for shifen_vip.go */
/*
modification history
--------------------
2016/08/03, by wuzhenxing, create
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
func Test_GetVipsByShifen_1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"data":{"name":"jpaas-movie.e.shifen.com","service_name":"movie","rms_product":"nuomi","noah_path":"BAIDU_LBS_nuomi","rms_department":"LBS","ipInfo":[{"ip":"123.125.115.198","status":"reserve","idc":"TC","isp":"\u8054\u901a"}]},"msg":"ok","cache":"ok"}`)
	}))
	defer ts.Close()

	shifenVip, err := getVipsByShifen("jpaas-movie.e.shifen.com", ts.URL)
	if err != nil {
		t.Errorf("err in getVipsByShifen(): %s", err.Error())
		return
	}

	if len(shifenVip.IpInfos) != 1 {
		t.Errorf("err in getVipsByShiFen(): length of shifenVip.IpInfos should be 1")
	}
}

// server error
func Test_GetVipsByShifen_2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getVipsByShifen("jpaas-movie.e.shifen.com", ts.URL)
	if err == nil || err.Error() != "http_util.Read(): status code:500" {
		t.Errorf("err should be http_util.Read(): status code:500")
	}
}

// not json
func Test_GetVipsByShifen_3(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getVipsByShifen("jpaas-movie.e.shifen.com", ts.URL)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// server error
func Test_GetVipsByShifen_4(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"data":{"name":"jpaas-movie.e.shifen.com","service_name":"movie","rms_product":"nuomi","noah_path":"BAIDU_LBS_nuomi","rms_department":"LBS","ipInfo":[{"ip":"123.125.115.198","status":"reserve","idc":"TC","isp":"\u8054\u901a"}]},"msg":"ok","cache":"ok"}`)
	}))
	defer ts.Close()

	_, err := getVipsByShifen("movie.e.shifen.com", ts.URL)
	errInfo := "shifen name in shifenVip Raw [jpaas-movie.e.shifen.com] conflict with shifen [movie.e.shifen.com]"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// server error
func Test_GetVipsByShifen_5(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":0,"data":{"name":"jpaas-movie.e.shifen.com","service_name":"movie","rms_product":"nuomi","noah_path":"BAIDU_LBS_nuomi","rms_department":"LBS","ipInfo":[{"ip":"123.125.115.198","status":"reserve","idc":"TC","isp":"\u8054\u901a"},{"ip":"123.125.115.198","status":"reserve","idc":"TC","isp":"\u8054\u901a"}]},"msg":"ok","cache":"ok"}`)
	}))
	defer ts.Close()

	_, err := getVipsByShifen("jpaas-movie.e.shifen.com", ts.URL)
	errInfo := "ip info repeated [123.125.115.198]"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}
