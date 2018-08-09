/* delay_summary_test.go - test for delay_summary.go    */
/*
modification history
--------------------
2014/9/9, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package delay_counter

import (
    "testing"
)

import (
    "www.baidu.com/golang-lib/log"
)

func TestDelaySummary(t *testing.T) {
    log.Init("test", "DEBUG", "./log", true, "D", 5)

    var counter DelaySummary

    BUCKET_SIZE := 3
    BUCKET_NUM := 10

    // test Init()
    counter.Init(BUCKET_SIZE, BUCKET_NUM)

    if len(counter.Counters) != 11 {
        t.Error("len(counter.Counters) should be 11")
    }
    
    // test Add()
    counter.Add(4000)
    counter.Add(2000)    

    if counter.Counters[0] != 1 {
        t.Error("len(counter.Counters[0]) should be 1")
    }

    if counter.Counters[1] != 1 {
        t.Error("len(counter.Counters[1]) should be 1")
    }

    // test CalcAvg()
    counter.CalcAvg()

    if counter.Ave != 3000 {
        t.Error("counter.Ave should be 3000")
    }
    
    // test Copy()
    var counterCopy DelaySummary
    counterCopy.Copy(counter)
    
    if counterCopy.BucketSize != BUCKET_SIZE || 
            counterCopy.BucketNum != BUCKET_NUM ||
            counterCopy.Count != 2 ||
            counterCopy.Sum != 6000 ||
            counterCopy.Ave != 3000 ||
            len(counterCopy.Counters) != (BUCKET_NUM + 1) ||
            counterCopy.Counters[0] != 1 ||
            counterCopy.Counters[1] != 1 {
        t.Error("error in counter.Copy()")
    }

    // test calcSum()
    var counter2 DelaySummary
    counter2.Init(BUCKET_SIZE*2, BUCKET_NUM)
    if err := counter.calcSum(counter2); err == nil {
        t.Error("should return error in counter.calcSum()")
    }

    counter2.Init(BUCKET_SIZE, BUCKET_NUM*2)
    if err := counter.calcSum(counter2); err == nil {
        t.Error("should return error in counter.calcSum()")
    }

    counter2.Init(BUCKET_SIZE, BUCKET_NUM)
    counter2.Add(4000)
    counter2.Add(8000)
    if err := counter.calcSum(counter2); err != nil {
        t.Error("error in counter.calcSum()")
    }
    if counter.Counters[0] != 1 {
        t.Error("len(counter.Counters[0]) should be 1")
    }
    if counter.Counters[1] != 2 {
        t.Error("len(counter.Counters[1]) should be 2")
    }
    if counter.Counters[2] != 1 {
        t.Error("len(counter.Counters[2]) should be 1")
    }
    if counter.Ave != 4500 {
        t.Error("counter.Ave should be 4500")
    }

    // test Clear()
    counter.Clear()
    if counter.Count != 0 || counter.Sum != 0 || counter.Ave != 0 {
        t.Error("error in counter.Clear()")
    }
    
    for i := 0; i < BUCKET_NUM + 1; i ++ {
        if counter.Counters[i] != 0 {
            t.Error("error in counter.Clear()")
        }
    }
    
    log.Logger.Close()
}
