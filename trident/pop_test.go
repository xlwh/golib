/* pop_test.go - test for pop.go */
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

func TestGetPop(t *testing.T) {
	pops, err := GetPop()
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}

	if len(pops) == 0 {
		t.Errorf("empty result returned")
	}
}
