/* bns_instance.go - get instance by service name	*/
/*
modification history
--------------------
2015/7/29, by Zhang Miao, create
2015/8/3, by Zhang Miao, move to golang-lib
2015/8/3, by Guang Yao, add options support for bns cmd
2015/9/1, by Guang Yao, add GetBnsInstancesIP
*/
/*
DESCRIPTION
get instance by bns service name
*/
package bns_util

import (
    "bytes"
    "errors"
    "net"
    "os/exec"
    "strings"
)

/*
get instance by bns service name

Params:
    - service: name of service
    - options: options for get_instance_by_service, e.g., ["-i", "-a"]

Returns:
    (instances, error)
*/
func GetBnsInstances(service string, options []string) ([]string, error) {
    // assemble params
    params := append(options, service)
    
    // get instance with cmd "get_instance_by_service"
    out, err := exec.Command("get_instance_by_service", params...).Output()
    if err != nil {
        return nil, err
    }

    // split instance names
	byteInstances := bytes.Split(out, []byte{'\n'})

    instances := make([]string, 0)
    for _, instance := range byteInstances {
        if len(instance) != 0 {
            instances = append(instances, string(instance))
        }
    }

    return instances, nil
}

/*
get instances' ip addresses by bns service name

Params:
    - service: name of service
Returns:
    (ips, error)
*/
func GetBnsInstancesIP(service string) ([]string, error) {
    // get instances from bns
    instances, err := GetBnsInstances(service, []string{"-i"})
    if err != nil {
        return nil, err
    }

    // parse ips from instances
    ips, err := parseInstancesIP(instances)
    if err != nil {
        return nil, err
    }
    
    return ips, nil
}

/*
parse ip address from bns output line

Params:
    - instances: output lines of bns cmd
Returns:
    (ips, error)
*/
func parseInstancesIP(instances []string) ([]string, error) {
    ips := make([]string, len(instances))

    // parse each line
    for index, instance := range instances {
        // check output format
        instanceSegs := strings.Split(instance, " ")
        if len(instanceSegs) != 2 {
            // output is not as expected
            return nil, errors.New("unexpected bns output: " + instance)
        }

        // check IP format
        ip := instanceSegs[1]
        if net.ParseIP(ip) == nil {
            return nil, errors.New("invalid IP format:" + ip)
        }

        // put ip in the table
        ips[index] = ip
    }

    return ips, nil
}
