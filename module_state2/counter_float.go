/* counter.go - counters  */
/*
modification history
--------------------
2014/4/24, by Zhang Miao, create
*/
/*
DESCRIPTION
*/

package module_state2

/* flat counters for float64   */
type FloatCounters map[string]float64

// create new Counters for float64
func NewFloatCounters() FloatCounters {
	floatCounters := make(FloatCounters)
	return floatCounters
}
