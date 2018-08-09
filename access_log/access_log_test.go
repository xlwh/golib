/* access_log_test.go - test for access_log.go */
/*
modification history
--------------------
2014/6/9, by Zhang Miao, create
*/
package access_log

import (
    "testing"
    "time"
)

// test LoggerInit(prefix, logDir)
func TestLog_1(t *testing.T) {
    logger, err := LoggerInit("test", "./log", "M", 2)
    
    if err != nil {
        t.Error("LoggerInit() fail")
    }
    
    for i:=0; i < 100; i = i + 1 {
        logger.Info("info msg: %d", i)
        
        // time.Sleep(10 * time.Second)
    }
    
    time.Sleep(100 * time.Millisecond)
}

// test LoggerInit2(fileName, logDir)
func TestLog_2(t *testing.T) {
    logger, err := LoggerInit2("test.output", "./log", "M", 2)
    
    if err != nil {
        t.Error("LoggerInit() fail")
    }
    
    for i:=0; i < 100; i = i + 1 {
        logger.Info("info msg: %d", i)
        
        // time.Sleep(10 * time.Second)
    }
    
    time.Sleep(100 * time.Millisecond)
}

// test LoggerInit3(filePath)
func TestLog_3(t *testing.T) {
    logger, err := LoggerInit3("./log/test.out", "M", 2)
    
    if err != nil {
        t.Error("LoggerInit() fail")
    }
    
    for i:=0; i < 100; i = i + 1 {
        logger.Info("info msg: %d", i)
        
        // time.Sleep(10 * time.Second)
    }
    
    time.Sleep(100 * time.Millisecond)
}