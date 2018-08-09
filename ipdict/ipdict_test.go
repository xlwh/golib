/* ipdict.go - test file of ipdict */
/*
modification history
--------------------
2014/7/14, by Li Bingyi, create
*/
/*
DESCRIPTION
*/

package ipdict

import (
    "fmt"
    "testing"
)

import (
    "www.baidu.com/golang-lib/net_util"
)

type ipStr struct {
    start string
    end   string
}

type ipStrs []ipStr

// util func for unit test
// load ip to IPItems from struct string
func loadIPStr(ips ipStrs) (*IPItems, error) {

    ipItems, err := NewIPItems(1000,1000)
    if err != nil {
        return nil, err
    }
    for _, ip := range ips {
        startIP := net_util.ParseIPv4(ip.start)
        endIP := net_util.ParseIPv4(ip.end)
        err := ipItems.InsertPair(startIP, endIP)
        if err != nil {
            return nil, err
        }
    }

    return ipItems, nil
}

// util func for unit test
// check whether two ipPairs is equal
// return true if equal else return false
func checkEqual(src, dst ipPairs) bool {
    if len(src) != len(dst) {
        fmt.Println("checkEqual(): len not equal!")
        return false
    } else {
        for i := 0; i < len(src); i++ {
            if src[i].startIP != dst[i].startIP {
                fmt.Printf("checkEqual(): start element [%d] and [%d] are not equal!\n",
                    src[i].startIP, dst[i].startIP)
                return false
            }

            if src[i].endIP != dst[i].endIP {
                fmt.Printf("checkEqual(): end element [%d] and [%d] are not equal!\n",
                    src[i].endIP, dst[i].endIP)
                return false
            }

        }
    }

    return true
}

// < case
func TestLess_Case0(t *testing.T) {
    var p ipPairs

    ipStr1 := "1.1.1.1"
    ipStr2 := "2.2.2.2"

    ip1, _ := net_util.IPv4StrToUint32(ipStr1)
    ip2, _ := net_util.IPv4StrToUint32(ipStr2)

    p = append(p, ipPair{ip1, ip1})
    p = append(p, ipPair{ip2, ip2})

    if p.Less(0, 1) {
        t.Errorf("Less(): %s >= %s", ipStr1, ipStr2)
    }

}

//  = case
func TestLess_Case1(t *testing.T) {
    var p ipPairs

    ipStr := "1.1.1.1"

    ip, _ := net_util.IPv4StrToUint32(ipStr)

    p = append(p, ipPair{ip, ip})
    p = append(p, ipPair{ip, ip})

    if !p.Less(0, 1) {
        t.Errorf("Less(): %s < %s", ipStr, ipStr)
    }

}

//  > case
func TestLess_Case2(t *testing.T) {
    var p ipPairs

    ipStr1 := "2.2.2.2"
    ipStr2 := "1.1.1.1"

    ip1, _ := net_util.IPv4StrToUint32(ipStr1)
    ip2, _ := net_util.IPv4StrToUint32(ipStr2)

    p = append(p, ipPair{ip1, ip1})
    p = append(p, ipPair{ip2, ip2})

    if !p.Less(0, 1) {
        t.Errorf("Less(): %s < %s", ipStr1, ipStr2)
    }
}

// normal case
func TestSwap_Case0(t *testing.T) {
    var p ipPairs

    ipStr1 := "1.1.1.1"
    ipStr2 := "2.2.2.2"

    ip1, _ := net_util.IPv4StrToUint32(ipStr1)
    ip2, _ := net_util.IPv4StrToUint32(ipStr2)

    p = append(p, ipPair{ip1, ip1})
    p = append(p, ipPair{ip2, ip2})

    p.Swap(0, 1)

    if ip1 != p[1].startIP || ip1 != p[1].endIP {
        t.Errorf("Swap(): %s and %s swap failed!", ipStr1, ipStr2)
    }

    if ip2 != p[0].startIP || ip2 != p[0].endIP {
        t.Errorf("Swap(): %s and %s swap failed!", ipStr1, ipStr2)
    }

}

// startIP < endIP case
func TestInsert_Case0(t *testing.T) {
    ipItems, err := NewIPItems(1000,1000)
    if err != nil {
        t.Error(err.Error())
    }

    startIPStr := "1.1.1.1"
    endIPStr := "2.2.2.2"

    startIP := net_util.ParseIPv4(startIPStr)
    endIP := net_util.ParseIPv4(endIPStr)

    err = ipItems.InsertPair(startIP, endIP)
    if err != nil {
        t.Errorf("insert(): %s!", err.Error())
    }
}

// startIP = endIP case
func TestInsert_Case1(t *testing.T) {
    ipItems, err := NewIPItems(1000,1000)
    if err != nil {
        t.Error(err.Error())
    }

    startIPStr := "1.1.1.1"
    endIPStr := "1.1.1.1"

    startIP := net_util.ParseIPv4(startIPStr)
    endIP := net_util.ParseIPv4(endIPStr)

    err = ipItems.InsertPair(startIP, endIP)
    if err != nil {
        t.Error(err.Error())
    }
}

