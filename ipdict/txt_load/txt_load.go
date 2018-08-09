/* txt_load.go - thread-safe loading IP txt file to IP items in memory */

/*
modification history
--------------------
2014/7/9, by Li Bingyi, create
2014/9/11, by zhangjiyang01@baidu.com, modify
    - use hash table for single IP
*/
/*
DESCRIPTION

    Provides for thread-safely loading IP txt file to IP items in memory

    File format:
    ------------------------
    start_ip [end_ip]
    start_ip [end_ip]
    # this is a comment line
    start_ip [end_ip]
    ...        ...
    ------------------------

    File line is composed by start_ip and end_ip, seprated by space[s] or tab[s], end_ip is optinal
    If there is no end_ip, considered start_ip instead.
    start_ip must >= end_ip by alphabetical order
    Line begins with '#' is recognized as a comment line
    Empty line is ignored too.
    Spaces and tabs are allowed in the leading and trailing of the line

    File example:
    ------------------------
    10.26.74.55 10.26.74.255
    10.12.14.2  10.12.14.50

    10.21.34.5  10.23.77.100
    # this is a comment line
    10.23.77.88 10.23.77.240
    ------------------------

Usage:

    import "www.baidu.com/golang-lib/ipdict"

    var fileLoader ipdict.TxtFileLoader
    var err error
    var ipItems *ipdict.IPItems

    table := ipdict.NewIPTable()

    fileLoader := ipdict.NewTxtFileLoader(fileName)

    ipItems, err = fileLoader.CheckAndLoad("")

    if err != nil {
        // Err handler
    } else {
        // Switch items
        t.Update(ipItems)
    }
*/

package txt_load

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "net"
    "os"
    "strings"
)

import (
    "www.baidu.com/golang-lib/ipdict"
    "www.baidu.com/golang-lib/net_util"
)

var (
    // file version not change, needn't load the file
    ERR_NO_NEED_UPDATE = errors.New("Version no change no need update")
    // line num of file larger than maxline configured
    ERR_MAX_LINE_EXCEED = errors.New("Max line exceed")
    // wrong meta info
    ERR_WRONG_META_INFO = errors.New("Wrong meta info")
)

type TxtFileLoader struct {
    fileName string
    maxLine  int
}

func NewTxtFileLoader(fileName string) *TxtFileLoader {
    f := new(TxtFileLoader)
    f.fileName = fileName
    f.maxLine = -1
    return f
}

// set max line num
func (f *TxtFileLoader) SetMaxLine(maxLine int) {
    f.maxLine = maxLine
}

/*
   checkSplit checks line split format
   legal start ip and end ip is seprated by space[s]/tab[s]
*/
func checkSplit(line string, sep string) (net.IP, net.IP, error) {
    var startIPStr, endIPStr string
    var startIP, endIP net.IP

    segments := strings.SplitN(line, sep, 2)
    segLen := len(segments)

    // Segments[0] : start ip string
    // Segments[1] : end ip string(start ip string instead when no end ip string found)
    if segLen == 1 {
        startIPStr, endIPStr = segments[0], segments[0]
    } else if len(segments) == 2 {
        startIPStr = strings.Trim(segments[0], " \t")
        endIPStr = strings.Trim(segments[1], " \t")
    } else {
        return nil, nil, fmt.Errorf("checkSplit(): err, line is: %s", line)
    }

    // startIPStr format err
    if startIP = net_util.ParseIPv4(startIPStr); startIP == nil {
        return nil, nil, fmt.Errorf("checkSplit(): line %s format err", line)
    }

    // endIPStr format err
    if endIP = net_util.ParseIPv4(endIPStr); endIP == nil {
        return nil, nil, fmt.Errorf("checkSplit(): line %s format err", line)
    }

    return startIP, endIP, nil
}

// checkLine checks line format
func checkLine(line string) (net.IP, net.IP, error) {
    var startIP, endIP net.IP
    var err error

    // check space split segment
    startIP, endIP, err = checkSplit(line, " ")
    if err != nil {
        // check tab split segment
        startIP, endIP, err = checkSplit(line, "\t")
        if err != nil {
            return nil, nil, fmt.Errorf("checkLine(): err, %s", err.Error())
        }
    }

    return startIP, endIP, err
}

/* check Version num and load IP txt file to IP items in memory */
func (f TxtFileLoader) CheckAndLoad(curVersion string) (*ipdict.IPItems, error) {
    var startIP, endIP net.IP

    fileName := f.fileName
    // get file Version and lineNum
    metaInfo, err := getFileInfo(fileName)
    if err != nil {
        return nil, fmt.Errorf("loadFile(): %s %s", fileName, err.Error())
    }
    newVersion := metaInfo.Version
    singleIPNum := metaInfo.SingleIPNum
    pairIPNum := metaInfo.PairIPNum

    // if singleIPNum + pairIPNum > maxLine
    // use maxline for singleIPNum and pairIPNum(protect malloc failed)
    // but the dict will still cut off by maxLine
    if f.maxLine != -1 && singleIPNum+pairIPNum > f.maxLine {
        singleIPNum = f.maxLine
        pairIPNum = f.maxLine
    }

    // check version
    if newVersion == curVersion && newVersion != "" {
        return nil, ERR_NO_NEED_UPDATE
    }

    // init counter for singleIP & pairIP
    singleIPCounter := 0
    pairIPCounter := 0
    lineCounter := 0
    // open file
    file, err := os.Open(fileName)
    if err != nil {
        return nil, fmt.Errorf("loadFile(): %s, %s", fileName, err.Error())
    }
    defer file.Close()
    // create ipItems
    ipItems, err := ipdict.NewIPItems(singleIPNum, pairIPNum)
    if err != nil {
        return nil, err
    }
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
            return nil, fmt.Errorf("loadFile(): err, %s", err.Error())
        }

        // insert start ip and end ip into dict
        if bytes.Compare(startIP, endIP) == 0 {
            // startIp == endIP insert single
            err = ipItems.InsertSingle(startIP)
            singleIPCounter += 1
        } else {
            err = ipItems.InsertPair(startIP, endIP)
            pairIPCounter += 1
        }
        if err != nil {
            return nil, fmt.Errorf("loadFile(): err, %s", err.Error())
        }

        // check if lineCounter > maxLine or not
        lineCounter += 1
        if f.maxLine != -1 && lineCounter > f.maxLine {
            //sort dict
            ipItems.Sort()
            ipItems.Version = newVersion
            return ipItems, ERR_MAX_LINE_EXCEED
        }

        // if ipcounter > max ipnum
        if singleIPCounter > singleIPNum || pairIPCounter > pairIPNum {
            //sort dict
            ipItems.Sort()
            ipItems.Version = newVersion
            return ipItems, ERR_MAX_LINE_EXCEED
        }
    }

    err = scanner.Err()
    // Scan meets error
    if err != nil {
        return nil, fmt.Errorf("loadFile(): err, %s", err.Error())
    }

    // Load succ, sort dict
    ipItems.Sort()
    ipItems.Version = newVersion
    return ipItems, nil
}
