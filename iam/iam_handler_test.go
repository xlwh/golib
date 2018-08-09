// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	iamURL        = "testURL"
	iamUser       = "testUser"
	iamPassword   = "testPassword"
	iamDomainName = "testDomainName"
	iamDomainId   = "testDomainId"
	iamTimeout    = 30
)

var handler *RequestHandler

func init() {
	handler = &RequestHandler{}
	handler.iamURL = iamURL
	handler.iamUser = iamUser
	handler.iamPassword = iamPassword
	handler.iamVersion = IAMVersion
	handler.iamDomainName = iamDomainName
	handler.iamDomainId = iamDomainId
	timeout := time.Duration(time.Duration(iamTimeout) * time.Second)
	handler.timeout = timeout
}

func TestBuildValidateRequestBody(t *testing.T) {
	requestArgs := RequestArguments{}
	requestArgs.Method = "testMethod"
	requestArgs.URI = "testUri"
	headers := make(map[string]string)
	headers["Host"] = "testHost"
	headers["XBceDate"] = "testXBceDate"
	requestArgs.SignHeaders = headers
	requestArgs.Authorization = "testAuthorization"

	buffer, err := buildValidateRequestBody(requestArgs)
	assert.Nil(t, err)
	reqBody := &ValidateRequestBody{}
	err = json.Unmarshal(buffer.Bytes(), reqBody)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	assert.Equal(t, requestArgs.Method, reqBody.Auth.Request.Method)
	assert.Equal(t, requestArgs.URI, reqBody.Auth.Request.URI)
	assert.NotNil(t, reqBody.Auth.Request.Headers)
	assert.Equal(t, "testHost", reqBody.Auth.Request.Headers["Host"])
	assert.Equal(t, "testXBceDate", reqBody.Auth.Request.Headers["XBceDate"])
	assert.Equal(t, requestArgs.Authorization, reqBody.Auth.Authorization)
	assert.NotNil(t, reqBody.Auth.Request.Params)
	assert.Equal(t, 0, len(reqBody.Auth.Request.Params))

	// test if query params is not nil
	queryParams := make(map[string]string)
	queryParams["testKey1"] = "testValue1"
	queryParams["testKey2"] = "testValue2"
	requestArgs.QueryParams = queryParams
	buffer, err = buildValidateRequestBody(requestArgs)
	assert.Nil(t, err)
	err = json.Unmarshal(buffer.Bytes(), reqBody)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	assert.Equal(t, 2, len(reqBody.Auth.Request.Params))
	assert.Equal(t, requestArgs.QueryParams["testKey1"], reqBody.Auth.Request.Params["testKey1"])
	assert.Equal(t, requestArgs.QueryParams["testKey2"], reqBody.Auth.Request.Params["testKey2"])
}

func TestBuildValidateRequest(t *testing.T) {
	serviceToken := "testServiceToken"
	requestArgs := RequestArguments{}
	requestArgs.Method = "testMethod"
	requestArgs.URI = "testUri"
	headers := make(map[string]string)
	headers["Host"] = "testHost"
	headers["XBceDate"] = "testXBceDate"
	requestArgs.SignHeaders = headers
	requestArgs.Authorization = "testAuthorization"
	request, err := handler.buildValidateRequest(serviceToken, requestArgs)
	assert.Nil(t, err)

	expectURL := fmt.Sprintf("%s/%s/BCE-CRED/accesskeys", iamURL, IAMVersion)
	assert.Equal(t, request.URL.Path, expectURL)
	assert.Equal(t, 4, len(request.Header))
	assert.Equal(t, serviceToken, request.Header.Get(authTokenHeader))
}

func TestBuildGetServiceTokenRequestBody(t *testing.T) {
	buffer, err := handler.buildGetServiceTokenRequestBody()
	assert.Nil(t, err)

	reqBody := &ServiceTokenRequestBody{}
	json.Unmarshal(buffer.Bytes(), reqBody)
	assert.Equal(t, iamUser, reqBody.Auth.Identity.Password.User.Name)
	assert.Equal(t, iamPassword, reqBody.Auth.Identity.Password.User.Password)
	assert.Equal(t, iamDomainName, reqBody.Auth.Identity.Password.User.Domain.Name)
	assert.Equal(t, iamDomainId, reqBody.Auth.Scope.Domain.ID)
	assert.Equal(t, 1, len(reqBody.Auth.Identity.Methods))
	assert.Equal(t, "password", reqBody.Auth.Identity.Methods[0])
}

