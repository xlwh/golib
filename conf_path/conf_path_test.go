/* conf_path_test.go - test for conf_path.go    */
/*
modification history
--------------------
2014/8/21, by Zhang Miao, create
2014/9/29, by Zhang Miao, move from go-bfe to golang-lib
*/
/*
DESCRIPTION
*/
package conf_path

import (
    "testing"
)

func TestConfPathProc(t *testing.T) {
    confRoot := "/home/work/go-bfe/conf"
    
    // Case1: confPath is absolution path
    confPath := "/home/work/waf/waf.conf"
    confPath = ConfPathProc(confPath, confRoot)
	
	if confPath != "/home/work/waf/waf.conf" {
	    t.Errorf("err in ConfPathProc() for absolute path")
	}
	
    // Case2: confPath is relative path
    confPath = "waf/waf.conf"
    confPath = ConfPathProc(confPath, confRoot)
	
	if confPath != "/home/work/go-bfe/conf/waf/waf.conf" {
	    t.Errorf("err in ConfPathProc() for relative path")
	}	
}

