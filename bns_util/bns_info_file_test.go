/* bns_info_file_test.go - test for bns_info_file.go  */
/*
modification history
--------------------
2015/7/29, by Zhang Miao, create
2015/8/3, by Zhang Miao, move to golang-lib
*/
/*
DESCRIPTION
*/
package bns_util

import (
	"reflect"
	"testing"
)

// test for BnsInfoSave() and BnsInfoLoad()
func TestBnsInfoFile(t *testing.T) {
	rootPath := "./"
	service := "gslb-scheduler.BFE.all"
	instances := []string{"cq01-gslb-sch-1.cq01", "tc-gslb-sch-1.tc", "yf-gslb-sch-1.yf01"}
	
	// save to file
	err := BnsInfoSave(rootPath, service, instances)
	if err != nil {
		t.Errorf("err in BnsInfoSave(%s, %s):%s", rootPath, service, err.Error())
		return
	}

	// load from file
	instancesNew, err := BnsInfoLoad(rootPath, service)
	if err != nil {
		t.Errorf("err in BnsInfoLoad(%s, %s):%s", rootPath, service, err.Error())
		return
	}
	
	if !reflect.DeepEqual(instances, instancesNew) {
		t.Errorf("instancesNew != instances, instancesNew=%s", instancesNew)
	}
}
