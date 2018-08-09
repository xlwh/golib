/* bns_sync_table.go - a table for bns instances auto sync with bns */
/*
modification history
--------------------
2016/2/6, by Guang Yao, create
*/
/*
DESCRIPTION
*/
package bns_sync_table

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	"www.baidu.com/golang-lib/bns_util"
	"www.baidu.com/golang-lib/log"
	"www.baidu.com/golang-lib/module_state2"
)

const (
	BNS_ALL_INFO_ITEMS = 11 // the item count from "get_instance_by_service -a bnsname"
)

type BnsInstance struct {
	// all possible fields using "get_instance_by_service -a bnsname"
	Hostname    string            // hostname of instance
	IP          string            // [ip]
	Service     string            // [ServiceName]
	Port        int64             // [port]
	Status      int64             // [status]
	Load        int64             // [load]
	Offset      int64             // [offset]
	Tags        map[string]string // [tags]
	Extra       string            // [extra]
	ContainerId string            // [container-id]
	DeployPath  string            // [deploy-path]
}

// keys for recording state
type SyncStateParams struct {
	TotalCountKey     string // the key for total resolve count in State
	ErrCountKey       string // the key for err resolve count in State
	SucessCountKey    string // the key for sucess resolve count in State
	LastUpdateTimeKey string // the key for last update time in State
}

type BnsSyncTable struct {
	bnsName        string // bns service/service group name to sync
	updateInterval int    // interval to update instance info; in seconds
	storepath      string // file location to save/load bns instances info

	state *SyncState // for update resolve state

	lock    sync.RWMutex            // protect the following fields
	nameMap map[string]*BnsInstance // hostname => instance
	ipMap   map[string]*BnsInstance // host ip => instance
}

// return a BnsSyncTable
// Params:
// 		- bnsName: bnsName to sync
// 		- updateInterval: sync interval
// 		- storepath: file location to save/load bns instances info
// 		- state: for recoding state in sync; allow nil
// 		- stateParams: keys for recording state; allow nil if state is nil
//
// Returns:
// 		(table, err)
func NewBnsSyncTable(bnsName string, updateInterval int, storepath string,
	state *module_state2.State, stateParams *SyncStateParams) (*BnsSyncTable, error) {
	table := new(BnsSyncTable)

	table.bnsName = bnsName
	table.updateInterval = updateInterval
	table.storepath = storepath

	// init state
	syncState, err := NewSyncState(state, stateParams)
	if err != nil {
		return nil, fmt.Errorf("err in NewSyncState(): %v", err)
	}
	table.state = syncState

	// init instances in table
	err = table.init()
	if err != nil {
		return nil, fmt.Errorf("err in init(): %v", err)
	}

	// start sync
	go table.updateRoutine()

	return table, nil
}

// check whether an instance is in table by hostname
func (table *BnsSyncTable) HasHost(hostname string) bool {
	table.lock.RLock()
	defer table.lock.RUnlock()

	_, ok := table.nameMap[hostname]
	return ok
}

// check whether an instance is in table by ip
func (table *BnsSyncTable) HasIP(IP string) bool {
	table.lock.RLock()
	defer table.lock.RUnlock()

	_, ok := table.ipMap[IP]
	return ok
}

// get all instance
func (table *BnsSyncTable) GetAll() []*BnsInstance {
	ret := make([]*BnsInstance, 0)

	table.lock.RLock()
	defer table.lock.RUnlock()

	for _, instance := range table.nameMap {
		ret = append(ret, instance)
	}

	return ret
}

// init the table
func (table *BnsSyncTable) init() error {
	instances, err := resolveBns(table.bnsName)
	if err == nil {
		table.update(instances)
		// update to sucess count
		table.state.IncSucess()
		// set update time
		table.state.SetLastUpdate(time.Now().Format("20060102150405"))
		return nil
	}

	// err occurs; try to load from file
	table.state.IncErr()
	log.Logger.Warn("err in resolveBns(%s): %v", table.bnsName, err)
	instances, err = BnsInstancesLoad(table.storepath, table.bnsName)
	if err == nil {
		table.update(instances)
		return nil
	}

	// fail at last
	log.Logger.Error("err in BnsInstancesLoad(): %v", err)
	return err
}

