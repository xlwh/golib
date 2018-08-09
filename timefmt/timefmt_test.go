/* timefmt_test.go - test for timefmt.go  */
/*
modification history
--------------------
2014/4/24, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package timefmt

import (
	"testing"
)

func TestCurrTimeGet(t *testing.T) {
	curr := CurrTimeGet()
	println(curr)
}

func TestUnixTimeToDateTime(t *testing.T) {
	dateTime := UnixTimeToDateTime(1419238997)
	if dateTime != 20141222170317 {
		t.Errorf("err in UnixTimeToDateTime(), ok:20141222170317, now:%d", dateTime)
	}
}

func TestTimestampSplit(t *testing.T) {
	// good case
	timestr := "20160808144120"
	date, time, err := TimestampSplit(timestr)
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}
	if date != "20160808" {
		t.Errorf("unexpected split result(date): %s", date)
	}
	if time != "144120" {
		t.Errorf("unexpected split result(time): %s", time)
	}

	// bad case 1:
	timestr = "2016080814412011111"
	_, _, err = TimestampSplit(timestr)
	if err == nil {
		t.Errorf("err should happen for : %s", timestr)
	}

	// bad case 2:
	timestr = "20160808144"
	_, _, err = TimestampSplit(timestr)
	if err == nil {
		t.Errorf("err should happen for : %s", timestr)
	}
}

func TestStrToUnix(t *testing.T) {
	ts, err := StrToUnix("2017-05-10 19:35:41")
	if err != nil {
		t.Errorf("unexpected error")
		return
	}
	if ts != 1494416141 {
		t.Errorf("wrong timestamp")
		return
	}
}

func TestUnixToStr(t *testing.T) {
	if UnixToStr(1494416141) != "2017-05-10 19:35:41" {
		t.Errorf("wrong time string")
		return
	}
}
