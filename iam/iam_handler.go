// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	IAMVersion = "v3"

	// headers
	tokenHeader       = "x-subject-token"
	contentTypeHeader = "Content-Type"
	acceptHeader      = "Accept"
	authTokenHeader   = "X-Auth-Token"
	requestID         = "X-Bce-Request-Id"

	// request method
	MethodPost = "POST"
	MethodGet  = "GET"

	// response status code
	statusOk      = 200
	statusCreated = 201
)

type IamHandler interface {
	getUserSecretKey(userAccessKey string, serviceToken string) (AccessKey, error)
	getServiceToken() (Token, error)
	validateRequest(serviceToken string, requestArgs RequestArguments) (Token, error)
}

type RequestHandler struct {
	timeout       time.Duration
	iamURL        string
	iamVersion    string
	iamUser       string
	iamPassword   string
	iamDomainName string
	iamDomainId   string
}

func (handler RequestHandler) getUserSecretKey(userAccessKey string,
	serviceToken string) (AccessKey, error) {
	var accessKey AccessKey
	iamRequest, err := handler.buildGetUserSecretKeyRequest(userAccessKey, serviceToken)
	if err != nil {
		return accessKey, err
	}

	iamResponse, err := handler.doIamRequest(iamRequest)
	if err != nil {
		return accessKey, fmt.Errorf("Get user secret key failed: %v", err)
	}
	if iamResponse.StatusCode != statusOk {
		return accessKey, fmt.Errorf("Get user secret key failed. Error code: %d, Response: %v",
			iamResponse.StatusCode, iamResponse)
	}
	defer iamResponse.Body.Close()

	body, err := extractUserSecretKeyResponseBody(iamResponse)
	if err != nil {
		return accessKey, err
	}
	accessKey = body.Key

	return accessKey, nil
}

func (handler RequestHandler) getServiceToken() (Token, error) {
	var token Token
	iamRequest, err := handler.buildGetServiceTokenRequest()
	if err != nil {
		return token, err
	}

	iamResponse, err := handler.doIamRequest(iamRequest)
	if err != nil {
		return token, fmt.Errorf("Get service token failed: %v", err)
	}
	if iamResponse.StatusCode != statusCreated {
		return token, fmt.Errorf("Get service token failed. Error code: %d, Response: %v.",
			iamResponse.StatusCode, iamResponse)
	}
	defer iamResponse.Body.Close()

	iamAuthResponseBody, err := extractIamAuthResponseBody(iamResponse)
	if err != nil {
		return token, err
	}
	token = iamAuthResponseBody.IamToken
	token.TokenID = iamResponse.Header.Get(tokenHeader)
	return token, nil
}

func (handler RequestHandler) validateRequest(serviceToken string,
	requestArgs RequestArguments) (Token, error) {
	var token Token
	iamRequest, err := handler.buildValidateRequest(serviceToken, requestArgs)
	if err != nil {
		return token, err
	}

	iamResponse, err := handler.doIamRequest(iamRequest)
	if err != nil {
		return token, fmt.Errorf("Validate request failed: %v. RequestArgs: %+v.",
			err, requestArgs)
	}
	if iamResponse.StatusCode != statusOk {
		return token, fmt.Errorf("Validate request failed: %v. Error code: %d. RequestArgs: %+v. "+
			"Response: %v", err, iamResponse.StatusCode, requestArgs, iamResponse)
	}
	defer iamResponse.Body.Close()

	iamAuthResponseBody, err := extractIamAuthResponseBody(iamResponse)
	if err != nil {
		return token, err
	}
	token = iamAuthResponseBody.IamToken
	token.TokenID = iamResponse.Header.Get(tokenHeader)
	return token, nil
}

func (handler RequestHandler) doIamRequest(iamRequest *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: handler.timeout}
	return client.Do(iamRequest)
}

func (handler RequestHandler) buildValidateRequest(serviceToken string, requestArgs RequestArguments) (*http.Request, error) {
	validateRequestUrl := fmt.Sprintf("%s/%s/BCE-CRED/accesskeys", handler.iamURL, handler.iamVersion)
	validateRequestBody, err := buildValidateRequestBody(requestArgs)
	if err != nil {
		return nil, err
	}
	iamRequest, err := http.NewRequest(MethodPost, validateRequestUrl, validateRequestBody)
	if err != nil {
		return nil, fmt.Errorf("Build http request failed: %v", err)
	}
	iamRequest.Header.Add(contentTypeHeader, "application/json")
	iamRequest.Header.Add(acceptHeader, "application/json")
	iamRequest.Header.Add(authTokenHeader, serviceToken)
	iamRequest.Header.Add(requestID, requestArgs.RequestID)
	return iamRequest, nil
}

