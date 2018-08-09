/* noah_encode_test.go - test for noah_encode.go    */
/*
modification history
--------------------
2014/9/28, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package noah_encode

import (
    "fmt"
    "strings"
    "testing"
)

type testData struct {
    a   int
    b   string
    c   int32
}

func TestEncode(t *testing.T) {
    var data testData

    data.a = 123
    data.b = "456"
    data.c = 789

    buf, err := Encode(data)
    
    if err != nil {
        errStr := fmt.Sprintf("err in Encode():%s", err.Error())
        t.Error(errStr)
        return
    } 

    str := string(buf)
    str = strings.TrimSuffix(str, "\n")
    strs := strings.Split(str, "\n")
    
    strMap := map[string]bool {
        "a:123":true,
        "b:456":true,
        "c:789":true,
    }
    
    for _, str = range strs {
        _, ok := strMap[str]
        if !ok {
            t.Error("err in Encode(): result is not expected")
            return
        }
        
        delete(strMap, str)
    }
}
