/* delay_summary.go - for summarizing time delay    */
/*
modification history
--------------------
2014/3/20, by Zhang Miao, create
2014/9/5,  by Zhang Miao, move from waf-server to golang-lib
*/
/*
DESCRIPTION
*/
package delay_counter

import (
    "bytes"
    "fmt"
)

// for holding data in recent several seconds
type DelaySummary struct {
    BucketSize  int     // size of each delay bucket, e.g., 1(ms) or 2(ms)
    BucketNum   int     // number of bucket    
    
    Count       int64   // total number of samples
    Sum         int64   // in Microsecond
    Ave         int64   // in Microsecond
    
    Counters    []int64 // counters for each bucket
                        // for bucketSize == 1ms, BucketNum == 5
                        // Counters are for 0-1, 1-2, 2-3, 3-4, 4-5, >5
}

// initialize DelaySummary
func (dc *DelaySummary) Init(bucketSize int, bucketNum int) {
    dc.BucketSize = bucketSize
    dc.BucketNum = bucketNum
    dc.Counters = make([]int64, bucketNum + 1)
}

// calculate average for DelaySummary
func (dc *DelaySummary) CalcAvg() {
    if dc.Count != 0 {
        dc.Ave = dc.Sum / dc.Count
    }
}

// clear counters
func (dc *DelaySummary) Clear() {
    dc.Count = 0
    dc.Sum = 0
    dc.Ave = 0
    
    for i := 0; i <= dc.BucketNum; i ++ {
        dc.Counters[i] = 0
    }
}

// add one new data to DelaySummary    
//
// duration - in Microsecond (10^-6)
func (dc *DelaySummary) Add(duration int64) {
    if duration < 0 {
        // this will lead to panic, so add protection here
        // this should not happen
        return
    }
    
    dc.Count += 1
    dc.Sum += duration

    // calc slot for duration
    slot := duration / int64(dc.BucketSize * 1000)

    if int(slot) < dc.BucketNum {
        dc.Counters[slot] += 1
    } else {
        dc.Counters[dc.BucketNum] += 1
    }    
}

// make a copy of src DelaySummary
func (dc *DelaySummary) Copy(src DelaySummary) {
    dc.BucketSize = src.BucketSize
    dc.BucketNum = src.BucketNum
    
    dc.Count = src.Count
    dc.Sum = src.Sum
    dc.Ave = src.Ave

    if dc.Counters == nil || len(dc.Counters) != len(src.Counters) {
        dc.Counters = make([]int64, dc.BucketNum + 1)
    }

    for i := 0; i <= dc.BucketNum; i ++ {
        dc.Counters[i] = src.Counters[i]
    }    
}

// calculate sum of DelaySummay
func (dc *DelaySummary) calcSum(dc2 DelaySummary) error {
    if dc.BucketSize != dc2.BucketSize || dc.BucketNum != dc2.BucketNum {
        return fmt.Errorf("bucket size or num not match")
    }

    dc.Count += dc2.Count
    dc.Sum += dc2.Sum
    dc.CalcAvg()
    for i := 0; i <= dc.BucketNum; i++ {
        dc.Counters[i] += dc2.Counters[i]
    }

    return nil
}

// return noah string (i.e., lines of key:value) for DelaySummary
// 
// Params:
//      - buf: buf to write string
//      - prefix: prefix add to key, e.g., prefix='Past', key='Sum', output='Past_Sum'
func (dc *DelaySummary) NoahString(buf *bytes.Buffer, prefix string) {
    // BucketSize
    str := fmt.Sprintf("%s_BucketSize:%d\n", prefix, dc.BucketSize)
    buf.WriteString(str)
    // BucketNum
    str = fmt.Sprintf("%s_BucketNum:%d\n", prefix, dc.BucketNum)
    buf.WriteString(str)
    // Count
    str = fmt.Sprintf("%s_Count:%d\n", prefix, dc.Count)
    buf.WriteString(str)
    // Sum
    str = fmt.Sprintf("%s_Sum:%d\n", prefix, dc.Sum)
    buf.WriteString(str)
    // Ave
    str = fmt.Sprintf("%s_Ave:%d\n", prefix, dc.Ave)
    buf.WriteString(str)
    // Counters
    for i := 0; i <= dc.BucketNum; i ++ {
        str = fmt.Sprintf("%s_Counters_%d:%d\n", prefix, i, dc.Counters[i])
        buf.WriteString(str)        
    }    
}
