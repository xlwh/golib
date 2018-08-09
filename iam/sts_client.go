// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao (xiaoyuanhao@baidu.com)

package iam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	hostHeader = "Host"
	dateHeader = "X-Bce-Date"
)

type StsClienter interface {
	AssumeRole(accountID string, withToken bool, tokenExpiredTime int) (AssumeRoleResponse, error)
}

type StsClient struct {
	stsURL          string
	accessKey       string
	secretAccessKey string
	client          *http.Client
}

/*
 * STS server is different with IAM server
 * Arguments:
 * - accountID string, the user account id.
 * - withToken bool, be true if need to return OpenStack token.
 * - tokenExpiredTime int, the expired time of token
 *
 * Return:
 * - AssumeRoleResponse, include OpenStackToken if withToken is true.
 * - error is nil if success
 */
func (c *StsClient) AssumeRole(accountID string, roleName string, withToken bool,
	tokenExpiredTime int) (AssumeRoleResponse, error) {
	urlPattern := "%s/v1/credential?assumeRole&accountId=%s&roleName=%s&durationSeconds=%d"
	if withToken {
		urlPattern = urlPattern + "&withToken=true"
	}
	url := fmt.Sprintf(urlPattern, c.stsURL, accountID, roleName, tokenExpiredTime)

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return AssumeRoleResponse{}, err
	}
	request.Header.Set(contentTypeHeader, "application/json")
	request.Header.Set(hostHeader, request.Host)
	authorization := c.getTokenAuthorization(accountID, roleName,
		tokenExpiredTime, withToken, request)
	request.Header.Set("authorization", authorization)
	response, err := c.client.Do(request)
	if err != nil {
		return AssumeRoleResponse{}, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return AssumeRoleResponse{}, err
	}
	if response.StatusCode != http.StatusOK {
		return AssumeRoleResponse{}, fmt.Errorf("send assumeRole request got statusCode: %d, "+
			"body: %v", response.StatusCode, string(responseBody))
	}

	tokenResponse := AssumeRoleResponse{}
	if err := json.Unmarshal(responseBody, &tokenResponse); err != nil {
		return AssumeRoleResponse{}, err
	}
	if tokenResponse.Token.ID != "" {
		tokenResponse.OpenStackToken = tokenResponse.Token.ID
	}
	return tokenResponse, nil
}

// Fetch token request need a authorization
func (c *StsClient) getTokenAuthorization(accountID string, roleName string, tokenExpiredTime int,
	withToken bool, request *http.Request) string {
	timestamp := time.Now()
	queryParams := map[string]string{
		"assumeRole":      "",
		"accountId":       accountID,
		"durationSeconds": strconv.Itoa(tokenExpiredTime),
		"roleName":        roleName,
	}
	if withToken {
		queryParams["withToken"] = "true"
	}
	signArgs := SignArguments{
		UserAccessKey: c.accessKey,
		UserSecretKey: c.secretAccessKey,
		Method:        "POST",
		Path:          "/v1/credential",
		SignHeaders:   []string{hostHeader, contentTypeHeader},
		QueryParams:   queryParams,
		Headers: map[string]string{
			hostHeader:        request.Header.Get(hostHeader),
			contentTypeHeader: request.Header.Get(contentTypeHeader),
		},
		TimeStamps:          timestamp,
		ExpirationInSeconds: 1800,
	}
	sign := NewBceSigner()
	return sign.Sign(signArgs)
}

func NewStsClient(url string, serviceAccessKey string, serviceSecretAccessKey string,
	timeoutInSecond uint) *StsClient {
	return &StsClient{
		stsURL:          url,
		accessKey:       serviceAccessKey,
		secretAccessKey: serviceSecretAccessKey,
		client:          &http.Client{Timeout: time.Duration(timeoutInSecond) * time.Second},
	}
}
