// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"fmt"
	"time"
)

const (
	// definition of time
	serviceTokenCachedTimeout = 600
)

type IamClienter interface {
	GetUserSecretKey(userAccessKey string) (AccessKey, error)
	ValidateRequest(requestArgs RequestArguments) (Token, error)
}

type IamClient struct {
	serviceTokenCache      string
	serviceTokenExpireTime time.Time
	handler                IamHandler
}

// get user secret access key by user access key
func (iamClient *IamClient) GetUserSecretKey(userAccessKey string) (AccessKey, error) {
	var accessKey AccessKey
	serviceToken, err := iamClient.getServiceToken()
	if err != nil {
		return accessKey, err
	}
	return iamClient.handler.getUserSecretKey(userAccessKey, serviceToken)
}

func (iamClient *IamClient) getServiceToken() (string, error) {
	// get cached service token first
	if cachedServiceToken := iamClient.getCachedServiceToken(); len(cachedServiceToken) > 0 {
		return cachedServiceToken, nil
	}

	// get a new service token if cached service token is expired
	token, err := iamClient.handler.getServiceToken()
	if err != nil {
		return "", err
	}
	err = iamClient.updateCachedServiceToken(token.TokenID, token.ExpiresAt)
	if err != nil {
		return "", err
	}
	return token.TokenID, nil
}

// validate user request by signature
func (iamClient *IamClient) ValidateRequest(requestArgs RequestArguments) (Token, error) {
	var token Token
	serviceToken, err := iamClient.getServiceToken()
	if err != nil {
		return token, err
	}
	return iamClient.handler.validateRequest(serviceToken, requestArgs)
}

func (iamClient *IamClient) getCachedServiceToken() string {
	if time.Now().UTC().After(iamClient.serviceTokenExpireTime) {
		iamClient.serviceTokenCache = ""
	}
	return iamClient.serviceTokenCache
}

func (iamClient *IamClient) updateCachedServiceToken(serviceToken string,
	expireTimeStr string) error {
	iamClient.serviceTokenCache = serviceToken
	expireTime, err := time.Parse(time.RFC3339Nano, expireTimeStr)
	if err != nil {
		return fmt.Errorf("Parse token expire time failed: %v", err)
	}
	duration := time.Duration(time.Duration(serviceTokenCachedTimeout) * time.Second)
	// iam expired time subtract service token timeout is the real expired time
	iamClient.serviceTokenExpireTime = expireTime.Add(-duration)
	return nil
}

func NewClient(iamURL string, iamUser string, iamPassword string,
	iamDomainName string, iamDomainId string, iamTimeout uint) IamClienter {
	iamClient := &IamClient{}
	iamClient.handler = NewHandler(iamURL, iamUser, iamPassword,
		iamDomainName, iamDomainId, iamTimeout)

	return iamClient
}
