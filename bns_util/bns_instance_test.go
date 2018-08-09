/* bns_instance_test.go - test for bns_instance.go    */
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
	"runtime"
    "strings"
    "testing"
)

func contains(s []string, e string) bool {
    for _, a := range s { if a == e { return true } }
    return false
}

// normal case
func Test_GetBnsInstances_case1(t *testing.T) {
	if runtime.GOOS == "windows" {
		// this test can only run in linux
		return
	}

    // test no options
    instances, err := GetBnsInstances("zhangmiao.BFE.ai", nil)
    if err != nil {
        t.Errorf("err in GetBnsInstances('zhangmiao.BFE.ai'):err=%s", err.Error())
        return
    }
    if len(instances) != 1 {
        t.Errorf("len(instances) should be 1, now it's %d", len(instances))
    }
	
	if !contains(instances, "m1-op-bfe-test09.m1") {
		t.Errorf("instances should contain m1-op-bfe-test09.m1")
	}
    
    // test options
    instances, err = GetBnsInstances("zhangmiao.BFE.ai", []string{"-i", "-a"})
    if err != nil {
        t.Errorf("err in GetBnsInstances('zhangmiao.BFE.ai', ['-i', '-a']):err=%s", err.Error())
        return
    }
    if len(instances) != 1 {
        t.Errorf("len(instances) should be 1, now it's %d", len(instances))
    }
	// we cannot ensure constant ip assignment; thus only check the prefix
    validFlag := false
    for _, instance := range instances {
	    if strings.HasPrefix(instance, "m1-op-bfe-test09.m1") {
            validFlag = true
        }
    }
	if !validFlag {
        t.Errorf("one of the instances should have prefix 'm1-op-bfe-test09.m1'")
    }
}

// normal case
func Test_GetBnsInstancesIP_case1(t *testing.T) {
	if runtime.GOOS == "windows" {
		// this test can only run in linux
		return
	}

    ips, err := GetBnsInstancesIP("zhangmiao.BFE.ai")
    if err != nil {
        t.Errorf("err in GetBnsInstancesIP('zhangmiao.BFE.ai'): %s", err.Error())
        return
    }

    if len(ips) != 1 {
        t.Errorf("len(ips) should be 1, now it's %d", len(ips))
    }

    if !contains(ips, "10.42.220.52") {
        t.Errorf("ips should contain 10.42.220.52")
    }
}
