// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	signArgs := SignArguments{
		UserAccessKey:       "testAK",
		UserSecretKey:       "testSK",
		Method:              "testMethod",
		Path:                "testPath",
		ExpirationInSeconds: 1000,
	}
	signArgs.TimeStamps = time.Date(2016, time.March, 22, 14, 0, 0, 0, time.UTC)
	signArgs.Headers = map[string]string{
		"testHeader1": "testHeaderValue1",
		"testHeader2": "testHeaderValue2",
	}
	signArgs.QueryParams = map[string]string{
		"testQuery1": "testQueryValue1",
		"testQuery2": "testQueryValue2",
		"testQuery3": "testQueryValue2",
		"testQuery4": "testQueryValue4",
	}

	// test if sign headers is empty
	signArgs.SignHeaders = make([]string, 0)
	expectAuthorization1 := "bce-auth-v1/testAK/2016-03-22T14:00:00Z/1000//7a80989ed969f301c6ed62c69e1b43eaa7e8bbb0439a32bd63ffba01e4c65b20"
	sign := NewBceSigner()
	authorization1 := sign.Sign(signArgs)
	assert.Equal(t, expectAuthorization1, authorization1)

	// test if sign headers is not empty
	signHeaders := make([]string, 2)
	signHeaders[0] = "testHeader1"
	signHeaders[1] = "testHeader2"
	signArgs.SignHeaders = signHeaders
	expectAuthorization2 := "bce-auth-v1/testAK/2016-03-22T14:00:00Z/1000/testheader1;testheader2/ee85887f183f54b5b2ffebc45e2777adbb9050885417fcc900020570422c110d"
	authorization2 := sign.Sign(signArgs)
	assert.Equal(t, expectAuthorization2, authorization2)
}

func TestBuildSignKey(t *testing.T) {
	signArgs := SignArguments{
		UserAccessKey:       "testAK",
		UserSecretKey:       "testSK",
		ExpirationInSeconds: 1000,
	}
	signArgs.TimeStamps = time.Date(2016, time.March, 22, 14, 0, 0, 0, time.UTC)

	expectSignKey := "0cb24dbeee3cac010e25854608c42b9147b6eb3545794369e75fefcb1d8ca588"
	expectSignKeyInfo := "bce-auth-v1/testAK/2016-03-22T14:00:00Z/1000"
	signKey, signKeyInfo := buildSignKey(signArgs)
	assert.Equal(t, expectSignKey, signKey)
	assert.Equal(t, expectSignKeyInfo, signKeyInfo)
}

func TestCanonicalURI(t *testing.T) {
	// test if path does not have prefix of slash
	path1 := "test!Path1/test,Path2"
	expectPath := "/test%21Path1/test%2CPath2"
	canonicalPath1 := canonicalURI(path1)
	assert.Equal(t, expectPath, canonicalPath1)

	// test if path has prefix of slash
	path2 := "/test!Path1/test,Path2"
	canonicalPath2 := canonicalURI(path2)
	assert.Equal(t, expectPath, canonicalPath2)
}

func TestCanonicalQueryString(t *testing.T) {
	params := map[string]string{
		"test!Params1":  "testParamsValue1",
		"testParams2":   "testParams,Value2",
		"authorization": "testAuthorization",
	}

	expectCanonicalQueryString := "test%21Params1=testParamsValue1&testParams2=testParams%2CValue2"
	canonicalQuery := canonicalQueryString(params)
	assert.Equal(t, expectCanonicalQueryString, canonicalQuery)
}

func TestCanonicalHeaders(t *testing.T) {
	defaultHeaders := make([]string, 0)
	defaultHeaders = append(defaultHeaders, "host")
	defaultHeaders = append(defaultHeaders, "content-md5")
	defaultHeaders = append(defaultHeaders, "content-length")
	defaultHeaders = append(defaultHeaders, "content-type")

	bceSigner := &BceSign{
		headerNeedToSign: defaultHeaders,
	}
	headers := map[string]string{
		"host":         "testHost",
		"test!Header1": "testHeaderValue1",
		"testHeader2":  "testHeader,Value2",
	}

	// test if sign headers is empty
	signHeaders := make([]string, 0)
	expectHeaderString1 := "host:testHost"
	canonicalHeadersString1 := bceSigner.canonicalHeaders(headers, signHeaders)
	assert.Equal(t, expectHeaderString1, canonicalHeadersString1)

	// test if sign headers is not empty
	signHeaders = append(signHeaders, "testheader2")
	expectHeaderString2 := "host:testHost\ntestheader2:testHeader%2CValue2"
	canonicalHeadersString2 := bceSigner.canonicalHeaders(headers, signHeaders)
	assert.Equal(t, expectHeaderString2, canonicalHeadersString2)
}
