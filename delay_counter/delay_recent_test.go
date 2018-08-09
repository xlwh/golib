/* delay_recent_test.go - test for delay_recent.go    */
/*
modification history
--------------------
2014/4/2, by Zhang Miao, create
2014/9/5, by Zhang Miao, move from waf-server to golang-lib
*/
/*
DESCRIPTION
*/
package delay_counter

import (
    "time"
    "testing"
)

import (
    "www.baidu.com/golang-lib/log"
)

func TestDelayRecent(t *testing.T) {
    log.Init("test", "DEBUG", "./log", true, "D", 5)

    var delayTable DelayRecent
    
    // initialize the table
    // interval=20, bucketSize=1, bucketNum=10
    delayTable.Init(20, 1, 10)
    
    start := time.Now()
    
    // try to get when table is empty
    _, err1 := delayTable.GetJson()
    if err1 != nil {
        t.Error("Error in DelayTableGet()")
    }

    // try to invoke sub() of DelayTable
    delayTable.AddBySub(start, time.Now())
    
    // try to get again
    _, err2 := delayTable.GetJson()
    if err2 != nil {
        t.Error("Error in DelayTableGet()")
    }    

    // try to invoke add() of DelayTable
    duration := time.Now().Sub(start).Nanoseconds() / 1000
    delayTable.Add(duration)
    
    // try to get again
    _, err2 = delayTable.GetJson()
    if err2 != nil {
        t.Error("Error in DelayTableGet()")
    }
    
    log.Logger.Close()
}

func TestFormatOutput(t *testing.T) {
    var delay DelayRecent
    delay.Init(20, 1, 100)

    params := map[string][]string {
        "format" : []string{"json"},
    }
    _, err := (&delay).FormatOutput(params)
    if err != nil {
        t.Errorf("FormatOutDR(): testcase 0 : %s", err.Error())
    }

    params = map[string][]string {
        "format" : []string{"noah"},
    }

    _, err = (&delay).FormatOutput(params)
    if err != nil {
        t.Errorf("FormatOutDR(): testcase 1 : %s", err.Error())
    }

    params = map[string][]string {
        "format" : []string{"no_noah"},
    }

    _, err = (&delay).FormatOutput(params)
    if err == nil {
        t.Errorf("FormatOutDR(): testcase 2 should return error!")
    }
}
