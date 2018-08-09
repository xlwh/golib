// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	strList := make([]string, 0)
	strList = append(strList, "test1")
	strList = append(strList, "test2")

	assert.Equal(t, true, stringInSlice("test1", strList))
	assert.Equal(t, true, stringInSlice("test2", strList))
	assert.Equal(t, false, stringInSlice("test3", strList))
}

func TestNomalizeString(t *testing.T) {
	testString1 := "test space. test alphanumeric string. 123 abc."
	expectString1 := "test%20space.%20test%20alphanumeric%20string.%20123%20abc."
	assert.Equal(t, expectString1, NomalizeString(testString1))

	testString2 := "test unreserved characters - . _ ~"
	expectString2 := "test%20unreserved%20characters%20-%20.%20_%20~"
	assert.Equal(t, expectString2, NomalizeString(testString2))

	testString3 := "test persent encode characters! @ # $ % ^ & *"
	expectString3 := "test%20persent%20encode%20characters%21%20%40%20%23%20%24%20%25%20%5E%20%26%20%2A"
	assert.Equal(t, expectString3, NomalizeString(testString3))
}

func TestNomalizePath(t *testing.T) {
	testString1 := "test/space./test/alphanumeric/string./123/abc."
	expectString1 := "test/space./test/alphanumeric/string./123/abc."
	assert.Equal(t, expectString1, nomalizePath(testString1))

	testString2 := "test/unreserved/characters/-/./_/~"
	expectString2 := "test/unreserved/characters/-/./_/~"
	assert.Equal(t, expectString2, nomalizePath(testString2))

	testString3 := "test/persent/encode/characters!/@#$%^&*"
	expectString3 := "test/persent/encode/characters%21/%40%23%24%25%5E%26%2A"
	assert.Equal(t, expectString3, nomalizePath(testString3))
}

func TestConvertMap(t *testing.T) {
	testMap := make(url.Values)
	testMap["testQuery1"] = make([]string, 0)
	testMap["testQuery2"] = make([]string, 0)

	testMap["testQuery1"] = append(testMap["testQuery1"], "123")
	testMap["testQuery2"] = append(testMap["testQuery2"], "456")

	newMap := ConvertMap(testMap)
	assert.EqualValues(t, "123", newMap["testQuery1"])
	assert.EqualValues(t, "456", newMap["testQuery2"])
}
