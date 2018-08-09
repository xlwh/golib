/* time_wait.go - implement WaitTill()  */
/*
modification history
--------------------
2014/12/22, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package time_wait

import (
	"time"
)

/*
Wait until toTime

Params:
    - toTime: time to wait until. the number of seconds elapsed since January 1, 1970 UTC.
*/
func WaitTill(toTime int64) {
	waitSecs := toTime - time.Now().Unix()
	if waitSecs > 0 {
		time.Sleep(time.Second * time.Duration(waitSecs))
	}
}

/*
Calc the nearest time from now, given cycle and offset

Params:
	- cycle: cycle in seconds
	- offset: offset of the next time; in seconds

Return:
	- timestamp of next time
*/
func CalcNextTime(cycle int64, offset int64) int64 {
	current := time.Now().Unix()

	if current%cycle == 0 {
		return current + offset
	} else {
		return current - current%cycle + cycle + offset
	}
}
