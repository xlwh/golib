/* instances_file.go - file for storing bns instance info */
/*
modification history
--------------------
2016/2/6, by Guang Yao, create based on bns_util/bns_info_file.go
*/
/*
DESCRIPTION
*/
package bns_sync_table

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// generate file path by service name and root path
func filePathGen(service, rootPath string) string {
	filePath := path.Join(rootPath, service)
	return filePath
}

/*
save bns instance info to file

Params:
    - rootPath: root path to store info file
    - service: service name, e.g., gslb-scheduler.BFE.all
    - instance: instances

Returns:
    error
*/
func BnsInstancesSave(rootPath string, service string, instances []*BnsInstance) error {
	// convert instances to json string
	jsonStr, err := json.Marshal(instances)
	if err != nil {
		return fmt.Errorf("json.Marshal():%s", err.Error())
	}

	// create root dir, if not exist
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		if os.MkdirAll(rootPath, 0777) != nil {
			return fmt.Errorf("os.MkdirAll(%s):%s", rootPath, err.Error())
		}
	}

	// save data to disk
	// generate full file path
	filePath := filePathGen(service, rootPath)

	// save to file
	err = ioutil.WriteFile(filePath, jsonStr, 0644)
	if err != nil {
		return fmt.Errorf("ioutil.WriteFile(%s):%s", filePath, err.Error())
	}

	return nil
}

/*
load bns instance info from file

Returns:
    (instances, error)
*/
func BnsInstancesLoad(rootPath string, service string) ([]*BnsInstance, error) {
	// generate full file path
	filePath := filePathGen(service, rootPath)

	// read all data from file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile(%s):%s", filePath, err.Error())
	}

	// json decode
	var instances []*BnsInstance
	err = json.Unmarshal(data, &instances)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal():%s", err.Error())
	}

	return instances, nil
}
