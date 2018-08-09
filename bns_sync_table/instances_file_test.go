/* instances_file_test.go - test for instances_file.go  */
/*
modification history
--------------------
2016/2/14, Guang Yao, create
*/
/*
DESCRIPTION
*/
package bns_sync_table

import (
	"reflect"
	"testing"
)

// test for BnsInstancesSave() and BnsInstancesLoad()
func TestInstancesFile(t *testing.T) {
	rootPath := "./test_data/"
	service := "small.bfe.yf"
	instances := []*BnsInstance{
		&BnsInstance{
			"yf-s-bfe01.yf01",
			"10.36.22.24",
			"small.BFE.yf",
			8900,
			0,
			-1,
			1,
			map[string]string{
				"cpu_type": "E5-2620-0",
			},
			"",
			"",
			"",
		},
		&BnsInstance{
			"yf-s-bfe01.yf00",
			"10.38.159.43",
			"small.BFE.yf",
			8900,
			0,
			-1,
			1,
			map[string]string{
				"cpu_type": "E5-2620-0",
			},
			"",
			"",
			"",
		},
	}

	// save to file
	err := BnsInstancesSave(rootPath, service, instances)
	if err != nil {
		t.Errorf("err in BnsInstancesSave(%s, %s):%s", rootPath, service, err.Error())
		return
	}

	// load from file
	instancesNew, err := BnsInstancesLoad(rootPath, service)
	if err != nil {
		t.Errorf("err in BnsInstancesLoad(%s, %s):%s", rootPath, service, err.Error())
		return
	}

	if !reflect.DeepEqual(instances, instancesNew) {
		t.Errorf("instancesNew != instances, instancesNew=%v", instancesNew)
	}
}