func TestBuildGetServiceTokenRequest(t *testing.T) {
	request, err := handler.buildGetServiceTokenRequest()
	assert.Nil(t, err)
	expectURL := fmt.Sprintf("%s/%s/auth/tokens", iamURL, IAMVersion)
	assert.Equal(t, expectURL, request.URL.Path)
	assert.Equal(t, 2, len(request.Header))
}

func TestBuildGetUserSecretKeyRequest(t *testing.T) {
	userAccessKey := "testAK"
	serviceToken := "testServiceToken"
	request, err := handler.buildGetUserSecretKeyRequest(userAccessKey, serviceToken)
	assert.Nil(t, err)

	expectURL := fmt.Sprintf("%s/%s/BCE-CRED/accesskeys/%s", iamURL, IAMVersion, userAccessKey)
	assert.Equal(t, expectURL, request.URL.Path)
	assert.Equal(t, 3, len(request.Header))
	assert.Equal(t, serviceToken, request.Header.Get(authTokenHeader))
}

func TestDoIamRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	anyMethod := "anyMethod"
	iamRequest1, err := http.NewRequest(anyMethod, ts.URL, nil)
	if err != nil {
		t.Fatalf("Create request failed.\n")
	}
	iamResponse1, err := handler.doIamRequest(iamRequest1)
	if err != nil {
		t.Fatalf("Do request failed: %v", err)
	}
	defer iamResponse1.Body.Close()
	respBody, err := ioutil.ReadAll(iamResponse1.Body)
	if err != nil {
		t.Fatalf("Read response body failed: %v", err)
	}
	assert.Equal(t, "Hello, client\n", string(respBody))
}

func TestGetServiceTokenSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleGetServiceTokenSuccess))
	defer ts.Close()

	duration := time.Duration(time.Duration(1000) * time.Second)
	expectCachedServiceToken := "testCachedServiceToken"
	iamClient := &IamClient{
		serviceTokenCache:      expectCachedServiceToken,
		serviceTokenExpireTime: time.Now().UTC().Add(duration),
		handler: NewHandler(ts.URL, iamUser, iamPassword,
			iamDomainName, iamDomainId, iamTimeout),
	}
	// test if iam client cached service token is not expired
	serviceToken1, err1 := iamClient.getServiceToken()
	assert.Nil(t, err1)
	assert.Equal(t, expectCachedServiceToken, serviceToken1)

	// test if iam client cached service token is expired
	expectServiceToken := "testServiceToken"
	expectExpiresAt := "2016-03-18T16:22:00Z"
	iamClient.serviceTokenExpireTime = time.Now().UTC()
	serviceToken2, err2 := iamClient.getServiceToken()
	assert.Nil(t, err2)
	assert.Equal(t, expectServiceToken, serviceToken2)
	assert.Equal(t, expectExpiresAt, iamClient.serviceTokenExpireTime.Format(time.RFC3339))
}

func TestGetServiceTokenFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleFailed))
	defer ts.Close()

	iamClient := &IamClient{
		handler: NewHandler(ts.URL, iamUser, iamPassword,
			iamDomainName, iamDomainId, iamTimeout),
	}
	serviceToken, err := iamClient.getServiceToken()
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(serviceToken))
}

func TestValidateRequestSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleValidateRequestSuccess))
	defer ts.Close()

	// do not call getServiceToken in validateRequest by setting a not expired service token
	duration := time.Duration(time.Duration(1000) * time.Second)
	cachedServiceToken := "testCachedServiceToken"
	iamClient := &IamClient{
		serviceTokenCache:      cachedServiceToken,
		serviceTokenExpireTime: time.Now().UTC().Add(duration),
		handler: NewHandler(ts.URL, iamUser, iamPassword,
			iamDomainName, iamDomainId, iamTimeout),
	}

	requestArgs := RequestArguments{}
	token, err := iamClient.ValidateRequest(requestArgs)
	assert.Nil(t, err)
	assert.Equal(t, "testUserToken", token.TokenID)
	assert.Equal(t, "2016-03-18T16:32:00.000000Z", token.ExpiresAt)
	assert.Equal(t, "test_project_id", token.Project.ID)
	assert.Equal(t, "test_project_name", token.Project.Name)
	assert.Equal(t, "test_user_id", token.User.ID)
	assert.Equal(t, "test_user_name", token.User.Name)
}

func TestValidateRequestFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleFailed))
	defer ts.Close()

	// do not call getServiceToken in validateRequest by setting a not expired service token
	duration := time.Duration(time.Duration(1000) * time.Second)
	cachedServiceToken := "testCachedServiceToken"
	iamClient := &IamClient{
		serviceTokenCache:      cachedServiceToken,
		serviceTokenExpireTime: time.Now().UTC().Add(duration),
		handler: NewHandler(ts.URL, iamUser, iamPassword,
			iamDomainName, iamDomainId, iamTimeout),
	}

	requestArgs := RequestArguments{}
	token, err := iamClient.ValidateRequest(requestArgs)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(token.ExpiresAt))
	assert.Equal(t, 0, len(token.Project.ID))
	assert.Equal(t, 0, len(token.Project.Name))
	assert.Equal(t, 0, len(token.User.ID))
	assert.Equal(t, 0, len(token.User.Name))
}

func TestGetUserSecretKeySuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleGetUserSecretKeySuccess))
	defer ts.Close()

	// do not call getServiceToken in validateRequest by setting a not expired service token
	duration := time.Duration(time.Duration(1000) * time.Second)
	cachedServiceToken := "testCachedServiceToken"
	iamClient := &IamClient{
		serviceTokenCache:      cachedServiceToken,
		serviceTokenExpireTime: time.Now().UTC().Add(duration),
		handler: NewHandler(ts.URL, iamUser, iamPassword,
			iamDomainName, iamDomainId, iamTimeout),
	}

	userAccessKey := "testAK"
	accessKey, err := iamClient.GetUserSecretKey(userAccessKey)
	assert.Nil(t, err)
	assert.Equal(t, "testAccess", accessKey.Access)
	assert.Equal(t, "testCredentialID", accessKey.CredentialID)
	assert.Equal(t, "testProjectID", accessKey.ProjectID)
	assert.Equal(t, "testSecret", accessKey.Secret)
	assert.Equal(t, "testUserID", accessKey.UserID)
}

func TestGetUserSecretKeyFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleFailed))
	defer ts.Close()

	// do not call getServiceToken in validateRequest by setting a not expired service token
	duration := time.Duration(time.Duration(1000) * time.Second)
	cachedServiceToken := "testCachedServiceToken"
	iamClient := &IamClient{
		serviceTokenCache:      cachedServiceToken,
		serviceTokenExpireTime: time.Now().UTC().Add(duration),
		handler: NewHandler(ts.URL, iamUser, iamPassword,
			iamDomainName, iamDomainId, iamTimeout),
	}

	userAccessKey := "testAK"
	accessKey, err := iamClient.GetUserSecretKey(userAccessKey)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(accessKey.Access))
	assert.Equal(t, 0, len(accessKey.CredentialID))
	assert.Equal(t, 0, len(accessKey.ProjectID))
	assert.Equal(t, 0, len(accessKey.Secret))
	assert.Equal(t, 0, len(accessKey.UserID))
}

func handleGetServiceTokenSuccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(tokenHeader, "testServiceToken")
	respBody := IamAuthResponseBody{
		IamToken: Token{
			ExpiresAt: "2016-03-18T16:32:00.000000Z",
		},
	}
	jsonBody, _ := json.Marshal(respBody)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBody)
}

func handleValidateRequestSuccess(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(tokenHeader, "testUserToken")
	respBody := IamAuthResponseBody{
		IamToken: Token{
			ExpiresAt: "2016-03-18T16:32:00.000000Z",
		},
	}
	respBody.IamToken.Project.ID = "test_project_id"
	respBody.IamToken.Project.Name = "test_project_name"
	respBody.IamToken.User.ID = "test_user_id"
	respBody.IamToken.User.Name = "test_user_name"
	jsonBody, _ := json.Marshal(respBody)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

func handleGetUserSecretKeySuccess(w http.ResponseWriter, r *http.Request) {
	respBody := SecretKeyResponseBody{
		Key: AccessKey{
			Access:       "testAccess",
			CredentialID: "testCredentialID",
			ProjectID:    "testProjectID",
			Secret:       "testSecret",
			UserID:       "testUserID",
		},
	}
	jsonBody, _ := json.Marshal(respBody)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBody)
}

func handleFailed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}
