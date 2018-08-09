/* ipdict.go - Provides for thread-safe IP dict operation
/*
modification history
--------------------
2014/7/7, by Li Bingyi, create
2014/9/11, by zhangjiyang01@baidu.com, modify
    - use hash table for single IP
*/
/*
DESCRIPTION

    Provides for thread-safe IP dict operation .
    Dict in memory is organized as a descending struct slice according the startIP
    EndIP >= startIP in each item pair
    -------------
    startIP endIP
    startIP endIP
    startIP endIP
      ...    ...
    -------------

    See txt_load.go for more detail usage info.

Usage:
    import (
        "net"
    )

    import (
        import "www.baidu.com/golang-lib/ipdict"
    )

    var hit   bool
    var table *ipdict.IPTable
    var items *ipdict.IPItems
    var startIP, endIP, srcIP net.IP

    // Create table
    table = ipdict.NewIPTable()

    // Create items
    items = ipdict.NewIPItems()

    // Insert IP pair into the items
    items.InsertPair(startIP, endIP)
    items.InsertSingle(IP)
    // Switch items
    table.Update(items)

    // Search items
    hit = table.Search(srcIP)
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
    "www.baidu.com/golang-lib/hash_set"
    "www.baidu.com/golang-lib/net_util"
)

const (
    IP_LENGTH = 4 // length of net.IP after converted by net_util.ParseIPV4
)

var ZERO_IP = uint32(0)

/* implement Hash method for hashSet
 * convert net.IP to type uint64
 */
func Hash(ip []byte) uint64 {
    var num uint32
    ipReader := bytes.NewReader(ip)

    binary.Read(ipReader, binary.BigEndian, &num)
    return uint64(num)
}

/* struct for Pair IP */
type ipPair struct {
    startIP uint32
    endIP   uint32
}

type ipPairs []ipPair

/* IPItems manage single IP(hashSet) and ipPairs */
type IPItems struct {
    ipSet   *hash_set.HashSet
    items   ipPairs
    Version string
}

/* create new IPItems */
func NewIPItems(maxSingleIPNum int, maxPairIPNum int) (*IPItems, error) {
    // maxSingleIPNum && maxPairIPNum must >= 0 
    if maxSingleIPNum < 0 || maxPairIPNum < 0 {
        return nil, fmt.Errorf("SingleIPNum/PairIPNum must >= 0")
    }
    
    var err error
    ipItems := new(IPItems)
    
    // create a hashSet for single IPs
    isFixedSize := true  // ip address is fixed size(IP_LENGTH)
    maxSingleIPNum += 1  // +1, hash_set don't support maxSingleIPNum == 0 
    ipItems.ipSet, err = hash_set.NewHashSet(maxSingleIPNum, IP_LENGTH, isFixedSize, Hash)
    if err != nil {
        return nil, err
    }

    // create item array for pair IPs
    ipItems.items = make(ipPairs, 0, maxPairIPNum)
    return ipItems, nil
}

/* IPItems should implement Len() for calling sort.Sort(items) */
func (items ipPairs) Len() int {
    return len(items)
}

/* IPItems should implement Less(int, int) for calling sort.Sort(items) */
func (items ipPairs) Less(i, j int) bool {
    return items[i].startIP >= items[j].startIP
}

/* IPItems should implement Swap(int, int) for calling sort.Sort(items) */
func (items ipPairs) Swap(i, j int) {
    items[i], items[j] = items[j], items[i]
}

/* checkMerge merge items between index i and j in sorted items.
   If items[i] and items[j] can merge, then merge all items between index i and j
   Others do not merge.
   Constraint: j > i, items[j].startIP >= items[i].startIP
*/
func (ipItems *IPItems) checkMerge(i, j int) int {
    var mergedNum int

    items := ipItems.items

    if items[j].endIP >= items[i].startIP {
        items[i].startIP = items[j].startIP
        if items[j].endIP >= items[i].endIP {
            items[i].endIP = items[j].endIP
        }

        items[j].startIP = ZERO_IP
        items[j].endIP = ZERO_IP

        mergedNum++

        // Merge items [i+1, j)
        for k := i + 1; k < j; k++ {
            if items[k].endIP == ZERO_IP {
                continue
            }

            items[k].startIP = ZERO_IP
            items[k].endIP = ZERO_IP
            mergedNum++
        }
    }

    return mergedNum
}

/* mergeItems provides for merging sorted items
   1. Sorted dict
    startIPStr   endIPStr
   ------------------------
   10.26.74.55 10.26.74.255
   10.23.77.88 10.23.77.240
   10.21.34.5  10.23.77.100
   10.12.14.2  10.12.14.50
   ------------------------
   2. Merged sorted dict
    startIPStr   endIPStr
   ------------------------
   10.26.74.55 10.26.74.255
   10.21.34.5  10.23.77.240
   10.12.14.2  10.12.14.50
   0.0.0.0     0.0.0.0
   ------------------------
*/
func (ipItems *IPItems) mergeItems() int {
    var mergedNum int

    items := ipItems.items
    length := len(items)

    for i := 0; i < length-1; i++ {

        if items[i].endIP == ZERO_IP {
            continue
        }

        for j := i + 1; j < length; j++ {
            if items[j].endIP == ZERO_IP {
                continue
            }

            mergedNum += ipItems.checkMerge(i, j)
        }
    }

    return mergedNum
}

/* InsertPair provides insert startIP,endIP into IpItems */
func (ipItems *IPItems) InsertPair(startIP, endIP net.IP) error {
    if bytes.Compare(startIP, endIP) == 1 {
        return fmt.Errorf("Insert(): err startIPStr %s > endIPStr %s",
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

    ipItems.items = append(ipItems.items, ipPair{sIP, eIP})
    return nil
}

/* InsertSingle single ip into ipitems */
func (ipItems *IPItems) InsertSingle(ip net.IP) error {
    return ipItems.ipSet.Add(ip)
}

/*
   Sort provides for sorting dict according startIP by descending order
   1. Origin dict
    startIPStr   endIPStr
   ------------------------
   10.26.74.55 10.26.74.255
   10.12.14.2  10.12.14.50
   10.21.34.5  10.23.77.100
   10.23.77.88 10.23.77.240
   ------------------------
   2. Sorted dict
    startIPStr   endIPStr
   ------------------------
   10.26.74.55 10.26.74.255
   10.23.77.88 10.23.77.240
   10.21.34.5  10.23.77.100
   10.12.14.2  10.12.14.50
   ------------------------
   3. Merged sorted dict
    startIPStr   endIPStr
   ------------------------
   10.26.74.55 10.26.74.255
   10.21.34.5  10.23.77.240
   10.12.14.2  10.12.14.50
   0.0.0.0     0.0.0.0
   ------------------------
   4. Dict after resliced
    startIPStr   endIPStr
   ------------------------
   10.26.74.55 10.26.74.255
   10.21.34.5  10.23.77.240
   10.12.14.2  10.12.14.50
   ------------------------
*/
func (ipItems *IPItems) Sort() {

    // Sort items according startIP by descending order
    sort.Sort(ipItems.items)

    // Merge item lines
    mergedNum := ipItems.mergeItems()
    length := len(ipItems.items) - mergedNum

    // Sort items according startIP by descending order
    sort.Sort(ipItems.items)

    // Reslice
    ipItems.items = ipItems.items[0:length]
}

/* get ip num of IPItems */
func (ipItems *IPItems) Length() int {
    num := len(ipItems.items)
    num += ipItems.ipSet.Len()

    return num
}
