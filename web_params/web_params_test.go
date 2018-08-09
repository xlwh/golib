/* web_params_test.go - test for web_params.go  */
/*
modification history
--------------------
2014/9/23, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package web_params

import (
    "testing"
)

func TestParamsValueGet(t *testing.T) {
    params := make(map[string][]string)

    params["format"] = []string{"json", "noah"}
    
    value, err := ParamsValueGet(params, "form")
    if err == nil {
        t.Error("err in ParamsValueGet(), should no value for 'form'")
    }
    
    value, err = ParamsValueGet(params, "format")
    if value != "json" {
        t.Error("err in ParamsValueGet(), value should be 'json'")
    }    
}

func TestParamsMultiValueGet(t *testing.T) {
    params := make(map[string][]string)

    params["format"] = []string{"json", "noah"}
    
    value, err := ParamsMultiValueGet(params, "form")
    if err == nil {
        t.Error("err in ParamsMultiValueGet(), should no value for 'form'")
    }
    
    value, err = ParamsMultiValueGet(params, "format")
    if value == nil {
        t.Error("err in ParamsMultiValueGet(), value should not be nil")
        return
    }
    if len(value) != 2 || value[0] != "json" || value[1] != "noah" {
        t.Error("err in ParamsMultiValueGet(), err in value for 'format'")
    }
}
