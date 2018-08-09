/* access_log.go - encapsulation of log4go to print access log  */
/*
modification history
--------------------
2014/06/09, by Zhang Miao, create
2014/12/18, by Zhang Miao, modify, add LoggerInit2() and LoggerInitWithFormat2()
                                   use fileName instead of prefix
*/
/*
DESCRIPTION
log: encapsulation for log4go

Usage:
    import "www.baidu.com/golang-lib/access_log"

    // One log file will be generated in ./log: test.log
    // The log will rotate, and there is support for backup count 
    logger := access_log.LoggerInit("test", "./log", "midnight", 5)
     
    logger.Info("msg1")
     
    // it is required, to work around bug of log4go
    time.Sleep(100 * time.Millisecond) 
*/
package access_log

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

import "code.google.com/p/log4go"

/* logDirCreate(): check and create dir if nonexist   */
func logDirCreate(logDir string) error {
    if _, err := os.Stat(logDir); os.IsNotExist(err) {
        /* create directory */
        err = os.MkdirAll(logDir, 0777)
        if err != nil {
            return err
        }
    }
    return nil
}

// fullPathGen - generate full path of the file
func fullPathGen(fileName, logDir string) string {
    // remove the last '/'
    strings.TrimSuffix(logDir, "/")    
    fullPath := filepath.Join(logDir, fileName)    
    return fullPath    
}

// generate fileName from prefix
func prefix2Name(prefix string) string {
    return prefix + ".log"
}

/*
* LoggerInit - initialize logger
*
* PARAMS:
*   - prefix: Name of log file will be prefix.log
*   - logDir: directory for log. It will be created if noexist
*   - when: 
*       "M", minute
*       "H", hour
*       "D", day
*       "MIDNIGHT", roll over at midnight
*   - backupCount: If backupCount is > 0, when rollover is done, no more than 
*       backupCount files are kept - the oldest ones are deleted.
*
* RETURNS: 
*   (logger, nil), if succeed
*   (logger, error), if fail
*/
func LoggerInit(prefix string, logDir string, when string, backupCount int) (log4go.Logger, error) {
    fileName := prefix2Name(prefix)
    return LoggerInit2(fileName, logDir, when, backupCount)
}

// similar to LoggerInit
// instead of prefix, fileName should be provided
func LoggerInit2(fileName, logDir, when string, backupCount int) (log4go.Logger, error) {
    accessDefaultFormat := "%M"
    return LoggerInitWithFormat2(fileName, logDir, when, backupCount, accessDefaultFormat)
}

// similar to LoggerInit2
// instead of (fileName, logDir), filePath should be provided
func LoggerInit3(filePath, when string, backupCount int) (log4go.Logger, error) {
    logDir, fileName := filepath.Split(filePath)
    return LoggerInit2(fileName, logDir, when, backupCount)
}

/*
* LoggerInitWithFormat - initialize logger
*
* PARAMS:
*   - prefix: Name of log file will be prefix.log
*   - logDir: directory for log. It will be created if noexist
*   - when: 
*       "M", minute
*       "H", hour
*       "D", day
*       "MIDNIGHT", roll over at midnight
*   - backupCount: If backupCount is > 0, when rollover is done, no more than 
*       backupCount files are kept - the oldest ones are deleted.
*   - format: log4j supported format
* RETURNS: 
*   (logger, nil), if succeed
*   (logger, error), if fail
*/
func LoggerInitWithFormat(prefix, logDir, when string, backupCount int, 
                          format string) (log4go.Logger, error) {
    fileName := prefix2Name(prefix)
    return LoggerInitWithFormat2(fileName, logDir, when, backupCount, format)
}

// similar to LoggerInit
// instead of prefix, fileName should be provided
func LoggerInitWithFormat2(fileName, logDir, when string, backupCount int, 
                          format string) (log4go.Logger, error) {
    var logger log4go.Logger

    // check value of when is valid
    if ! log4go.WhenIsValid(when) {
        log4go.Error("LoggerInitWithFormat(): invalid value of when(%s)", when)
        return logger, fmt.Errorf("invalid value of when: %s", when)        
    }
    
    // change when to upper
    when = strings.ToUpper(when)
              
    // check, and create dir if nonexist
    if err := logDirCreate(logDir); err != nil {
        log4go.Error("Init(), in logDirCreate(%s)", logDir)
        return logger, err
    }

    // create logger
    logger = make(log4go.Logger)
        
    // create file writer for all log
    fullPath := fullPathGen(fileName, logDir)
    logWriter := log4go.NewTimeFileLogWriter(fullPath, when, backupCount)
    if logWriter == nil {
        return logger, fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fullPath)
    }
    logWriter.SetFormat(format)
    logger.AddFilter("log", log4go.INFO, logWriter)

    return logger, nil
}

/*
* LoggerInitWithSvr - initialize logger with remote log server
*
* PARAMS:
*   - progName: program name
*   - loggerName: logger name
*   - network: using "udp" or "unixgram"
*   - svrAddr: remote address
*
* RETURNS: 
*   (logger, nil), if succeed
*   (logger, error), if fail
*/
func LoggerInitWithSvr(progName string, loggerName string, 
                       network string, svrAddr string) (log4go.Logger, error) {
    var logger log4go.Logger
    
    /* create file writer for all log   */
    name := fmt.Sprintf("%s_%s", progName, loggerName)
    
    // create logger
    logger = make(log4go.Logger)
    
    logWriter := log4go.NewPacketWriter(name, network, svrAddr, log4go.LogFormat)
    if logWriter == nil {
        return nil, fmt.Errorf("error in log4go.NewPacketWriter(%s)", name)
    }
    logger.AddFilter(name, log4go.INFO, logWriter)
    
    return logger, nil
}
