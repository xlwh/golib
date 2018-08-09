/* iptable.go - Provides for thread-safe IP dict operation */
/*
modification history
--------------------
2016/12/19, by Zhang Jiyang, create
*/
/*
DESCRIPTION

Usage: 
    import "www.baidu.com/golang-lib/iptable"
    
    var fileLoader ipdict.TxtFileLoader
    var err error
    var ipdict *iptable.IPDict

    table := iptable.NewIPTable()

    fileLoader := ipTable.NewTxtFileLoader(filePath)

    ipDict, err = fileLoader.CheckAndLoad("")

    if err != nil {
        // Err handler
    } else {
        // Update items
        table.Update(ipDict)
    }
*/

package iptable

import (
	"net"
	"sync"
)

// IPTable Provide thread-safe ipdict
type IPTable struct {
	ipDict  *IPDict // ipdict
	version string  // version
	lock    sync.RWMutex
}

// NewIPTable create new iptable instance
func NewIPTable() *IPTable {
	t := new(IPTable)
	return t
}

// Update
func (t *IPTable) Update(ipDict *IPDict, version string) {
	t.lock.Lock()
	t.ipDict = ipDict
	t.version = version
	t.lock.Unlock()
}

// Search
func (t *IPTable) Search(ip net.IP) (string, bool) {
	t.lock.RLock()
	ipDict := t.ipDict
	t.lock.RUnlock()

	return ipDict.Search(ip)
}

// Version, get version
func (t *IPTable) Version() string {
	t.lock.RLock()
	version := t.version
	t.lock.RUnlock()

	return version
}
