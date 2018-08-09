/* delay_recent.go - recent delay summary   */
/*
modification history
--------------------
2014/3/20, by Zhang Miao, create
2014/9/5,  by Zhang Miao, move from waf-server to golang-lib
2015/6/15, by Li Bingyi, move FormatOutput() from waf-server to golang-lib
*/
/*
DESCRIPTION
*/
package delay_counter

import (
    "bytes"
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

import (
	"www.baidu.com/golang-lib/module_state2"
    "www.baidu.com/golang-lib/web_params"
)

type DelayRecent struct {
    lock        sync.Mutex    
    
    interval        int         // interval of making switch
    
    currTime    time.Time       
    current     DelaySummary    // data for current minute
    
    pastTime    time.Time
    past        DelaySummary    // data for last minute

	// for noah output
    NoahKeyPrefix   string      // prefix for noah key
	ProgramName		string		// program name
}

/* for json output  */
type DelayOutput struct {
    Interval        int
    NoahKeyPrefix   string
	ProgramName		string

    CurrTime    string
    Current     DelaySummary

    PastTime    string
    Past        DelaySummary
}

/* initialize delay table
 *
 * Params:
 *      - interval: interval for move current to past
 *      - bucketSize: size of each delay bucket, e.g., 1(ms) or 2(ms)
 *      - number of bucket  
 */
func (t *DelayRecent) Init(interval int, bucketSize int, bucketNum int) {
    t.currTime = time.Now()
    // adjust time
    t.currTime = t.currTime.Truncate(time.Duration(interval) * time.Second)

    t.interval = interval
    
    // initialize DelayCounters
    t.current.Init(bucketSize, bucketNum)
    t.past.Init(bucketSize, bucketNum)
}

/* prefix is used for Noah Key generate */
func (t *DelayRecent) SetNoahKeyPrefix(prefix string) {
    t.NoahKeyPrefix = prefix
}

/* program is also used for Noah Key generate */
func (t *DelayRecent) SetProgramName(programName string) {
    t.ProgramName = programName
}

/* add one new data to the table, by providing start time and end time  */
func (t *DelayRecent) AddBySub(start time.Time, end time.Time) {
    /* get duration from start to now, in Microsecond   */
    duration := end.Sub(start).Nanoseconds() / 1000

    t.Add(duration)
}

// clear counters
func (t *DelayRecent) Clear() {
    t.current.Clear()
    t.past.Clear()
}

/* add one new data to the table.    
 * 
 * Params:
 *      - duration: delay duration, in Microsecond (10^-6)
 */
func (t *DelayRecent) Add(duration int64) {
    t.lock.Lock()    
    defer t.lock.Unlock()  

    t.trySwitch()
    t.current.Add(duration)
}

/* add one new data to the table.    
 * 
 * Params:
 *      - duration: time duration of delay (in Nanosecond)
 */
func (t *DelayRecent) AddDuration(duration time.Duration) {
    delay := int64(duration/time.Microsecond)
    t.Add(delay)
}

// check and switch DelayRecent
func (t *DelayRecent) trySwitch() {
    now := time.Now()
    if (t.currTime.Unix() / int64(t.interval)) != (now.Unix() / int64(t.interval)) {
        /* they are not in the same minute, do a switch */
        t.pastTime = t.currTime
        t.currTime = now
        
        t.past.Copy(t.current)
        t.current.Clear()       // clear t.current
    }
}

func (t *DelayRecent) get() DelayOutput {
    var retVal DelayOutput

    t.lock.Lock()
    defer t.lock.Unlock()

    t.trySwitch()
    retVal.Interval = t.interval
    retVal.CurrTime = fmt.Sprintf(t.currTime.Format("2006-01-02 15:04:05"))
    retVal.Current.Copy(t.current)
    retVal.PastTime = fmt.Sprintf(t.pastTime.Format("2006-01-02 15:04:05"))
    retVal.Past.Copy(t.past)
    
    // set noah key prefix and program name
    retVal.NoahKeyPrefix = t.NoahKeyPrefix
	retVal.ProgramName = t.ProgramName

    return retVal        
}

/* get counter from table    */
func (t *DelayRecent) Get() DelayOutput {
    retVal := t.get()    

    // calc average
    retVal.Current.CalcAvg()
    retVal.Past.CalcAvg()
    
    return retVal
}

// get data in the table, return with json string
func (t *DelayRecent) GetJson() ([]byte, error) {
    d := t.Get()
    return d.GetJson()
}

// get data in the table, return with noah string (i.e., lines of key:value)
func (t *DelayRecent) GetNoah() []byte {
    d := t.Get()
    return d.GetNoah()
}

// get data in the table, return with noah string, with program name
func (t *DelayRecent) GetNoahWithProgramName() []byte {
    d := t.Get()
    return d.GetNoahWithProgramName()
}

// format output according format value in params
func (t *DelayRecent) FormatOutput(params map[string][]string) ([]byte, error) {
    format, err := web_params.ParamsValueGet(params, "format")
    if err != nil {
        format = "json"
    }

    switch format {
    case "json", "hier_json":
        return t.GetJson()
    case "noah":
        return t.GetNoah(), nil
	case "noah_with_program_name":
		return t.GetNoahWithProgramName(), nil
    default:
        return nil, fmt.Errorf("format not support: %s", format)
    }
}

// calculate sum of DelayOutput
func (d *DelayOutput) Sum(d2 DelayOutput) error {
    if d.Interval != d2.Interval {
        return fmt.Errorf("Interval not match")
    }

    if err := d.Current.calcSum(d2.Current); err != nil {
        return err
    }
    if err := d.Past.calcSum(d2.Past); err != nil {
        return err
    }

    if d.CurrTime < d2.CurrTime {
        d.CurrTime = d2.CurrTime
    }
    if d.PastTime < d2.PastTime {
        d.PastTime = d2.PastTime
    }
    return nil
}

// get json string for DelayOutput
func (d *DelayOutput) GetJson() ([]byte, error) {
    return json.Marshal(d)
}

// generate noah key prefix
func (d *DelayOutput) noahKeyPrefixGen(key string, withProgramName bool) string{
	return module_state2.NoahKeyGen(key, d.NoahKeyPrefix, d.ProgramName, withProgramName)
}

// get noah string for DelayOutput, without program name
func (d *DelayOutput) GetNoah() []byte {
	return d.getNoah(false)
}

// get noah string for DelayOutput, with program name
func (d *DelayOutput) GetNoahWithProgramName() []byte {
	return d.getNoah(true)
}

// get noah string for DelayOutput
func (d *DelayOutput) getNoah(withProgramName bool) []byte {
    // convert to noah string
    var buf bytes.Buffer

    // current
    str := d.noahKeyPrefixGen("Current", withProgramName)
    d.Current.NoahString(&buf, str)
    
    // past
    str = d.noahKeyPrefixGen("Past", withProgramName)
    d.Past.NoahString(&buf, str)

    return buf.Bytes()
}
