// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	prefixOfBce = "x-bce-"
)

type BceSigner interface {
	Sign(signArgs SignArguments) string
}

type BceSign struct {
	headerNeedToSign []string
}

// Create the authorization
func (sign BceSign) Sign(signArgs SignArguments) string {
	signKey, signKeyInfo := buildSignKey(signArgs)
	// build canonical request
	canonicalURI := canonicalURI(signArgs.Path)
	canonicalQueryString := canonicalQueryString(signArgs.QueryParams)
	// lower the sign headers makes sure they can be signed
	signHeaders := make([]string, 0)
	for _, signHeader := range signArgs.SignHeaders {
		signHeaders = append(signHeaders, strings.ToLower(signHeader))
	}
	// sort sign headers as iam server required
	sort.Strings(signHeaders)
	canonicalHeaders := sign.canonicalHeaders(signArgs.Headers, signHeaders)
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s", strings.ToUpper(signArgs.Method),
		canonicalURI, canonicalQueryString, canonicalHeaders)

	newSign := hmac.New(sha256.New, []byte(signKey))
	newSign.Write([]byte(canonicalRequest))
	signature := hex.EncodeToString(newSign.Sum(nil))

	authorization := ""
	if len(signArgs.SignHeaders) == 0 {
		authorization = fmt.Sprintf("%s//%s", signKeyInfo, signature)
	} else {
		signHeadersString := strings.Join(signHeaders, ";")
		authorization = fmt.Sprintf("%s/%s/%s", signKeyInfo, signHeadersString, signature)
	}

	return authorization
}

func buildSignKey(signArgs SignArguments) (string, string) {
	// convert the time with iam time format
	formatTimestamp := signArgs.TimeStamps.UTC().Format(time.RFC3339)
	signKeyInfo := fmt.Sprintf("bce-auth-v1/%s/%s/%d", signArgs.UserAccessKey, formatTimestamp,
		signArgs.ExpirationInSeconds)
	// calculate the sign key with user secret key
	signKey := hmac.New(sha256.New, []byte(signArgs.UserSecretKey))
	signKey.Write([]byte(signKeyInfo))

	return hex.EncodeToString(signKey.Sum(nil)), signKeyInfo
}

// formatting the URL with signing protocol.
func canonicalURI(path string) string {
	if len(path) == 0 {
		return "/"
	} else if strings.HasPrefix(path, "/") {
		return nomalizePath(path)
	}
	return "/" + nomalizePath(path)
}

// formatting the query string with signing protocol.
func canonicalQueryString(params map[string]string) string {
	querys := make([]string, 0)
	for key, value := range params {
		if key == "authorization" {
			continue
		}
		// canonical the key and value of a query
		key = NomalizeString(key)
		value = NomalizeString(value)
		querys = append(querys, fmt.Sprintf("%s=%s", key, value))
	}
	sort.Strings(querys)
	return strings.Join(querys, "&")
}

// formatting the headers from the request with signing protocol.
func (sign BceSign) canonicalHeaders(headers map[string]string, signHeaders []string) string {
	canonicalHeaders := make([]string, 0)
	headersToSign := sign.getSignHeaders(headers, signHeaders)
	for key, value := range headersToSign {
		key = NomalizeString(key)
		value = NomalizeString(value)
		canonicalHeaders = append(canonicalHeaders, fmt.Sprintf("%s:%s", key, value))
	}
	sort.Strings(canonicalHeaders)
	return strings.Join(canonicalHeaders, "\n")
}

func (sign BceSign) getSignHeaders(headers map[string]string,
	signHeaders []string) map[string]string {
	headersToSign := make(map[string]string)
	for key, value := range headers {
		if value = strings.TrimSpace(value); len(value) == 0 {
			continue
		}
		key = strings.TrimSpace(strings.ToLower(key))
		if strings.HasPrefix(key, prefixOfBce) || stringInSlice(key, signHeaders) ||
			stringInSlice(key, sign.headerNeedToSign) {
			headersToSign[key] = value
		}
	}
	return headersToSign
}

func NewBceSigner() BceSigner {
	bceSigner := &BceSign{}
	headers := make([]string, 0)
	headers = append(headers, "host")
	headers = append(headers, "content-md5")
	headers = append(headers, "content-length")
	headers = append(headers, "content-type")
	bceSigner.headerNeedToSign = headers
	return bceSigner
}