// startIP > endIP case
func TestInsert_Case2(t *testing.T) {
    ipItems, err := NewIPItems(1000,1000)
    if err != nil {
        t.Error(err.Error())
    }

    startIPStr := "2.2.2.2"
    endIPStr := "1.1.1.1"

    startIP := net_util.ParseIPv4(startIPStr)
    endIP := net_util.ParseIPv4(endIPStr)

    err = ipItems.InsertPair(startIP, endIP)
    if err == nil {
        t.Error(err.Error())
    }
}

func TestCheckMerge_Case0(t *testing.T) {
    ips := ipStrs{
        {
            "10.26.74.55",
            "10.26.74.255",
        },
        {
            "0.0.0.0",
            "0.0.0.0",
        },
        {
            "10.12.14.2",
            "10.26.74.105",
        },
    }

    ipItems, err := loadIPStr(ips)
    if err != nil {
        t.Error(err.Error())
    }
    
    ret := ipItems.checkMerge(0, 2)
    if ret != 1 {
        t.Errorf("checkMerge(): failed! ret:%d", ret)
    }

}

// len 0 case
func TestMergeItems_Case0(t *testing.T) {
    ips := ipStrs{}

    ipItems, err := loadIPStr(ips)
    if err != nil {
        t.Error(err.Error())
    }

    if ipItems.mergeItems() != 0 {
        t.Errorf("mergeItems(): failed!")
    }
}


// convert ip string to uint32 format, without error handle
func ipv4StrToUint32(ipStr string) uint32 {
    ip := net_util.ParseIPv4(ipStr)
    
    var ipNum uint32
    var tmp uint32

    for i, b := range ip {
        tmp = uint32(b)
        ipNum = ipNum | (tmp << uint((3-i)*8))
    }


    return ipNum
}


// normal case
func TestSort_Case0(t *testing.T) {
    ips := ipStrs{
        {
            "10.26.74.55",
            "10.26.74.255",
        },
        {
            "10.21.34.5",
            "10.23.77.100",
        },
        {
            "10.12.14.2",
            "10.12.14.50",
        },
    }

    IPs := ipPairs{
        {
            ipv4StrToUint32("10.26.74.55"),
            ipv4StrToUint32("10.26.74.255"),
        },
        {
            ipv4StrToUint32("10.21.34.5"),
            ipv4StrToUint32("10.23.77.100"),
        },
        {
            ipv4StrToUint32("10.12.14.2"),
            ipv4StrToUint32("10.12.14.50"),
        },
    }

    ipItems, err := loadIPStr(ips)
    if err != nil {
        t.Error(err.Error())
    }

    ipItems.Sort()

    if !checkEqual(ipItems.items, IPs) {
        t.Errorf("checkEqual(): failed!")
    }
}

// merge case
func TestSort_Case1(t *testing.T) {

    ips := ipStrs{
        {
            "10.26.74.55",
            "10.26.74.255",
        },
        {
            "10.23.77.88",
            "10.23.77.240",
        },
        {
            "10.21.34.5",
            "10.23.77.100",
        },
        {
            "10.12.14.2",
            "10.12.14.50",
        },
    }

    IPs := ipPairs{
        {
            ipv4StrToUint32("10.26.74.55"),
            ipv4StrToUint32("10.26.74.255"),
        },
        {
            ipv4StrToUint32("10.21.34.5"),
            ipv4StrToUint32("10.23.77.240"),
        },
        {
            ipv4StrToUint32("10.12.14.2"),
            ipv4StrToUint32("10.12.14.50"),
        },
    }

    ipItems, err := loadIPStr(ips)
    if err != nil {
        t.Error(err.Error())
    }

    ipItems.Sort()

    if !checkEqual(ipItems.items, IPs) {
        t.Errorf("checkEqual(): failed!")
    }

}

// merge case
func TestSort_Case2(t *testing.T) {

    ips := ipStrs{
        {
            "10.26.74.55",
            "10.26.74.255",
        },
        {
            "10.23.74.8",
            "10.26.74.55",
        },
    }

    IPs := ipPairs{
        {
            ipv4StrToUint32("10.23.74.8"),
            ipv4StrToUint32("10.26.74.255"),
        },
    }

    ipItems, err := loadIPStr(ips)
    if err != nil {
        t.Error(err)
    }

    ipItems.Sort()

    if !checkEqual(ipItems.items, IPs) {
        t.Errorf("checkEqual(): failed!")
    }

}

// total merge case
func TestSort_Case3(t *testing.T) {

    ips := ipStrs{
        {
            "10.26.74.55",
            "10.26.74.255",
        },
        {
            "10.23.77.88",
            "10.23.77.240",
        },
        {
            "10.21.34.5",
            "10.23.77.100",
        },
        {
            "10.12.14.2",
            "10.30.74.5",
        },
    }

    IPs := ipPairs{
        {
            ipv4StrToUint32("10.12.14.2"),
            ipv4StrToUint32("10.30.74.5"),
        },
    }

    ipItems, err := loadIPStr(ips)
    if err != nil {
        t.Error(err.Error())
    }

    ipItems.Sort()

    if !checkEqual(ipItems.items, IPs) {
        t.Errorf("checkEqual(): failed!")
    }
}
