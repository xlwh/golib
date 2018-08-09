/* iptable_test.go -
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
    "testing"
)

import (
    "www.baidu.com/golang-lib/net_util"
)

func loadIP(ips ipPairs) (*IPItems, error) {

    ipItems, err := NewIPItems(1000, 1000)
    if err != nil {
        return nil, err
    }

    for _, ip := range ips {
        startIP := net_util.Uint32ToIPv4(ip.startIP)
        endIP := net_util.Uint32ToIPv4(ip.endIP)
        err := ipItems.InsertPair(startIP, endIP)
        if err != nil {
            return nil, err
        }
    }

    return ipItems, nil
}

/* Update provides for thread-safe switching items */
func TestUpdate(t *testing.T) {
    table := NewIPTable()

    ipItems, err := NewIPItems(1000, 1000)
    if err != nil {
        t.Error(err.Error()) 
    }

    zeroIP := uint32(0)

    ipItems.items = append(ipItems.items, ipPair{zeroIP, zeroIP})
    table.Update(ipItems)
    
    if ipItems.Length() != 1 {
        t.Errorf("TestItemLength): itemNum [%d] != 1", ipItems.Length())
    }

}

/* Search provides for binary search IP in dict */
func TestSearch(t *testing.T) {
    // Create table
    table := NewIPTable()

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
        {
            ipv4StrToUint32("2.2.2.2"),
            ipv4StrToUint32("2.2.2.2"),
        },
    }

    ipItems, err := loadIP(IPs)
    if err != nil {
        t.Error(err.Error())
    }

    ipItems.Sort()

    // Switch items
    table.Update(ipItems)

    // Search items
    if !table.Search(net_util.ParseIPv4("10.12.14.12")) {
        t.Errorf("TestSearch(): 10.12.14.12 not hit")
    }

    if !table.Search(net_util.ParseIPv4("2.2.2.2")) {
        t.Errorf("TestSearch(): 2.2.2.2 not hit")
    }

    if table.Search(net_util.ParseIPv4("1.1.1.1")) {
        t.Errorf("TestSearch(): 1.1.1.1 hit")
    }
}