// the routine to update the instances
func (table *BnsSyncTable) updateRoutine() {
	for {
		// wait
		time.Sleep(time.Duration(table.updateInterval) * time.Second)

		// resolve
		table.state.IncTotal()
		instances, err := resolveBns(table.bnsName)
		if err != nil {
			// update to err count and log
			table.state.IncErr()
			log.Logger.Error("err in resolveBns(%s):%v", table.bnsName, err)
		} else {
			// update to state
			table.state.IncSucess()
			table.state.SetLastUpdate(time.Now().Format("20060102150405"))

			// update instances
			table.update(instances)

			// save to file; only log in case of error
			err = BnsInstancesSave(table.storepath, table.bnsName, instances)
			if err != nil {
				log.Logger.Warn("err in BnsInstancesSave(%s, %s): %v", table.storepath,
					table.bnsName, err)
			}
		}
	}
}

// update instances
func (table *BnsSyncTable) update(instances []*BnsInstance) {
	// build hostname map and ip map
	nameMap := make(map[string]*BnsInstance)
	ipMap := make(map[string]*BnsInstance)
	for _, instance := range instances {
		nameMap[instance.Hostname] = instance
		ipMap[instance.IP] = instance
	}

	table.lock.Lock()
	defer table.lock.Unlock()

	// update to table
	table.nameMap = nameMap
	table.ipMap = ipMap
}

// resolve bns and return instances
func resolveBns(bnsName string) ([]*BnsInstance, error) {
	// get bns output
	outputLines, err := bns_util.GetBnsInstances(bnsName, []string{"-a"})
	if err != nil {
		return nil, fmt.Errorf("err in GetBnsInstances(%s): %v", bnsName, err)
	}

	// parse lines to build instances
	ret := make([]*BnsInstance, 0)
	for _, line := range outputLines {
		instance, err := parseBnsAllInfoLine(line)
		if err != nil {
			return nil, fmt.Errorf("err in parseBnsAllInfoLine(%s): %v", line, err)
		}
		ret = append(ret, instance)
	}

	return ret, nil
}

// parse a output line from "get_instance_by_service -a bnsname"
func parseBnsAllInfoLine(line string) (*BnsInstance, error) {
	// expected format:
	// hkg02-l-bfe03.hkg02 10.242.123.45 bfe-hk.BFE.hk 8900 0 -1 2 \
	// interface:clientip,keepalive:0,weight:10 extra container /tmp
	lineItems := strings.Split(line, " ")
	if len(lineItems) != BNS_ALL_INFO_ITEMS {
		return nil, fmt.Errorf("unexpected bns output format: %s, %d", line, len(lineItems))
	}

	instance := new(BnsInstance)

	instance.Hostname = lineItems[0]
	instance.IP = lineItems[1]
	instance.Service = lineItems[2]

	// parse port
	port, err := strconv.ParseInt(lineItems[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("err in parse port: %v", err)
	}
	instance.Port = port

	// parse Status
	status, err := strconv.ParseInt(lineItems[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("err in parse status: %v", err)
	}
	instance.Status = status

	// parse Load
	load, err := strconv.ParseInt(lineItems[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("err in parse load: %v", err)
	}
	instance.Load = load

	// parse Offset
	offset, err := strconv.ParseInt(lineItems[6], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("err in parse offset: %v", err)
	}
	instance.Offset = offset

	// parse tags
	tags := lineItems[7]
	if tags == "" {
		instance.Tags = make(map[string]string)
	} else {
		// parse tags to map
		tagMap, err := parseBnsTags(tags)
		if err != nil {
			return nil, fmt.Errorf("err in parse tags: %v", err)
		}
		instance.Tags = tagMap
	}

	instance.Extra = lineItems[8]
	instance.ContainerId = lineItems[9]
	instance.DeployPath = lineItems[10]

	return instance, nil
}

// parse tags to map
func parseBnsTags(tagStr string) (map[string]string, error) {
	// result to return
	ret := make(map[string]string)

	// expected format:
	// interface:clientip,keepalive:0,weight:10
	tagItems := strings.Split(tagStr, ",")
	for _, tagItem := range tagItems {
		keyValue := strings.Split(tagItem, ":")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("err in parse tagItem: %v", tagItem)
		}

		key := keyValue[0]
		value := keyValue[1]

		ret[key] = value
	}

	return ret, nil
}
