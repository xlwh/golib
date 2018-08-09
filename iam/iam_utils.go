// Copyright 2016 Baidu Inc. All Rights Reserved.
// Authors: Xiao Yuanhao(xiaoyuanhao@baidu.com)

package iam

import (
	"fmt"
	"strings"
)

var (
	nomalizeMap map[byte]string
)

// initialize the nomalizeMap
// the unreserved characters include:
//  - 'a' to 'z'
//  - 'A' to 'Z'
//  - '0' to '9'
//  - '-', '.', '_', '~'
func init() {
	nomalizeMap = make(map[byte]string)

	// persent encode characters
	for i := 0; i < 256; i++ {
		nomalizeMap[byte(i)] = fmt.Sprintf("%%%02X", i)
	}
	// unreserved characters
	for i := 'a'; i <= 'z'; i++ {
		nomalizeMap[byte(i)] = string(i)
	}
	for i := 'A'; i <= 'Z'; i++ {
		nomalizeMap[byte(i)] = string(i)
	}
	for i := '0'; i <= '9'; i++ {
		nomalizeMap[byte(i)] = string(i)
	}
	nomalizeMap['-'] = string('-')
	nomalizeMap['.'] = string('.')
	nomalizeMap['_'] = string('_')
	nomalizeMap['~'] = string('~')
}

func nomalizePath(oldString string) string {
	return strings.Replace(NomalizeString(oldString), "%2F", "/", -1)
}

// nomalize a string for use in BCE web service APIs
func NomalizeString(oldString string) string {
	newString := ""
	for _, v := range []byte(oldString) {
		newString = newString + nomalizeMap[v]
	}
	return newString
}

func stringInSlice(a string, list []string) bool {
	for _, v := range list {
		if a == v {
			return true
		}
	}
	return false
}

// convert the go query and header format to iam query and header format
//  - go format is map[string][]string
//  - iam fomat is map[string]string
func ConvertMap(oldMap map[string][]string) map[string]string {
	newMap := make(map[string]string)
	for key, value := range oldMap {
		newMap[key] = value[0]
	}
	return newMap
}
