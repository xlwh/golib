/* counter_slice.go - get diff of two counters  */
/*
modification history
--------------------
2014/4/24, by Zhang Miao, create
2014/11/12, by Li Bingyi, modify.
    - move codes from waf-server for periodically get counter slice
2015/6/15, by Li Bingyi, move FormatOutput() from waf-server to golang-lib
*/
/*
DESCRIPTION

Usage:
    import "www.baidu.com/golang-lib/module_state2"
    
    var counter module_state2.Counter
    var counterSlice *module_state2.CounterSlice
    var state *module_state2.State

    // usage 1: get diff once
    counterSlice.Set(counter)
    // make some update to counter here
    counterSlice.Set(counter)
    // get diff between update
    diff := counterSlice.Get()

    // usage 2: update diff periodically and get when needed
    var examCnt examCounter
    // update diff periodically
    counterSlice.Init(state, interval)
    // get diff
    diff := counterSlice.Get()
*/

package module_state2

import (
    "bytes"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

import (
    "www.baidu.com/golang-lib/web_params"
)

/* diff of two counters    */
type CounterSlice struct {
    lock            sync.Mutex
    
	lastTime        time.Time
	duration        time.Duration
    
    countersLast    Counters    //  last absolute counter
    countersDiff    Counters    //  diff in last duration
    
    noahKeyPrefix   string      //  for noah key
	programName		string		//  program name, e.g., 'go-bfe', for displaying variable in noah
}

type CounterDiff struct {
	LastTime    string      // time till
	Duration    int         // in second

    Diff        Counters
    
    NoahKeyPrefix   string  // for noah key
	ProgramName		string	// for program name
}

/* set for noah key prefix */
func (cs *CounterSlice) SetNoahKeyPrefix(prefix string) {
    cs.noahKeyPrefix = prefix
}

/* set program name	*/
func (cs *CounterSlice) SetProgramName(programName string) {
    cs.programName = programName
}

/* get noah key prefix */
func (cs *CounterSlice) GetNoahKeyPrefix() string {
    return cs.noahKeyPrefix
}

/* set to counter slice */
func (cs *CounterSlice) Set(counters Counters) {
    cs.lock.Lock()
    defer cs.lock.Unlock()

    if cs.countersLast == nil {
        // not initialized
        cs.lastTime = time.Now()
        cs.countersLast = counters.copy()
        cs.countersDiff = NewCounters()
    } else {
        now := time.Now()
        cs.duration = now.Sub(cs.lastTime)
        cs.lastTime = now
        
        cs.countersDiff = counters.diff(cs.countersLast)
        cs.countersLast = counters.copy()
    }    
}

/* get diff from counter slice   */
func (cs *CounterSlice) Get() CounterDiff {
    var retVal CounterDiff
    
    cs.lock.Lock()
    defer cs.lock.Unlock()

    if cs.countersLast == nil {
        retVal.Diff = NewCounters()
    } else {
        retVal.LastTime = cs.lastTime.Format("2006-01-02 15:04:05")
        retVal.Duration = int(cs.duration.Seconds())
        retVal.Diff = cs.countersDiff.copy()
    }
    
    retVal.NoahKeyPrefix = cs.noahKeyPrefix
	retVal.ProgramName = cs.programName
    
    return retVal
}

// get json format of counter diff
func (cs *CounterSlice) GetJson() ([]byte, error) {
    return json.Marshal(cs.Get())
}

func (cd CounterDiff) noahKeyGen(str string, withProgramName bool) string {
	return NoahKeyGen(str, cd.NoahKeyPrefix, cd.ProgramName, withProgramName)
}

// output noah string (lines of key:value) for CounterDiff
func (cd CounterDiff) NoahString() []byte {
	return cd.noahString(false)
}

// output noah string (lines of key:value) for CounterDiff, with program name
func (cd CounterDiff) NoahStringWithProgramName() []byte {
	return cd.noahString(true)
}

// output noah string (lines of key:value) for CounterDiff
func (cd CounterDiff) noahString(withProgramName bool) []byte {
	var buf bytes.Buffer

    for key, value := range cd.Diff {
        key = cd.noahKeyGen(key, withProgramName)
        str := fmt.Sprintf("%s:%d\n", key, value)
        buf.WriteString(str)
    }
        
    return buf.Bytes()
}

// format output according format value in params
func (cd *CounterDiff) FormatOutput(params map[string][]string) ([]byte, error) {
    format, err := web_params.ParamsValueGet(params, "format")
    if err != nil {
        format = "json"
    }

    switch format {
    case "json":
        return json.Marshal(cd)
    case "hier_json":
        return GetCdHierJson(cd)
    case "noah":
        return cd.NoahString(), nil
    case "noah_with_program_name":
        return cd.NoahStringWithProgramName(), nil
		default:
        return nil, fmt.Errorf("format not support: %s", format)
    }
}

// go-routine for periodically get counter slice
func (cs *CounterSlice) handleCounterSlice(s *State, interval int) {
    for {
        counter := s.GetCounters()
        cs.Set(counter)

        leftSeconds := NextInterval(time.Now(), interval)
        time.Sleep(time.Duration(leftSeconds) * time.Second)
    }
}

// init the counter diff
// Params:
//    - s: module State
//    - interval: interval to compute between two counters
// Notice: use this method only when you need to get diff between two counters periodically
func (cs *CounterSlice) Init(s *State, interval int) {
    go cs.handleCounterSlice(s, interval)
}
