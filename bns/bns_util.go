/* bns_util.go - utility functions for bns client */
/*
modification history
--------------------
2015/5/28, by Sijie Yang, create

2017/5/12, by Sijie Yang, modify
    - support local name conf
*/
/*
DESCRIPTION
*/
package bns

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

import (
	"code.google.com/p/goprotobuf/proto"
)

const (
	DefaultWeight = 1 // defualt weight for instance
)

// instance info
type Instance struct {
	Host   string // instance host
	Port   int    // instance port
	Weight int    // instance weight
}

// local name conf for service
type LocalNameConf map[string][]Instance

// local name conf file
type NameConf struct {
	Version string
	Config  LocalNameConf
}

var localNameConf LocalNameConf
var localNameLock sync.RWMutex

// load name conf file
func LoadLocalNameConf(filename string) error {
	// load local name conf
	nameConf, err := parseLocalNameConf(filename)
	if err != nil {
		return err
	}

	// update local name map
	localNameLock.Lock()
	localNameConf = nameConf
	localNameLock.Unlock()
	return nil
}

/* GetInstances - get instance addr and weight info
 * Params:
 *     - bnsClient: bns client instance
 *     - serviceName: service name
 *
 * Return:
 *     - instanceList: server instance list
 *     - error: error if fail
 */
func GetInstances(bnsClient *Client, serviceName string) ([]Instance, error) {
	// check local conf
	if instances, ok := getInstancesLocal(serviceName); ok {
		return instances, nil
	}

	// check bns service
	return getInstancesRemote(bnsClient, serviceName)
}

// get address for service (deprecated)
func GetAddr(bnsClient *Client, serviceName string) ([]string, error) {
	instances, err := GetInstances(bnsClient, serviceName)
	if err != nil {
		return nil, err
	}

	addrList := make([]string, 0)
	for _, instance := range instances {
		address := fmt.Sprintf("%s:%d", instance.Host, instance.Port)
		addrList = append(addrList, address)
	}

	return addrList, nil
}

// get address and weight for service (deprecated)
func GetAddrAndWeight(bnsClient *Client, serviceName string) (map[string]int32, error) {
	instances, err := GetInstances(bnsClient, serviceName)
	if err != nil {
		return nil, err
	}

	addrInfo := make(map[string]int32)
	for _, instance := range instances {
		address := fmt.Sprintf("%s:%d", instance.Host, instance.Port)
		addrInfo[address] = int32(instance.Weight)
	}

	return addrInfo, nil
}

// get instances from local name conf
func getInstancesLocal(serviceName string) ([]Instance, bool) {
	localNameLock.RLock()
	instances, ok := localNameConf[serviceName]
	localNameLock.RUnlock()

	return instances, ok
}

func parseLocalNameConf(filename string) (LocalNameConf, error) {
	var conf NameConf
	var err error

	// open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode the file
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&conf); err != nil {
		return nil, err
	}

	// check config
	err = checkLocalNameConf(conf)
	if err != nil {
		return nil, err
	}

	return conf.Config, nil
}

func checkLocalNameConf(conf NameConf) error {
	for name, instances := range conf.Config {
		for _, instance := range instances {
			if err := checkInstance(instance); err != nil {
				return fmt.Errorf("invalid instance for %s: %s", name, err)
			}
		}
	}
	return nil
}

func checkInstance(instance Instance) error {
	if len(instance.Host) == 0 {
		return fmt.Errorf("invalid host: %v", instance)
	}
	if instance.Port < 0 || instance.Port > 65535 {
		return fmt.Errorf("invalid port: %v", instance)
	}
	if instance.Weight < 0 {
		return fmt.Errorf("invalid weight: %v", instance)
	}
	return nil
}

// get instances from remote bns servcie
func getInstancesRemote(bnsClient *Client, serviceName string) ([]Instance, error) {
	instances := make([]Instance, 0)

	// get name response
	resp, err := getNameResp(bnsClient, serviceName)
	if err != nil {
		return nil, err
	}
	if resp.GetRetcode() < 0 {
		return nil, fmt.Errorf("LocalNamingResponse Retcode %d", resp.GetRetcode())
	}

	// parse name resposne
	for _, instanceInfo := range resp.InstanceInfo {
		instance, err := parseInstance(instanceInfo)
		if err != nil {
			continue
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

// get name response from bns service
func getNameResp(bnsClient *Client, serviceName string) (*LocalNamingResponse, error) {
	request := &LocalNamingRequest{ServiceName: proto.String(serviceName)}
	resp := new(LocalNamingResponse)
	err := bnsClient.Call(request, resp)
	return resp, err
}

// parse instance info
func parseInstance(info *InstanceInfo) (Instance, error) {
	var instance Instance

	// check instance status
	status := info.GetInstanceStatus()
	if status == nil || status.GetStatus() != 0 {
		return instance, fmt.Errorf("invalid status")
	}

	// parse instance info
	instance.Host = parseIp(info.GetHostIp())
	instance.Port = int(status.GetPort())
	instance.Weight = parseWeight(status.GetTags(), DefaultWeight)

	// check instance info
	if err := checkInstance(instance); err != nil {
		return instance, err
	}

	return instance, nil
}

// parse weight value from tags
func parseWeight(tags string, weight int) int {
	// Note: tags format: key:vaule,key2:value2, ... , keyN:valueN
	tagslice := strings.Split(tags, ",")
	for _, tag := range tagslice {
		if _, err := fmt.Sscanf(tag, "weight:%d", &weight); err == nil {
			return weight
		}
	}
	return weight
}

// parse Ip
func parseIp(ip uint32) string {
	netIp := net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
	return netIp.String()
}
