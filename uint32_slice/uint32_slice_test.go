/* uint32slice_test.go - test for uint32slice.go  */
/*
modification history
--------------------
2015/6/23, by Guang Yao, create
*/

package uint32_slice

import (
    "testing"
)

func TestSort(t *testing.T) {
    s := Uint32Slice{20, 15, 30}
    s_sorted := Uint32Slice{15, 20, 30}
    s.Sort()

    for i, _ := range s {
        if s[i] != s_sorted[i] {
            t.Errorf("Invalid Sort result: [%v]=>[%v]", Uint32Slice{20, 15, 30}, s)
        }
    }
}

func TestFindLeftMost(t *testing.T) {
    s := Uint32Slice{}
    as := NewAscendUint32Slice(s)
    if _, err := as.FindLeftMost(20); err == nil {
        t.Errorf("err should happen in search empty slice")
    }

    s = Uint32Slice{20, 25}
    as = NewAscendUint32Slice(s)
    if _, err := as.FindLeftMost(10); err == nil {
        t.Errorf("err should happen in search key not in scope")
    }

    s = Uint32Slice{20, 25}
    as = NewAscendUint32Slice(s)
    if key, err := as.FindLeftMost(20); err != nil {
        t.Errorf("err happen in find: %v", err)
    } else if key != 20 {
        t.Errorf("search by 20 in %v, 20 should be returned, but %d returned", s, key)
    }

    if key, err := as.FindLeftMost(21); err != nil {
        t.Errorf("err happen in find: %v", err)
    } else if key != 20 {
        t.Errorf("search by 21 in %v, 20 should be returned, but %d returned", s, key)
    }

    if key, err := as.FindLeftMost(30); err != nil {
        t.Errorf("err happen in find: %v", err)
    } else if key != 25 {
        t.Errorf("search by 30 in %v, 25 should be returned, but %d returned", s, key)
    }

    // large slice
    s = Uint32Slice{}
    for i := 0; i < 9999999; i++ {
        s = append(s, uint32(i*2))
    }
    as = NewAscendUint32Slice(s)
    if key, err := as.FindLeftMost(3333332); err != nil {
        t.Errorf("err happen in find: %v", err)
    } else if key != 3333332 {
        t.Errorf("search by 3333333 in [1*2...9999999*2], 3333332 should be"+
            " returned, but %d returned", key)
    }
}
