/* txt_load.go - thread-safe loading IP txt file to IP items in memory */

/*
modification history
--------------------
2016/12/21, by Zhang JiYang Create
*/
/*
DESCRIPTION

    Provides for thread-safely loading IP txt file to IP items in memory

    File format:
    ------------------------
    start_ip end_ip val1
    start_ip end_ip val2
    # this is a comment line
    start_ip end_ip val3
    ...        ...
    ------------------------

    File line is composed by start_ip, end_ip and value seprated by space[s]
    Line begins with '#' is recognized as a comment line
    Empty line is ignored too.
    Spaces and tabs are allowed in the leading and trailing of the line

    File example:
    ------------------------
    10.26.74.55 10.26.74.255 CN
    10.12.14.2  10.12.14.50  US

    10.21.34.5  10.23.77.100 EU
    # this is a comment line
    10.23.77.88 10.23.77.240 SG
    ------------------------

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
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

import (
	"www.baidu.com/golang-lib/log"
	"www.baidu.com/golang-lib/net_util"
)

// default max value length
var defaultValLen = 32

var (
	// file version not change, needn't load the file
	ERR_NO_NEED_UPDATE = errors.New("Version no change no need update")
	// wrong meta info
	ERR_WRONG_META_INFO = errors.New("Wrong meta info")
)

type TxtFileLoader struct {
	filePath  string
	maxLine   int // max line num
	maxValLen int // max value length of every value

	checkValFunc func(string) error
}

func NewTxtFileLoader(path string) *TxtFileLoader {
	f := new(TxtFileLoader)
	f.filePath = path
	f.maxLine = -1
	f.maxValLen = defaultValLen
	return f
}

// set max line num
func (f *TxtFileLoader) SetMaxLine(maxLine int) {
	f.maxLine = maxLine
}

// set check value function
func (f *TxtFileLoader) SetCheckValFunc(checkFunc func(string) error) {
	f.checkValFunc = checkFunc
}

// set max value length
func (f *TxtFileLoader) SetMaxValLen(maxValLen int) {
	f.maxValLen = maxValLen
}

// checkSplit: checks line split format
//   legal start ip and end ip is seprated by seps
// Params:
//   - line: line for split
//   - sep: separator for line
// Returns:
//   - (startIP, endIP, value, nil): if success
//   - (nil, nil, "", error); if failed
func checkSplit(line string, sep string) (net.IP, net.IP, string, error) {
	var startIPStr, endIPStr, value string
	var startIP, endIP net.IP

	segments := strings.SplitN(line, sep, 3)

	// segments[0] : start ip(string)
	// segments[1] : end ip(string)
	// segments[2] : value
	if len(segments) == 3 {
		startIPStr = strings.Trim(segments[0], " \t")
		endIPStr = strings.Trim(segments[1], " \t")
		value = strings.Trim(segments[2], " \t")
	} else {
		return nil, nil, value, fmt.Errorf("expect 3 segments, got %d", len(segments))
	}

	// startIPStr format err
	if startIP = net_util.ParseIPv4(startIPStr); startIP == nil {
		return nil, nil, value, fmt.Errorf("startIP, worng format %s", startIPStr)
	}

	// endIPStr format err
	if endIP = net_util.ParseIPv4(endIPStr); endIP == nil {
		return nil, nil, value, fmt.Errorf("endIP, wrong format %s", endIPStr)
	}

	return startIP, endIP, value, nil
}

/* check Version num and load IP txt file to IP items in memory */
func (f TxtFileLoader) CheckAndLoad(curVersion string) (*IPDict, string, error) {
	var startIP, endIP net.IP
	var value string

	path := f.filePath
	// get file Version and lineNum
	metaInfo, err := getFileInfo(path)
	if err != nil {
		return nil, "", fmt.Errorf("loadFile(): %s %s", path, err.Error())
	}
	newVersion := metaInfo.Version
	lineNum := metaInfo.LineNum

	// if newVersion equal curVersion, and newVersion is not empty no need to update
	if newVersion == curVersion && newVersion != "" {
		return nil, newVersion, ERR_NO_NEED_UPDATE
	}

	// if f.maxLine is set, check lineNum
	if f.maxLine != -1 && lineNum > f.maxLine {
		log.Logger.Info("file %s lineNum[%d] large than f.maxLine[%d]", path, lineNum, f.maxLine)
		lineNum = f.maxLine
	}

	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("loadFile(): %s, %s", path, err.Error())
	}
	defer file.Close()

	// create ipDict
	ipDict, err := NewIPDict(lineNum, f.maxValLen)
	if err != nil {
		return nil, "", err
	}

	// counter line num
	lineCounter := 0
	// scan the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCounter++
		// Remove all leading and trailing spaces and tabs
		line := strings.Trim(scanner.Text(), " \t")
		//Line begins with "#" is considered as a comment
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		// Check and split line
		startIP, endIP, value, err = checkSplit(line, " ")
		if err != nil {
			log.Logger.Warn("checkSplit(): line: %s lineNum: %d, err: %s", line, lineCounter, err.Error())
			continue
		}

		// check value format
		if fn := f.checkValFunc; fn != nil {
			if err := fn(value); err != nil {
				log.Logger.Warn("illegal value, line: %d val: %s err: %s", lineCounter, value, err.Error())
				continue
			}
		}

		// insert start ip and end ip into dict
		if err = ipDict.Add(startIP, endIP, value); err != nil {
			log.Logger.Warn("ipDict.Add():  [%s-%s %s] err: %s", startIP, endIP, value, err.Error())
			continue
		}
	}

	// Scan meets error
	if err = scanner.Err(); err != nil {
		return nil, newVersion, fmt.Errorf("loadFile(): err %s", err.Error())
	}

	// Load succ, sort dict
	if err = ipDict.SortAndCheck(); err != nil {
		return ipDict, newVersion, err
	}
	return ipDict, newVersion, nil
}
