// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type successHandler struct{}

func (h successHandler) getServiceToken() (Token, error) {
	return Token{
		TokenID:   "testServiceToken",
		ExpiresAt: "2016-03-18T16:32:00.000000Z",
	}, nil
}

func (h successHandler) getUserSecretKey(userAccessKey string, serviceToken string) (AccessKey, error) {
	return AccessKey{}, nil
}

func (h successHandler) validateRequest(serviceToken string, requestArgs RequestArguments) (Token, error) {
	return Token{}, nil
}

type serviceTokenFailedHandler struct {
	successHandler
}

func (h serviceTokenFailedHandler) getServiceToken() (Token, error) {
	return Token{}, errors.New("test get service token failed.")
}

type failedHandler struct {
	successHandler
}

func (h failedHandler) getUserSecretKey(userAccessKey string, serviceToken string) (AccessKey, error) {
	return AccessKey{}, errors.New("test get user secret key failed.")
}

func (h failedHandler) validateRequest(serviceToken string, requestArgs RequestArguments) (Token, error) {
	return Token{}, errors.New("test validate request failed.")
}

func TestGetCachedServiceToken(t *testing.T) {
	iamClient := &IamClient{}
	iamClient.serviceTokenCache = "testServiceToken"

	// test if service token is not expired
	duration := time.Duration(time.Duration(1000) * time.Second)
	iamClient.serviceTokenExpireTime = time.Now().UTC().Add(duration)
	serviceToken := iamClient.getCachedServiceToken()
	assert.Equal(t, iamClient.serviceTokenCache, serviceToken)

	// test if service token is expired
	iamClient.serviceTokenExpireTime = time.Now().UTC()
	serviceToken = iamClient.getCachedServiceToken()
	assert.Equal(t, 0, len(serviceToken))
}

func TestUpdateCachedServiceToken(t *testing.T) {
	iamClient := &IamClient{}
	serviceToken := "testServiceToken"
	expiredTime := "2016-03-17T15:14:05.123456Z"
	expectExpiresAt := "2016-03-17T15:04:05Z"
	err := iamClient.updateCachedServiceToken(serviceToken, expiredTime)
	assert.Nil(t, err)
	assert.Equal(t, serviceToken, iamClient.serviceTokenCache)
	assert.Equal(t, expectExpiresAt, iamClient.serviceTokenExpireTime.Format(time.RFC3339))

	// test if parse expiredTime failed
	expiredTime = "2016-03-17T15:04:05"
	err = iamClient.updateCachedServiceToken(serviceToken, expiredTime)
	assert.NotNil(t, err)
}

func TestClientGetUserSecretKey(t *testing.T) {
	iamClient := &IamClient{
		handler: &successHandler{},
	}
	_, err := iamClient.GetUserSecretKey("testAK")
	assert.Nil(t, err)

	iamClient.handler = &failedHandler{}
	_, err = iamClient.GetUserSecretKey("testAK")
	assert.NotNil(t, err)
	assert.Equal(t, "test get user secret key failed.", err.Error())
}

func TestClientValidateRequest(t *testing.T) {
	iamClient := &IamClient{
		handler: &successHandler{},
	}
	_, err := iamClient.ValidateRequest(RequestArguments{})
	assert.Nil(t, err)

	iamClient.handler = &failedHandler{}
	_, err = iamClient.ValidateRequest(RequestArguments{})
	assert.NotNil(t, err)
	assert.Equal(t, "test validate request failed.", err.Error())
}

func TestClientGetServiceToken(t *testing.T) {
	iamClient := &IamClient{
		handler: &successHandler{},
	}
	serviceToken, err := iamClient.getServiceToken()
	assert.Nil(t, err)
	assert.Equal(t, "testServiceToken", serviceToken)
	assert.Equal(t, "2016-03-18T16:22:00Z", iamClient.serviceTokenExpireTime.Format(time.RFC3339))
	assert.NotNil(t, iamClient.serviceTokenCache)

	iamClient.handler = &serviceTokenFailedHandler{}
	serviceToken, err = iamClient.getServiceToken()
	assert.NotNil(t, err)
	assert.Equal(t, "test get service token failed.", err.Error())
}
