/* ipdict.go - Provide startIP, endIP, value dict operation */
/*
modification history
--------------------
2016/12/19, by Zhang Jiyang, create
*/
/*
 */

package iptable

import (
	"errors"
	"fmt"
	"net"
	"sort"
)

import (
	"www.baidu.com/golang-lib/byte_pool"
	"www.baidu.com/golang-lib/net_util"
)

var (
	ErrDictFull    = errors.New("Dict is Full")
	ErrTooLargeVal = errors.New("Value is too large")
)

// IPPair struct for pair ip
type IPPair struct {
	startIP uint32 // start ip addr in uint32
	endIP   uint32 // end ip addr in uint32
	index   int    // value index at byte slice pool
}

// IPParis ip pairs slice
type IPPairs []IPPair

// Less func used by pkg sort
func (ip IPPairs) Less(i, j int) bool {
	return ip[i].startIP < ip[j].startIP
}

// Swap func used by pkg sort
func (ip IPPairs) Swap(i, j int) {
	ip[i], ip[j] = ip[j], ip[i]
}

// Len func used by pkg sort
func (ip IPPairs) Len() int {
	return len(ip)
}

// IPDict manage ip pairs
type IPDict struct {
	ipPairs     IPPairs
	valPool     *byte_pool.BytePool // byte pool store value
	maxItemSize int                 // max item size
	length      int                 // length of ipdict
	capacity    int                 // capacity of ipdict
}

func NewIPDict(capacity int, maxItemSize int) (*IPDict, error) {
	if capacity <= 0 || maxItemSize <= 0 {
		return nil, fmt.Errorf("capacify and maxItemSize should > 0, while %d:%d", capacity, maxItemSize)
	}

	d := new(IPDict)
	d.ipPairs = make(IPPairs, 0, capacity)
	d.valPool = byte_pool.NewBytePool(capacity, maxItemSize)
	d.maxItemSize = maxItemSize
	d.length = 0
	d.capacity = capacity

	return d, nil
}

// Add, add item into dict
func (d *IPDict) Add(startIP, endIP net.IP, val string) error {
	// check dict full
	if d.full() {
		return ErrDictFull
	}

	// check value length
	if len(val) > d.maxItemSize {
		return ErrTooLargeVal
	}

	var sIP, eIP uint32
	var err error

	// convert start ip
	if sIP, err = net_util.IPv4ToUint32(startIP.To4()); err != nil {
		return err
	}

	// convert end ip
	if eIP, err = net_util.IPv4ToUint32(endIP.To4()); err != nil {
		return err
	}

	// check, startIP should <= endIP
	if sIP > eIP {
		return fmt.Errorf("startIP shoule <= endIP while %s:%s", startIP.String(), endIP.String())
	}

	ipPair := IPPair{startIP: sIP, endIP: eIP, index: d.length}
	d.ipPairs = append(d.ipPairs, ipPair)
	d.valPool.Set(int32(d.length), []byte(val))
	d.length++

	return nil
}

func (d *IPDict) full() bool {
	return d.length >= d.capacity
}

// SortAndCheck, sort and check ipParis
// same ip have more than one value is not allowed,
// so endIP of former line must < startIP of behind line
//
// example:
//    1.1.1.1 3.3.3.3 val1
//    2.2.2.2 4.4.4.4 val2
// is not allowed, because IP 2.2.2.2 have two different values
//
func (d *IPDict) SortAndCheck() error {
	sort.Sort(d.ipPairs)

	// check ipPairs
	for i := 0; i < len(d.ipPairs)-1; i++ {
		pre := d.ipPairs[i]
		cur := d.ipPairs[i+1]

		if pre.endIP >= cur.startIP {
			return fmt.Errorf("confilt ipPairs %d-%d %d-%d", pre.startIP, pre.endIP, cur.startIP, cur.endIP)
		}
	}

	return nil
}

// Search IP in ipdict and return value
// Params:
//    - ip: net.IP, ip addr to search
// Returns:
//    - (val, true): if search success
//    - ("", false): if search failed
func (d *IPDict) Search(ip net.IP) (string, bool) {
	uIP, err := net_util.IPv4ToUint32(ip.To4())
	if err != nil {
		return "", false
	}

	i := sort.Search(d.length,
		func(i int) bool { return d.ipPairs[i].endIP >= uIP })

	if i < d.length && d.ipPairs[i].startIP <= uIP {
		return string(d.valPool.Get(int32(d.ipPairs[i].index))), true
	}

	return "", false
}
