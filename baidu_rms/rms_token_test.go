/* rms_token_test.go - test for rms_token.go */
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
func Test_GetRmsToken_1(t *testing.T) {
	appKey := "111"
	secretToken := "test_token"
	user := "test_user"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"access_token":"token","expires_in":604800,"uid":135907,"refresh_token":"refresh_token"}`)
	}))
	defer ts.Close()

	_, err := getRmsToken(ts.URL, appKey, secretToken, user)
	if err != nil {
		t.Errorf("err in getRmsToken(): %s", err.Error())
	}
}

// wrong user
func Test_GetRmsToken_2(t *testing.T) {
	appKey := "111"
	secretToken := "token"
	user := "wrong_user"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PostFormValue("user") == "wrong_user" {
			fmt.Fprintln(w, `{"error_code":0,"msg":"not allowed user wrong_user"}`)
		}
	}))
	defer ts.Close()

	_, err := getRmsToken(ts.URL, appKey, secretToken, user)
	errInfo := "rms err: not allowed user wrong_user"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// wrong token
func Test_GetRmsToken_3(t *testing.T) {
	appKey := "111"
	secretToken := "wrong_token"
	user := "test_user"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PostFormValue("secret_token") == "wrong_token" {
			fmt.Fprintln(w, `{"error_code":400,"msg":"The app credentials are invalid"}`)
		}
	}))
	defer ts.Close()

	_, err := getRmsToken(ts.URL, appKey, secretToken, user)
	errInfo := "rms err: The app credentials are invalid"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// wrong key
func Test_GetRmsToken_4(t *testing.T) {
	appKey := "wrong_key"
	secretToken := "test_token"
	user := "test_user"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.PostFormValue("app_key") == "wrong_key" {
			fmt.Fprintln(w, `{"error_code":400,"msg":"The app credentials are invalid"}`)
		}
	}))
	defer ts.Close()

	_, err := getRmsToken(ts.URL, appKey, secretToken, user)
	errInfo := "rms err: The app credentials are invalid"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}

// server error
func Test_GetRmsToken_5(t *testing.T) {
	appKey := "111"
	secretToken := "test_token"
	user := "test_user"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}))
	defer ts.Close()

	_, err := getRmsToken(ts.URL, appKey, secretToken, user)
	if err == nil || err.Error() != "http_util.Post(): status code:500" {
		t.Errorf("err should be http_util.Post(): status code:500")
	}
}

// not json
func Test_GetRmsToken_6(t *testing.T) {
	appKey := "111"
	secretToken := "test_token"
	user := "test_user"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not_json`)
	}))
	defer ts.Close()

	_, err := getRmsToken(ts.URL, appKey, secretToken, user)
	errInfo := "json.Unmarshal(): invalid character 'o' in literal null (expecting 'u')"
	if err == nil || err.Error() != errInfo {
		t.Errorf("err should be %s", errInfo)
	}
}
