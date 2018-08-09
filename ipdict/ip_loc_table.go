/* ip_loc_table.go - ip location table for dict-server */
/*
modification history
--------------------
2016/6/6, by Jiang Hui, create
*/
/*
DESCRIPTION
file fromat like:
#...
#version: v2.0.1
#...
#ipstart|ipend|country|isp|province|city|county|country confidence|isp confidence|province confidence|city confidence|county confidence
0.0.0.0|0.255.255.255|ZZ|None|None|None|None|100|0|0|0|0

*/
package ipdict

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sort"
)

import (
	"www.baidu.com/golang-lib/net_util"
)

const (
	IP_SIZE     = 4
	HEADER_LEN  = 8
	MAX_LINE    = 1000000
	MAX_LOC_LEN = 1024
)

//uppercasing the first letter for binary lib
type ipLocation struct {
	startIp  uint32
	endIp    uint32
	location []byte
}

//[]byte to string ,remove last 0 in []bytes
func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

type IpLocationTable struct {
	Version   string
	maxSize   uint32
	LocLen    uint32
	offset    uint32
	locations []byte
}

func NewIpLocationTable(maxSize uint32, LocLen uint32) (*IpLocationTable, error) {
	//maxSize max is MAX_LINE
	if maxSize == 0 || maxSize > MAX_LINE {
		return nil, fmt.Errorf("NewIpLocationTable caused by maxSize :%d", maxSize)
	}

	//LocLen max size is MAX_LOC_LEN
	if LocLen == 0 || LocLen > MAX_LOC_LEN {
		return nil, fmt.Errorf("NewIpLocationTable caused by LocLen :%d", LocLen)
	}

	ipLocTable := new(IpLocationTable)
	ipLocTable.maxSize = maxSize
	ipLocTable.offset = 0
	ipLocTable.LocLen = LocLen
	ipLocTable.locations = make([]byte, (HEADER_LEN+LocLen)*maxSize, (HEADER_LEN+LocLen)*maxSize)
	return ipLocTable, nil
}

//write ipLocation Struct to locations by [HeaderLen+t.LocLen]byte
func (t *IpLocationTable) writeStruct(idx uint32, ipLoc ipLocation) {
	sOffset := idx * (t.LocLen + HEADER_LEN)
	binary.LittleEndian.PutUint32(t.locations[sOffset:sOffset+IP_SIZE], ipLoc.startIp)
	binary.LittleEndian.PutUint32(t.locations[sOffset+IP_SIZE:sOffset+HEADER_LEN], ipLoc.endIp)
	copy(t.locations[sOffset+HEADER_LEN:sOffset+HEADER_LEN+t.LocLen], ipLoc.location)
}

//read ipLocation from locations by idx
func (t *IpLocationTable) readStruct(idx uint32) ipLocation {
	var ipLoc ipLocation
	ipLoc.location = make([]byte, t.LocLen)
	sOffset := idx * (t.LocLen + HEADER_LEN)
	ipLoc.startIp = binary.LittleEndian.Uint32(t.locations[sOffset : sOffset+IP_SIZE])
	ipLoc.endIp = binary.LittleEndian.Uint32(t.locations[sOffset+IP_SIZE : sOffset+HEADER_LEN])
	ipLoc.location = t.locations[sOffset+HEADER_LEN : sOffset+HEADER_LEN+t.LocLen]
	return ipLoc
}

//add ip location dict to locations buffer
//assume add startIP:EndIP have been sorted
//every add startIP:EndIP region does not overlap
func (t *IpLocationTable) Add(startIP, endIP net.IP, location string) error {
	if t.offset >= t.maxSize {
		return fmt.Errorf("Add():caused by table is full")
	}

	if bytes.Compare(startIP, endIP) == 1 {
		return fmt.Errorf("Add(): err startIPStr %s > endIPStr %s",
			startIP.String(), endIP.String())
	}

	// convert start IP to Uint32 format
	sIP, err := net_util.IPv4ToUint32(startIP)
	if err != nil {
		return err
	}

	// convert endIP to Uint32 format
	eIP, err := net_util.IPv4ToUint32(endIP)
	if err != nil {
		return err
	}

	//write unit(startip,endip,location) to locations buffer
	var loc ipLocation
	loc.startIp = sIP
	loc.endIp = eIP
	loc.location = make([]byte, t.LocLen)
	copy(loc.location[0:t.LocLen], location)
	t.writeStruct(uint32(t.offset), loc)

	t.offset++
	return nil
}

//binary search pool to find the ip's location
//search sort of array(order from small to large)
func (t *IpLocationTable) Search(cip net.IP) (string, error) {
	ip := cip.To4()
	ipAddr, err := net_util.IPv4ToUint32(ip)
	if err != nil {
		return "", err
	}

	indexLen := t.offset
	if indexLen == 0 {
		return "", fmt.Errorf("Search() error caused by locations is null")
	}

	idx := sort.Search(int(indexLen),
		func(i int) bool {
			s := uint32(i) * (HEADER_LEN + t.LocLen)
			e := uint32(i)*(HEADER_LEN+t.LocLen) + IP_SIZE
			b := t.locations[s:e]
			value := binary.LittleEndian.Uint32(b)
			return value >= ipAddr
		})

	//get idx corresponding ip section's first ip
	var fristIp uint32
	if uint32(idx) <= indexLen-1 {
		s := uint32(idx) * (HEADER_LEN + t.LocLen)
		e := uint32(idx)*(HEADER_LEN+t.LocLen) + IP_SIZE
		fristIp = binary.LittleEndian.Uint32(t.locations[s:e])
	}

	var preIdx uint32

	if uint32(idx) == indexLen {
		//consider ipAdd last element(uint32(idx) == indexLen)
		preIdx = uint32(indexLen - 1)
	} else if fristIp == ipAddr || idx == 0 {
		//consider ipAdd locate in frist section (idx == 0)
		//consider ipAdd is first ip in ip's section(fristIp == ipAddr)
		preIdx = uint32(idx)
	} else {
		//other think ipAdd location previous section
		preIdx = uint32(idx - 1)
	}

	//read unit(startip,endip,location) from locations buffer
	loc := t.readStruct(preIdx)
	if ipAddr <= loc.endIp && ipAddr >= loc.startIp {
		return byteString(loc.location[0:]), nil
	}
	return "", fmt.Errorf("Search() error caused by the ip's location does not exist")
}
