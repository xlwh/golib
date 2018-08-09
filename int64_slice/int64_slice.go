/* int64_slice.go - provide sort interfaces for int64 slice */
/*
modification history
--------------------
2016/2/2, by Guang Yao, create
*/
/*
DESCRIPTION
*/

package int64_slice

import (
	"sort"
)

type Int64Slice []int64

// funcs for sort int64 slice
func (s Int64Slice) Len() int {
	return len(s)
}

func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Int64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// sort a Int64Slice based on standard lib
func (s Int64Slice) Sort() {
	sort.Sort(s)
}
