/* bfe_backends_linux.go - get bfe backends info from bns instance */
/*
modification history
--------------------
2016/1/28, by Taochunhua, create
2017/12/01, by yuxiaofei, modify
    - add GetBnsInstanceWithDefaultWeight func
*/
/*
DESCRIPTION
*/
package bns_util

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

// backend info in each bns instance
// note: all members are pointers. It is compatible with go-bfe's BackendConf struct.
// see go-bfe package "bfe_config/bfe_cluster_conf/cluster_table_conf"
type BfeBackendConf struct {
	Name   *string // e.g., "tc-bae-static05.tc"
	Addr   *string // e.g., "10.26.35.33"
	Port   *int    // e.g., 8000
	Weight *int    // weight in load balance, e.g., 10
}

type BfeBackendConfList []*BfeBackendConf

// get a subcluster's backend conf.
// by calling and parse output of get_instance_by_service cmd
func GetBnsInstanceBfeBackends(subClusterBns string) (BfeBackendConfList, error) {
	var conf BfeBackendConfList

	out, err := exec.Command("get_instance_by_service", "-a", subClusterBns).Output()
	// get_instance_by_service return 0 on success
	// return 240 if any instance's status is not 0
	if err != nil && err.Error() != "exit status 240" {
		return conf, err
	}

	withDefaultWeight := false
	if err := parseBfeBackends(out, &conf, withDefaultWeight); err != nil {
		return conf, err
	}

	return conf, nil
}

// get a subcluster's backend conf.
// by calling and parse output of get_instance_by_service cmd
// if bns's backend has no weight, set default weight to 1
func GetBnsInstanceWithDefaultWeight(subClusterBns string) (BfeBackendConfList, error) {
	var conf BfeBackendConfList

	out, err := exec.Command("get_instance_by_service", "-a", subClusterBns).Output()
	// get_instance_by_service return 0 on success
	// return 240 if any instance's status is not 0
	if err != nil && err.Error() != "exit status 240" {
		return conf, err
	}

	withDefaultWeight := true
	if err := parseBfeBackends(out, &conf, withDefaultWeight); err != nil {
		return conf, err
	}

	return conf, nil
}

// parse get_instance_by_service output to BfeBackendConfList
// output format example:
// hostname ip bns port status unknown unknown unknown tag
// st01-orp-router03.st01 10.52.63.39 router.orp.tc 8080 0 -1 11 \
//      interface:clientip,keepalive:0,weight:10
// m1-orp-router01.m1 10.42.222.48 router.orp.tc 8080 0 -1 32
// NOTE: bns's output format may change?
func parseBfeBackends(out []byte, conf *BfeBackendConfList, withDefaultWeight bool) error {
	validFieldLen := 8
	if withDefaultWeight {
		validFieldLen = 7
	}

	instances := bytes.Split(out, []byte("\n"))
	for _, instance := range instances {
		fields := bytes.Split(instance, []byte(" "))
		fieldLen := len(fields)
		if fieldLen < validFieldLen {
			continue
		}

		name := string(fields[0])
		addr := string(fields[1])

		// port
		port, err := strconv.Atoi(string(fields[3]))
		if err != nil {
			port = 0
		}

		// status
		status, err := strconv.Atoi(string(fields[4]))
		if err != nil {
			status = -2
		}

		// weight
		weight := -1
		if fieldLen > 7 {
			// length > 7, may contain weight tag, try to get weight
			weight, err = parseBfeBackendWeight(fields[7])
			if err != nil {
				// parse backend weight failed, set weight to -1
				weight = -1
			}
		}
		// weight < 0 means getting weight failed
		if withDefaultWeight && weight < 0 {
			weight = 1
		}

		if len(name) == 0 || len(addr) == 0 || port <= 0 || weight <= 0 || status != 0 {
			continue
		}

		backendConf := &BfeBackendConf{
			Name:   &name,
			Addr:   &addr,
			Port:   &port,
			Weight: &weight,
		}

		*conf = append(*conf, backendConf)
	}

	if len(*conf) == 0 {
		return fmt.Errorf("get 0 avail instance from output")
	}

	return nil
}

// tags is of format k1:v1,k2:v2,...
// get value of key "weight"
// if weight not found, return error
func parseBfeBackendWeight(tags []byte) (int, error) {
	// if it has no weight, set weight to 0
	weight := 0
	got := false

	tagslice := bytes.Split(tags, []byte(","))
	for _, tag := range tagslice {
		n, _ := fmt.Sscanf(string(tag), "weight:%d", &weight)
		if n > 0 {
			got = true
			break
		}
	}

	if !got {
		return weight, errors.New("weight not found")
	}

	return weight, nil
}
