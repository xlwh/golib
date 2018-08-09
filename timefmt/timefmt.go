/* timefmt.go - output time in some format  */
/*
modification history
--------------------
2014/4/24, by Zhang Miao, create
2016/8/8, by Guang Yao, add TimestampSplit
*/
/*
DESCRIPTION

Usage:
    import "www.baidu.com/golang-lib/timefmt"
    println(timefmt.CurrTimeGet())  // 2006-01-02 15:04:05
*/
package timefmt

import (
	"fmt"
	"strconv"
	"time"
)

// get current time in format like 2006-01-02 15:04:05
func CurrTimeGet() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

/*
UnixTimeToDateTime - convert from unix time to date-and-time
                     e.g., from 13923223442(int64) to 20140611120132(int64)

Params:
    - unixTime: the number of seconds elapsed since January 1, 1970 UTC.

Returns:
    (datetime, error)
    datetime - yyyymmddhhmmss, e.g., 20140611120132 (int64)
*/
func UnixTimeToDateTime(unixTime int64) int64 {
	// convert from unix time to string of date-and-time
	timeStr := time.Unix(unixTime, 0).Format("20060102150405")

	// convert date-and-time from string to int
	timeInt, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		// this should not happen
		return 0
	}

	return timeInt
}

/*
Split a full time string in format "yyyyMMddHHmmss"(e.g., "20130411143020"), to
"yyyyMMdd" and "HHmmss"

Params:
     - timestamp: time string in format "yyyyMMddHHmmss"
Returns:
     ("yyyyMMdd", "HHmmss", error)
*/
func TimestampSplit(timestamp string) (string, string, error) {
	if len(timestamp) != 14 {
		return "", "", fmt.Errorf("length of timestamp is not as expected: %s", timestamp)
	}

	return timestamp[0:8], timestamp[8:14], nil
}

/*
convert time string of format "2016-01-11 06:12:33" to unix timestamp

Params:
     - timeStr: time string
Returns:
     (timestamp, err)
*/
func StrToUnix(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05 MST", timeStr+" CST")
	if err != nil {
		return 0, fmt.Errorf("wrong time string format")
	}
	return t.Unix(), nil
}

/*
convert unix timestamp to string of format "2006-01-02 15:04:05"

Params:
     - timestamp: unix timestamp
Returns:
     time string
*/
func UnixToStr(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}
