/* pop_traffic_test.go - test for pop_traffic.go */
/*
modification history
--------------------
2016/07/28, by Guang Yao, create
*/
/*
DESCRIPTION
*/

package trident

import (
	"testing"
)

func TestPopTrafficGet(t *testing.T) {
	pops := []Pop{
		Pop{"HZ01", "TELECOM"},
		Pop{"NJ01", "TELECOM"},
	}

	popTraffic, err := GetPopTraffic(pops)
	if err != nil {
		t.Errorf("GetPopTraffic(): %v", err)
	}

	if len(popTraffic) != 2 {
		t.Errorf("unexpected result: %+v", popTraffic)
	}
}
