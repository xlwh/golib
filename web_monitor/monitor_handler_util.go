/* monitor_handler_util.go - utils for web_monitor */
/*
modification history
--------------------
2014/12/16, by Sijie Yang, create
*/
/*
DESCRIPTION
    utility functions for simplifying use of the library
*/
package web_monitor 

import (
    "bytes"
    "encoding/json"
    "fmt"
    "runtime"
)

import (
    "www.baidu.com/golang-lib/delay_counter"
    "www.baidu.com/golang-lib/module_state2"
    "www.baidu.com/golang-lib/noah_encode"
)

// function prototype for getting CounterDiff/DelayOutput/StateData
type GetCounterDiffFunc func() *module_state2.CounterDiff
type GetDelayOutputFunc func() *delay_counter.DelayOutput
type GetStateDataFunc   func() *module_state2.StateData

/* CreateStateDataHandler - create monitor handler for StateData
 *
 * Params:
 *     - getter: func for getting StateData
 * 
 * Returns:
 *     - a monitor handler
 */
func CreateStateDataHandler(getter GetStateDataFunc) interface{} {
    return func(params map[string][]string) ([]byte, error) {
        var buff []byte
        var err error

        // get StateData
        state := getter()
        if state == nil {
            return nil, fmt.Errorf("GetStateDataFunc: invalid data")
        }

        // return encoded data 
        format := GetFormatParam(params)
        switch format {
            case "json":
                buff, err = json.Marshal(state)
            case "noah":
                buff = state.NoahString()
            case "noah_with_program_name":
                buff = state.NoahStringWithProgramName()
			default:
                err = fmt.Errorf("invalid format:%s", format)
        }
        return buff, err
    }
}

/* CreateDelayOutputHandler - create monitor handler for DelayRecent
 *
 * Params:
 *     - getter: func for getting DelayOutput
 * 
 * Returns:
 *     - a monitor handler
 */
func CreateDelayOutputHandler(getter GetDelayOutputFunc) interface{} {
    return func(params map[string][]string) ([]byte, error) {
        var buff []byte
        var err error

        // get DelayOutput
        delay := getter()
        if delay == nil {
            return nil, fmt.Errorf("GetDelayOutputFunc: invalid data")
        }

        // return encoded data 
        format := GetFormatParam(params)
        switch format {
            case "json":
                buff, err = delay.GetJson()
            case "noah":
                buff = delay.GetNoah()
            case "noah_with_program_name":
                buff = delay.GetNoahWithProgramName()
            default:
                err = fmt.Errorf("invalid format:%s", format)
        }
        return buff, err
    }
}

/* CreateCounterDiffHandler - create monitor handler for CounterDiff
 *
 * Params:
 *     - getter: func for getting CounterDiff
 * 
 * Returns:
 *     - a monitor handler
 */
func CreateCounterDiffHandler(getter GetCounterDiffFunc) interface{}  {
    return func(params map[string][]string) ([]byte, error) {
        var buff []byte
        var err error

        // get CounterDiff
        diff := getter()
        if diff == nil {
            return nil, fmt.Errorf("GetCounterDiffFunc: invalid data")
        }

        // return encoded data
        format := GetFormatParam(params)
        switch format {
            case "json":
                buff, err = json.Marshal(diff)
            case "noah":
                buff = diff.NoahString()
			case "noah_with_program_name":
				buff = diff.NoahStringWithProgramName()
            default:
                err = fmt.Errorf("invalid format:%s", format)
        }
        return buff, err
    }
}

/* CreateMemStatsHandler - create monitor handler for getting memory statistics
 *
 * Params:
 *     - keyPrefix: prefix of noah key, eg. <ServerName>_mem_stats
 *
 * Return:
 *     - a monitor handler
 */
func CreateMemStatsHandler(keyPrefix string) interface{} {
    return func (params map[string][]string) ([]byte, error) {
        var buff []byte
        var err error

        // get memory statistics
        var stat runtime.MemStats
        runtime.ReadMemStats(&stat)

        // return encoded data
        format := GetFormatParam(params)
        switch format {
            case "json":
                buff, err = json.Marshal(stat)
            case "noah":
                buff, err = MemStatsNoahEncode(stat, keyPrefix)
            default:
                err = fmt.Errorf("invalid format:%s", format)
        }
        return buff, err
    }
}

// get encode data of MemStats in noah format
func MemStatsNoahEncode(stat runtime.MemStats, keyPrefix string) ([]byte, error) {
    // for fields of baisc type
    buff, err := noah_encode.EncodeData(stat, keyPrefix, true)
    if err != nil {
        return nil, err
    }

    // for special fields
    prefix := keyPrefix
    if prefix != "" {
        prefix = prefix + "_"
    }
    var data bytes.Buffer

    // Note: stat.PauseNs is circular buffer of recent GC pause durations, 
    // most recent at [(NumGC+255) % 256]
    data.WriteString(fmt.Sprintf("%s%s:%d\n", prefix, "LastPauseNs", 
                     stat.PauseNs[(stat.NumGC+255) % 256]))

    buff = append(buff, data.Bytes()...)
    return buff, nil
}

// get format parameter
func GetFormatParam(params map[string][]string) string {
    format, err := ParamsValueGet(params, "format")
    if err != nil {
        format = "json" // default format is json
    }
    return format
}
