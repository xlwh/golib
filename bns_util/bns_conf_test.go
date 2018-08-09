/* bns_conf_test.go - test for bns_conf.go  */
/*
modification history
--------------------
2015/9/11, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package bns_util

import (
	"testing"
)

// test for BnsConfUpdate()
func TestBnsConfUpdate(t *testing.T) {
	// prepare data
	clusterID := "gslb-scheduler-debug.BFE.all"
	token := "271976d29412100ae749025f7312a99a"
	conf := `{"key":"this is a test"}`
	
	// do update
	err := BnsConfUpdate(clusterID, token, conf)
	if err != nil {
		t.Errorf("err in BnsConfUpdate():%s", err.Error())
		return
	}
}
