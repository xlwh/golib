/* json_test.go - unit test file for json.go */
/*
modification history
--------------------
2015/4/9, by zhangjiyang01, Create 
*/
/*
DESCRIPTION
*/
package file_util

import (
    "io/ioutil"
    "os"
    "testing"
)

func TestDumpJson(t *testing.T) {
    jsonObj := make(map[string]string)
    jsonObj["key"] = "value"

    if err := DumpJson(jsonObj, "testdata/jsonObj", 0744); err != nil {
        t.Error(err)
    }

    jsonFile, err := os.Open("testdata/jsonObj")
    if err != nil {
        t.Error(err)
    }

    if _, err = ioutil.ReadAll(jsonFile); err != nil {
        t.Error(err)
    }

}
