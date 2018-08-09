/* uint32_slice.go - provide sort and search interface for uint32 slice */
/*
modification history
--------------------
2015/6/23, by Guang Yao, create
*/
/*
DESCRIPTION
*/

package uint32_slice

import (
    "fmt"
    "sort"
)

type Uint32Slice []uint32

// funcs for sort uint32 slice
func (s Uint32Slice) Len() int {
    return len(s)
}

func (s Uint32Slice) Less(i, j int) bool {
    return s[i] < s[j]
}

func (s Uint32Slice) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

// sort a Uint32Slice based on standard lib
func (s Uint32Slice) Sort() {
    sort.Sort(s)
}

type AscendUint32Slice struct {
   slice []uint32
}

// generate a AscendUint32Slice from Uint32Slice
func NewAscendUint32Slice(s Uint32Slice) AscendUint32Slice {
    s.Sort()
    return AscendUint32Slice{s}
}

/*
search the largest value no larger than the key in an ascend slice, based on bisection
e.g., in [10, 20, 30, 35], if key = 25, 20 will be returned.

Params:
    - key: the key to search

Returns:
    - (left-most value in slice, error)
*/
func (s AscendUint32Slice) FindLeftMost(key uint32) (uint32, error) {
    // check whether slice is empty
    if s.slice == nil || len(s.slice) == 0 {
        return 0, fmt.Errorf("slice is empty")
    }

    // check whether key is less than the lower limit
    // key...lower...upper
    if key < s.slice[0] {
        return 0, fmt.Errorf("key[%d] not in scope", key)
    }

    // key is larger than the largest value
    // lower...upper...key
    if key > s.slice[len(s.slice)-1] {
        return s.slice[len(s.slice)-1], nil
    }

    // use binary search to find the leftmost value
    lowerIndex := 0
    upperIndex := len(s.slice) - 1
    curIndex := len(s.slice) / 2
    for {
        switch {
        // cur == key
        case s.slice[curIndex] == key:
            return key, nil
        // cur....key...cur+1
        case s.slice[curIndex] < key && s.slice[curIndex+1] > key:
            return s.slice[curIndex], nil
        // lower...key...cur...upper, decrease cur
        case s.slice[curIndex] > key:
            upperIndex = curIndex
            curIndex = (curIndex + lowerIndex) / 2
        // lower...cur...key...upper, increase cur
        case s.slice[curIndex] < key:
            lowerIndex = curIndex
            curIndex = (curIndex + upperIndex) / 2
        }
    }
}
