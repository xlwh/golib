/* dict_ip_location.go - ip location information for dict-server */
/*
modification history
--------------------
2016/6/6, by Jiang Hui, create
*/
/*
DESCRIPTION

*/

package loc_load

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

import (
	"www.baidu.com/golang-lib/ipdict"
	"www.baidu.com/golang-lib/net_util"
)

const (
	SEGMENT_LEN       = 7
	MAX_LOCATION_LINE = 1000000
	MAX_LOCATION_LEN  = 48
)

type IpLocDictFile struct {
	fileName string
	version  string
	maxLine  uint32
}

type ipSectionStr struct {
	startIp  net.IP
	endIp    net.IP
	location string
}

//maxline should in section[1,MAX_LOCATION_LINE]
func NewIpLocDictFile(fileName string, maxLine uint32) (*IpLocDictFile, error) {
	if maxLine == 0 || maxLine > MAX_LOCATION_LINE {
		return nil, fmt.Errorf("NewIpLocDictFile() err: maxline(%d) is invalid", maxLine)
	}

	txtFile := new(IpLocDictFile)
	txtFile.fileName = fileName
	txtFile.version = ""
	txtFile.maxLine = maxLine

	return txtFile, nil
}

// get version of dictfile
func (t *IpLocDictFile) getVersion() string {
	return t.version
}

//file begin with #
//the version contain in the line that begin with #
//example #version: v2.0.1
func (t *IpLocDictFile) parseVersion() (string, error) {
	file, err := os.Open(t.fileName)
	if err != nil {
		return "", fmt.Errorf("parseVersion(): %s, %s", t.fileName, err.Error())
	}
	defer file.Close()

	// scan the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Remove all leading and trailing spaces and tabs
		line := strings.Trim(scanner.Text(), " \t")
		//the version contain in the Line that begins with "#"
		if !strings.HasPrefix(line, "#") {
			return "", fmt.Errorf("parseVersion() err: file (%s) does not exist version", t.fileName)
		}

		//line like #version: v2.0.1
		if strings.Contains(line, "version:") {
			cols := strings.Split(line, ":")
			if len(cols) != 2 {
				return "", fmt.Errorf("parseVersion() err: should contain 2 fields seperated by : line (%s)", line)
			}
			version := strings.Trim(cols[1], " \t")
			return version, nil
		}
	}

	// Scan meets error
	err = scanner.Err()
	if err != nil {
		return "", fmt.Errorf("parseVersion(): err, %s", err.Error())
	}
	return "", fmt.Errorf("parseVersion() error")
}

//parse line format,
//like "1.0.2.0|1.0.3.255|CN|CHINANET|FUJIAN|None|None|83|90|85|0|0"
func checkSplit(line string, sep string) (ipSectionStr, error) {
	var startIPStr, endIPStr string
	var startIP, endIP net.IP
	var ipLocation string
	segments := strings.SplitN(line, sep, SEGMENT_LEN)

	// get start ip ,end ip, country code
	if len(segments) == SEGMENT_LEN {
		startIPStr = strings.Trim(segments[0], " \t")
		endIPStr = strings.Trim(segments[1], " \t")
		ipLocation = strings.Trim(segments[2], " \t")
		ipLocation += ":" + strings.Trim(segments[4], " \t")
		ipLocation += ":" + strings.Trim(segments[5], " \t")
	} else {
		return ipSectionStr{}, fmt.Errorf("checkSplit() err: should contain 7 fields seperated by | (%s)", line)
	}

	// startIPStr format err
	if startIP = net_util.ParseIPv4(startIPStr); startIP == nil {
		return ipSectionStr{}, fmt.Errorf("checkSplit() err: line (%s) format", line)
	}

	// endIPStr format err
	if endIP = net_util.ParseIPv4(endIPStr); endIP == nil {
		return ipSectionStr{}, fmt.Errorf("checkSplit() err: line (%s) format", line)
	}

	return ipSectionStr{startIP, endIP, ipLocation}, nil
}

//load  when only first load or version mismatch
//every unit in iplocationitem is (startip,endip,location)
//assume dict file have been sorted
func (t *IpLocDictFile) CheckAndLoad(curVersion string) (*ipdict.IpLocationTable, error) {
	var ipSec ipSectionStr
	var lineCounter uint32
	var err error

	// get file Version and lineNum
	t.version, err = t.parseVersion()
	if err != nil {
		return nil, fmt.Errorf("checkAndLoad() err: from file (%s) parseVersion error(%s)", t.fileName, err.Error())
	}

	// check version
	if t.version != "" && (t.version == curVersion) {
	    return nil, ErrNoNeedUpdate
	}

	// open file
	file, err := os.Open(t.fileName)
	if err != nil {
		return nil, fmt.Errorf("checkAndLoad() err:open (%s) error (%s)", t.fileName, err.Error())
	}
	defer file.Close()

	// create location table
	ipLocTable, err := ipdict.NewIpLocationTable(t.maxLine, MAX_LOCATION_LEN)
	if err != nil {
		return nil, fmt.Errorf("checkAndLoad() err: NewIpLocationTable error (%s)", err.Error())
	}

	ipLocTable.Version = t.version

	// scan the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Remove all leading and trailing spaces and tabs
		line := strings.Trim(scanner.Text(), " \t")
		//Line begins with "#" is considered as a comment
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		// Check line format
		ipSec, err = checkSplit(line, "|")
		if err != nil {
			return nil, fmt.Errorf("checkAndLoad() err: checkSplit (%s)", err.Error())
		}

		// check if lineCounter > maxLine or not
		lineCounter += 1
		if lineCounter > t.maxLine {
			return ipLocTable, ErrMaxLineExceed
		}

		// insert start ip, end ip, location into dict table
		err = ipLocTable.Add(ipSec.startIp, ipSec.endIp, ipSec.location)
		if err != nil {
			return nil, fmt.Errorf("checkAndLoad() err:add IpLocationTable error (%s)", err.Error())
		}
	}

	// Scan meets error
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("checkAndLoad() err:scan file error (%s)", err.Error())
	}

	return ipLocTable, nil
}
