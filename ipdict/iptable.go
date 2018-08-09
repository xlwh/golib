/* iptable.go - Provides for thread-safe IP table operation

   see txt_load.go for more detainfo.
/*
modification history
--------------------
2014/7/7, by Li Bingyi, create
2014/9/11, by zhangjiyang01@baidu.com, modify
    - use hash table for single IP
*/
/*
DESCRIPTION

See txt_load.go for more detail usage info.

Usage:

    // See txt_load.go for more detail info.
    import (
        "net"
    )

    import (
        import "www.baidu.com/golang-lib/ipdict"
    )

    var hit   bool
    var items *ipdict.IPItems
    var srcIP  net.IP

    // Create table
    table = ipdict.NewIPTable()

    // Switch items
    table.Update(items)

    // Search items
    hit = table.Search(srcIP)
*/

package ipdict

import (
	"net"
	"sort"
	"sync"

	"www.baidu.com/golang-lib/net_util"
)

type IPTable struct {
	lock    sync.Mutex
	ipItems *IPItems
}

func NewIPTable() *IPTable {
	table := new(IPTable)
	return table
}

func (t *IPTable) Version() string {
	t.lock.Lock()
	ipItems := t.ipItems
	t.lock.Unlock()

	if ipItems != nil {
		return ipItems.Version
	}
	return ""
}

/* Update provides for thread-safe switching items */
func (t *IPTable) Update(items *IPItems) {
	t.lock.Lock()
	t.ipItems = items
	t.lock.Unlock()
}

/* Search provides for binary search IP in dict */
func (t *IPTable) Search(srcIP net.IP) bool {
	var hit bool
	t.lock.Lock()
	ipItems := t.ipItems
	t.lock.Unlock()

	// check ipItems
	if ipItems == nil {
		return false
	}
	// convert ip to ipv4
	ip := srcIP.To4()

	// 1. check at the ip set
	if ipItems.ipSet.Exist(ip) {
		return true
	}

	// convert ip to uint32 format
	ipUint, err := net_util.IPv4ToUint32(ip)
	if err != nil {
		return false
	}

	// 2. check at the item array
	items := ipItems.items
	itemsLen := len(items)

	i := sort.Search(itemsLen,
		func(i int) bool { return items[i].startIP <= ipUint })

	if i < itemsLen {
		if items[i].endIP >= ipUint {
			hit = true
		}
	}

	return hit
}
