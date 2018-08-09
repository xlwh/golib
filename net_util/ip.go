/* ip.go - provide something different to net/ip.go of sys default  */
/*
modification history
--------------------
2014/8/29, by Zhang Miao, create
2015/4/8, by Zhang Miao, add IPv4ToUint32() and Uint32ToIPv4()
2015/6/23, by Guang Yao, add IPv4StrToUint32() and Uint32ToIPv4Str()
2017/12/11, by Taochunhua, add IsIPv4String() and IsPrivateIp()
*/
/*
DESCRIPTION
*/
package net_util

import (
    "bytes"
    "fmt"
    "net"
)

//IpRange - a structure that holds the start and end of a range of ip addresses
type IpRange struct {
	start net.IP
	end   net.IP
}

// private ip range
var privateRanges = []IpRange{
	IpRange{
		start: net.ParseIP("10.0.0.0").To4(),
		end:   net.ParseIP("10.255.255.255").To4(),
	},
	IpRange{
		start: net.ParseIP("172.16.0.0").To4(),
		end:   net.ParseIP("172.31.255.255").To4(),
	},
	IpRange{
		start: net.ParseIP("192.168.0.0").To4(),
		end:   net.ParseIP("192.168.255.255").To4(),
	},
}

// InRange - check to see if a given ip address is within a range given
func InRange(r IpRange, ip net.IP) bool {
	if bytes.Compare(ip, r.start) >= 0 && bytes.Compare(ip, r.end) <= 0 {
		return true
	}
	return false
}

/*
parse IP addr from string to net.IP

Params:
    - s: IP addr in string, e.g., "220.181.112.244"

Returns:
    IP addr in net.IP
*/
func ParseIPv4(s string) net.IP {
    ip := net.ParseIP(s)

    if ip != nil {
        ip = ip.To4()
    }

    return ip
}

/*
convert net.IP to uint32

e.g., 1.2.3.4 to 0x01020304

Params:
    - ipBytes: IPv4 addr in net.IP

Returns:
    IPv4 addr in uint32
*/
func IPv4ToUint32(ipBytes net.IP) (uint32, error) {
    if len(ipBytes) != 4 {
        return 0, fmt.Errorf("ip bytes len: %d", len(ipBytes))
    }

    var ipNum uint32
    var tmp uint32

    for i, b := range ipBytes {
        tmp = uint32(b)
        ipNum = ipNum | (tmp << uint((3-i)*8))
    }

    return ipNum, nil
}

/*
convert IPv4 string to uint32

e.g., "1.2.3.4" to 0x01020304

Params:
    - ipStr: IPv4 addr in string

Returns:
    IPv4 addr in uint32
*/
func IPv4StrToUint32(ipStr string) (uint32, error) {
    ip := ParseIPv4(ipStr)
    if ip == nil {
        return 0, fmt.Errorf("invalid IPv4 addr string: %s", ipStr)
    }

    return IPv4ToUint32(ip)
}

/*
convert uint32 net.IP

e.g., 0x01020304 to 1.2.3.4

Params:
    - ipNum: IPv4 addr in uint32

Returns:
    IPv4 addr in net.IP
*/
func Uint32ToIPv4(ipNum uint32) net.IP {
    var ipBytes [4]byte

    for i := 0; i < 4; i++ {
        ipBytes[3-i] = byte(ipNum & 0xFF)
        ipNum = ipNum >> 8
    }

    return net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3]).To4()
}

/*
convert uint32 to str

e.g., 0x01020304 to "1.2.3.4"

Params:
    - ipNum: IPv4 addr in uint32

Returns:
    IPv4 addr in string
*/
func Uint32ToIPv4Str(ipNum uint32) string {
    str := fmt.Sprintf("%d.%d.%d.%d", byte(ipNum>>24), byte(ipNum>>16), byte(ipNum>>8), byte(ipNum))

    return str
}

/*
Check input is ipv4 address or not.

param:
    - input: a string
return:
    bool
*/
func IsIPv4Address(input string) bool {
	ip := net.ParseIP(input).To4()
	if ip == nil {
		return false
	}
	return true
}

/*
Check to see if an ip is in a private subnet.

param:
    - input: an ip string
return:
    bool
*/
func IsPrivateIp(input string) bool {
	if ip := net.ParseIP(input).To4(); ip != nil {
		for _, r := range privateRanges {
			if InRange(r, ip) {
				return true
			}
		}
	}
	return false
}
