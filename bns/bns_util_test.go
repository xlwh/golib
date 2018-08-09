/* bns_util_test.go - unit test for bns_util */
/*
modification history
--------------------
2017/5/19, by Sijie Yang, create
*/
/*
DESCRIPTION
*/
package bns

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetLocalName(t *testing.T) {
	bnsClient := NewClient()

	// local local name conf
	filename := "./testdata/name_conf.data"
	err := LoadLocalNameConf(filename)
	if err != nil {
		t.Errorf("LoadLocalNameConf error: %s", err)
		return
	}

	tests := []struct {
		Name      string
		Err       error
		Instances []Instance
	}{
		{
			"s1.baidu.yf",
			nil,
			[]Instance{
				Instance{"10.1.1.1", 8080, 10},
				Instance{"10.1.1.2", 8080, 20},
			},
		},
		{
			"s2.baidu.yf",
			nil,
			[]Instance{
				Instance{"10.2.1.1", 8080, 10},
			},
		},
		{
			"s3.baidu.yf",
			errors.New("GetInstances fail"),
			nil,
		},
	}

	// run cases
	for i, tt := range tests {
		instances, err := GetInstances(bnsClient, tt.Name)
		if tt.Err == nil {
			if !reflect.DeepEqual(instances, tt.Instances) {
				t.Errorf("case %d expect %v, got %v", i, tt.Instances, instances)
			}
			continue
		}
		if err == nil {
			t.Errorf("case %d expect error", i)
		}
	}
}