func buildValidateRequestBody(requestArgs RequestArguments) (*bytes.Buffer, error) {
	validateRequestBody := &ValidateRequestBody{}
	validateRequestBody.Auth.Authorization = requestArgs.Authorization
	validateRequestBody.Auth.SecurityToken = requestArgs.SecurityToken
	validateRequestBody.Auth.Request.Method = requestArgs.Method
	validateRequestBody.Auth.Request.URI = requestArgs.URI
	if requestArgs.QueryParams == nil {
		requestArgs.QueryParams = make(map[string]string)
	}
	validateRequestBody.Auth.Request.Params = requestArgs.QueryParams
	if requestArgs.SignHeaders == nil {
		requestArgs.SignHeaders = make(map[string]string)
	}
	validateRequestBody.Auth.Request.Headers = requestArgs.SignHeaders
	jsonBody, err := json.Marshal(validateRequestBody)
	if err != nil {
		return nil, fmt.Errorf("Encode validate request body failed: %v", err)
	}
	return bytes.NewBuffer(jsonBody), nil
}

func extractIamAuthResponseBody(response *http.Response) (*IamAuthResponseBody, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Read iam response body failed: %v.", err.Error())
	}
	iamAuthResponseBody := &IamAuthResponseBody{}
	err = json.Unmarshal(body, iamAuthResponseBody)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal response body failed: %v", err)
	}
	return iamAuthResponseBody, nil
}

func (handler RequestHandler) buildGetServiceTokenRequest() (*http.Request, error) {
	serviceTokenUrl := fmt.Sprintf("%s/%s/auth/tokens", handler.iamURL, handler.iamVersion)
	reqBody, err := handler.buildGetServiceTokenRequestBody()
	if err != nil {
		return nil, err
	}
	iamRequest, err := http.NewRequest(MethodPost, serviceTokenUrl, reqBody)
	if err != nil {
		return nil, fmt.Errorf("Create request for getting service token failed: %v", err)
	}
	iamRequest.Header.Add(contentTypeHeader, "application/json")
	iamRequest.Header.Add(acceptHeader, "application/json")
	return iamRequest, nil
}

func (handler RequestHandler) buildGetServiceTokenRequestBody() (*bytes.Buffer, error) {
	methods := make([]string, 0)
	methods = append(methods, "password")
	serviceTokenRequestBody := ServiceTokenRequestBody{}
	serviceTokenRequestBody.Auth.Identity.Methods = methods
	serviceTokenRequestBody.Auth.Identity.Password.User.Domain.Name = handler.iamDomainName
	serviceTokenRequestBody.Auth.Identity.Password.User.Name = handler.iamUser
	serviceTokenRequestBody.Auth.Identity.Password.User.Password = handler.iamPassword
	serviceTokenRequestBody.Auth.Scope.Domain.ID = handler.iamDomainId

	jsonBody, err := json.Marshal(serviceTokenRequestBody)
	if err != nil {
		return nil, fmt.Errorf("Encode service token request body failed: %v", err)
	}
	reqBody := bytes.NewBuffer(jsonBody)
	return reqBody, nil
}

func (handler RequestHandler) buildGetUserSecretKeyRequest(userAccessKey string,
	serviceToken string) (*http.Request, error) {
	userSecretKeyUrl := fmt.Sprintf("%s/%s/BCE-CRED/accesskeys/%s", handler.iamURL,
		handler.iamVersion, userAccessKey)

	iamRequest, err := http.NewRequest(MethodGet, userSecretKeyUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Create request for getting user secret key failed: %v", err)
	}

	iamRequest.Header.Add(contentTypeHeader, "application/json")
	iamRequest.Header.Add(acceptHeader, "application/json")
	iamRequest.Header.Add(authTokenHeader, serviceToken)

	return iamRequest, nil
}

func extractUserSecretKeyResponseBody(response *http.Response) (*SecretKeyResponseBody, error) {
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Read iam response body from getting"+
			" user secret key failed: %v", err)
	}

	secretKeyResponseBody := &SecretKeyResponseBody{}
	err = json.Unmarshal(respBody, secretKeyResponseBody)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal body failed in getting user secret key, Error: %v", err)
	}
	return secretKeyResponseBody, nil
}

func NewHandler(iamURL string, iamUser string, iamPassword string,
	iamDomainName string, iamDomainId string, timeoutInSecond uint) IamHandler {
	timeout := time.Duration(time.Duration(timeoutInSecond) * time.Second)
	handler := &RequestHandler{}
	handler.timeout = timeout
	handler.iamURL = iamURL
	handler.iamVersion = IAMVersion
	handler.iamDomainName = iamDomainName
	handler.iamDomainId = iamDomainId
	handler.iamUser = iamUser
	handler.iamPassword = iamPassword
	return handler
}
