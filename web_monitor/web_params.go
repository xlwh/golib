/* web_params.go - operation related to web_monitor params  */
/*
modification history
--------------------
2014/7/8, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package web_monitor

import (
    "errors"
)

// get one (the first) value for given key in params
func ParamsValueGet(params map[string][]string, key string) (string, error) {
    values := params[key]
    
    if values == nil || len(values) == 0 {
        return "", errors.New("key not exist")
    }
    
    return values[0], nil
}

// get values for given key in params
func ParamsMultiValueGet(params map[string][]string, key string) ([]string, error) {
    values := params[key]
    
    if values == nil || len(values) == 0 {
        return nil, errors.New("key not exist")
    }

    return values, nil
}
