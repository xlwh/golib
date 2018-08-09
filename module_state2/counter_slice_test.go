/* counter_slice_test.go - test for counter_slice.go    */
/*
modification history
--------------------
2014/4/24, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package module_state2

import (
    "fmt"
    "testing"
    "time"
)

func TestCounterSliceGet(t *testing.T) {
    var cs CounterSlice
        
    diff := cs.Get()    
    if len(diff.Diff) != 0 {
        t.Error("data in diff should be zero")
    }
    
    counters := NewCounters()
    counters["test"] = 123
    
    cs.Set(counters)
    diff = cs.Get()
    if len(diff.Diff) != 0 {
        t.Error("data in diff should be zero")
    }

    time.Sleep(time.Second)
    
    counters["test"] = 223
    cs.Set(counters)
    diff = cs.Get()

    value, ok := diff.Diff["test"]
    if !ok || value != 100 {
        t.Error("the diff for test should be 100")
    }
    if diff.Duration != 1 {
        t.Error("duration should be 1")
    }
    fmt.Printf("diff=%v\n", diff)
}

func TestCounterDiff_NoahString(t *testing.T) {
    var diff CounterDiff
    
    // prepare data
    diff.LastTime = "1234"
    diff.Duration = 5678
    
    diff.Diff = NewCounters()
    diff.Diff.inc("counter", 1)
    
    // output noah string
    strOK := "counter:1\n"
    if string(diff.NoahString()) != strOK {
         t.Error("err in CounterDiff.NoahString()")
    }
}

func TestCounterDiff_NoahString_with_progName(t *testing.T) {
    var diff CounterDiff
    
    // prepare data
    diff.LastTime = "1234"
    diff.Duration = 5678
    
    diff.Diff = NewCounters()
    diff.Diff.inc("counter", 1)
    diff.ProgramName = "program"
	
    // output noah string
    strOK := "program.counter:1\n"
    if string(diff.NoahStringWithProgramName()) != strOK {
         t.Error("err in CounterDiff.NoahString()")
    }
}

func TestFormatOutput4CounterSlice(t *testing.T) {
    var err error
    var cs CounterSlice
    var diff CounterDiff

    sd := NewStateData()
    sd.SCounters.init([]string {
        "baidu",
    })

    diff = cs.Get()
    _, err = (&diff).FormatOutput(map[string][]string{"param" : []string{
        "no_format"},
    })

    if err != nil {
        t.Errorf("TestFormatOutCD_Case0(): %s", err.Error())
    }

    diff = cs.Get()
    _, err = (&diff).FormatOutput(map[string][]string{"format" : []string{
        "json"},
    })

    if err != nil {
        t.Errorf("TestFormatOutCD_Case0(): %s", err.Error())
    }

    diff = cs.Get()
    _, err = (&diff).FormatOutput(map[string][]string{"format" : []string{
        "hier_json"},
    })

    if err != nil {
        t.Errorf("TestFormatOutCD_Case0(): %s", err.Error())
    }

    diff = cs.Get()
    _, err = (&diff).FormatOutput(map[string][]string{"format" : []string{
        "noah"},
    })

    if err != nil {
        t.Errorf("TestFormatOutCD_Case0(): %s", err.Error())
    }

    diff = cs.Get()
    _, err = (&diff).FormatOutput(map[string][]string{"format" : []string{
        "no_exist"},
    })

    if err == nil {
        t.Error("TestFormatOutCD_Case0(): err should not equal nil")
    }
}
