/* txt_info.go - get file meta info */
/*
modification history
--------------------
2014/9/24, by zhangjiyang01@baidu.com, create
2014/10/08, by zhangjiyang01@baidu.com, modify
               - split lineNum to singleIPNum and pairIPNum
*/
/*
DESCRIPTION
    get the meta info(version, singleIPNum, pairIPNum) from the firstline

example:
#{ "version":"1.2.3.4","singleIPNum": 1, "pairIPNum":1 }
 1.1.1.1 2.2.2.2
 3.3.3.3
*/
package txt_load

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "net"
    "os"
    "strings"
)

type MetaInfo struct {
    Version     string
    SingleIPNum int // single IP num
    PairIPNum   int // Pair IP num
}

/*
* getFileInfo - get file meta info.
*
* get file meta info from first line, if failed, get actual IPNums
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

    return getActualFileInfo(path)
}

/*
* getCommentFileInfo - read the first Line, decode the json string, and return
*
* eg. #{ "version":"1.2.3.4","singleIPNum": 1234, "pairIPNum": 1234}
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
    metaInfo := &MetaInfo{SingleIPNum: -1, PairIPNum: -1}
    metaString := strings.Trim(firstLine, "#")
    decoder := json.NewDecoder(strings.NewReader(metaString))
    err = decoder.Decode(metaInfo)
    if err != nil {
        return nil, err
    }

    // check metaInfo
    err = checkMetaInfo(*metaInfo)
    if err != nil {
        return nil, err
    }

    return metaInfo, nil
}

// getActualFileInfo: cal meta info from file
func getActualFileInfo(path string) (*MetaInfo, error) {
    var startIP, endIP net.IP

    // open file
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("open(): %s, %s", path, err.Error())
    }
    defer file.Close()

    singleIPCounter := 0
    pairIPCounter := 0
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
        startIP, endIP, err = checkLine(line)
        if err != nil {
            return nil, fmt.Errorf("checkLine(): line[%s] err[%s]", line, err.Error())
        }

        // insert start ip and end ip into dict
        if bytes.Compare(startIP, endIP) == 0 {
            singleIPCounter += 1
        } else {
            pairIPCounter += 1
        }
    }
    err = scanner.Err()
    // Scan meets error
    if err != nil {
        return nil, fmt.Errorf("scan file: err, %s", err.Error())
    }

    return &MetaInfo{
        Version:     "",
        SingleIPNum: singleIPCounter,
        PairIPNum:   pairIPCounter,
    }, nil
}

/* check meta info */
func checkMetaInfo(Info MetaInfo) error {
    if Info.Version == "" {
        return fmt.Errorf("metaInfo:Version is empty string")
    }

    /* PairIPNum/SingleIPNum must >= 0 */
    if Info.PairIPNum < 0 || Info.SingleIPNum < 0 {
        return fmt.Errorf("metaInfo:PairIPNum || SingleIPNum < 0")
    }

    return nil
}
