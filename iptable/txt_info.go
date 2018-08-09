/* txt_info.go - get file meta info */
/*
modification history
--------------------
2016/12/21, by Zhang Jiyang create
*/
/*
DESCRIPTION
    get the meta info(version, LineNum) from the firstline

example:

#{ "version":"1.2.3.4","Linenum": 2}

1.1.1.1 2.2.2.2 val1
3.3.3.3 4.4.4.4 val2
*/
package iptable

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrWrongLineNum = errors.New("metaInfo: LineNum must >= 0")
	ErrWrongVersion = errors.New("metaInfo: empty Version")
)

type MetaInfo struct {
	Version string
	LineNum int // lineNum
}

/*
* getFileInfo - get file meta info.
*
* get file meta info from first line, if failed, get actual lineNum
*
* PARAMS:
*   - path: path of file
*
* RETURNS:
*   - (*MetaInfo, nil), if success ,return file metaInfo
*   - (nil error), if failed
 */
func getFileInfo(path string) (*MetaInfo, error) {
	// get meta info from comment(first line)
	if metaInfo, err := getCommentFileInfo(path); err == nil {
		return metaInfo, nil
	}

	// cal actual fileInfo
	return getActualFileInfo(path)
}

/*
* getCommentFileInfo - read the first Line, decode the json string, and return
*
* eg. #{ "version":"1.2.3.4","LineNum": 1234}
*
* PARAMS:
*   - path: path of file
*
* RETURNS:
*   - (*MetaInfo, nil), if success ,return file metaInfo
*   - (nil error), if failed
 */
func getCommentFileInfo(path string) (*MetaInfo, error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read the first line
	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	// get meta data
	firstLine := strings.Trim(string(line), " \t\r\n")
	if !strings.HasPrefix(firstLine, "#") {
		return nil, fmt.Errorf("firstLine don't contail meta info")
	}

	// decode the json string
	metaInfo := &MetaInfo{}
	metaString := strings.Trim(firstLine, "#")
	decoder := json.NewDecoder(strings.NewReader(metaString))
	if err = decoder.Decode(metaInfo); err != nil {
		return nil, err
	}

	// check metaInfo
	if err = checkMetaInfo(*metaInfo); err != nil {
		return nil, err
	}

	return metaInfo, nil
}

// getActualFileInfo: cal meta info from file
func getActualFileInfo(path string) (*MetaInfo, error) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open(): %s %s", path, err.Error())
	}
	defer file.Close()

	lineCounter := 0
	// scan the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Remove all leading and trailing spaces and tabs
		line := strings.Trim(scanner.Text(), " \t")
		//Line begins with "#" is considered as a comment
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		// Check format
		_, _, _, err = checkSplit(line, " ")
		if err != nil {
			return nil, fmt.Errorf("checkLine(): line[%s] err[%s]", line, err.Error())
		}

		lineCounter += 1
	}
	err = scanner.Err()
	// Scan meets error
	if err != nil {
		return nil, fmt.Errorf("scan file: err, %s", err.Error())
	}

	return &MetaInfo{
		Version: "",
		LineNum: lineCounter,
	}, nil
}

// check meta info
func checkMetaInfo(Info MetaInfo) error {
	if Info.Version == "" {
		return ErrWrongVersion
	}

	// lineNum must >= 0
	if Info.LineNum < 0 {
		return ErrWrongLineNum
	}

	return nil
}
