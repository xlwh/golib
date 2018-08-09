/* txt_load_test.go - unit test fo txt_load.go */
/*
modification history
--------------------
2016/12/21, by Zhang Jiyang create
*/
/*
 */
package iptable

import (
	"errors"
	"net"
	"strconv"
	"testing"
)

var notIntErr = errors.New("Type Error")

func testCheckIntFunc(value string) error {
	_, err := strconv.Atoi(value)
	if err != nil {
		return notIntErr
	}

	return nil
}

func TestLoadDictWrongCase(t *testing.T) {
	type testCase struct {
		name       string
		path       string
		curVersion string

		setCheckFunc bool

		expectErr     bool
		err           error
		expectVersion string
	}

	var testCases = []testCase{
		{
			name:          "normalCase1",
			path:          "testdata/ip_dict.data",
			expectErr:     false,
			expectVersion: "100",
		},
		{
			name:          "normalCase2",
			path:          "testdata/ip_dict.data",
			setCheckFunc:  true,
			expectErr:     false,
			expectVersion: "100",
		},
		{
			name:      "noSuchFileErr",
			path:      "testdata/no_exists.data",
			expectErr: true,
		},
		{
			name:          "withoutMetaCase",
			path:          "testdata/no_meta.data",
			expectErr:     false,
			expectVersion: "",
		},
		{
			name:          "noNeedUpdate",
			path:          "testdata/ip_dict.data",
			curVersion:    "100",
			expectErr:     true,
			err:           ERR_NO_NEED_UPDATE,
			expectVersion: "100",
		},
		{
			name:      "wrongIPDict",
			path:      "testdata/wrong_ip_dict.data",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		f := NewTxtFileLoader(tc.path)
		if tc.setCheckFunc {
			f.SetCheckValFunc(testCheckIntFunc)
		}

		_, version, err := f.CheckAndLoad(tc.curVersion)
		if tc.expectErr {
			if err == nil {
				t.Errorf("should load failed %s", tc.name)
			}
		} else {
			if err != nil {
				t.Errorf("should load success %s %s", tc.name, err.Error())
			}
		}

		// check err
		if tc.err != nil {
			if tc.err != err {
				t.Errorf("%s expect err %s, while %s", tc.name, tc.err, err)
			}
		}

		// check version
		if err == nil && tc.expectVersion != version {
			t.Errorf("%s expectVersion:version not same %s:%s", tc.name, tc.expectVersion, version)
		}
	}
}

func TestLoadAndSearchNormalCase(t *testing.T) {
	type fTestCase struct {
		name         string
		ip           string
		expectStatus bool
		expectVal    string
	}
	var fTestCases = []fTestCase{
		{
			name:         "normalCase1",
			ip:           "1.1.1.1",
			expectStatus: true,
			expectVal:    "val1",
		},
		{
			name:         "normalCase2",
			ip:           "2.2.2.2",
			expectStatus: true,
			expectVal:    "val1",
		},
		{
			name:         "normalCase3",
			ip:           "3.3.3.3",
			expectStatus: true,
			expectVal:    "val2",
		},
		{
			name:         "normalCase4",
			ip:           "3.3.3.4",
			expectStatus: true,
			expectVal:    "val2",
		},
		{
			name:         "wrongCase1",
			ip:           "6.6.6.6",
			expectStatus: false,
		},
	}
	f := NewTxtFileLoader("testdata/ip_dict.data")
	ipdict, version, err := f.CheckAndLoad("")
	if err != nil {
		t.Error(err)
	}

	expectVersion := "100"
	if version != expectVersion {
		t.Errorf("version should be %s, while %s", expectVersion, version)
		return
	}

	for _, testCase := range fTestCases {
		ip := testCase.ip
		expectStatus := testCase.expectStatus
		expectVal := testCase.expectVal

		val, ok := ipdict.Search(net.ParseIP(ip))
		if ok != expectStatus || val != expectVal {
			t.Errorf("val should be %s, while %s", expectVal, val)
		}
	}
}
