/* string_slice.go - []string related features  */
/*
modification history
--------------------
2015/8/13, by Zhang Miao, create
2015/12/1, by Guang Yao, add EqualIgnoreOrder
*/
/*
DESCRIPTION
*/
package string_slice

import (
	"sort"
)

// whether given string is in string slice
func InSlice(a string, slice []string) bool {
	for _, b := range slice {
		if a == b {
			return true
		}
	}
	return false
}

// sort for []string
type StringSlice []string

// funcs for sort string slice
func (s StringSlice) Len() int {
	return len(s)
}

func (s StringSlice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s StringSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// sort a Uint32Slice based on standard lib
func (s StringSlice) Sort() {
	sort.Sort(s)
}

// whether two string slice have same elements, ignoring order
func EqualIgnoreOrder(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	// sort and compare
	s1 := StringSlice(slice1)
	s2 := StringSlice(slice2)
	sort.Sort(s1)
	sort.Sort(s2)
	for i, str := range s1 {
		if str != s2[i] {
			return false
		}
	}

	return true
}
