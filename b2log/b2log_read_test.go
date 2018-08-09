/* b2log_read_test.go - test for b2log_read.go  */
/*
modification history
--------------------
2014/11/27, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package b2log

import (
    "io/ioutil"
    "testing"
)

// test of unsafe.Sizeof()
func Test_Sizeof_1(t *testing.T) {
    if HEADER_SIZE != 24 {
        t.Error("unsafe.Sizeof(Header)")
    }
}

// test of BuffParse(), case 1
// It's the normal situation
func Test_BuffParse_1(t *testing.T) {
    // read testing data from file
    data, err := ioutil.ReadFile("test_data/pb_access_1.log")
    if err != nil {
        t.Error("fail to open file for testing data")
        return
    }
    
    // parse b2log record from data
    records, buffer := BuffParse(data)
    
    if len(records) != 8 {
        t.Errorf("len(records) should be 8, but now it's %d", len(records))
    }
    if len(buffer) != 0 {
        t.Errorf("len(buffer) should be 0, but now it's %d", len(buffer))
    }
}

// test of BuffParse(), case 2
// magic number of first record is break
func Test_BuffParse_2(t *testing.T) {
    // read testing data from file
    data, err := ioutil.ReadFile("test_data/pb_access_2.log")
    if err != nil {
        t.Error("fail to open file for testing data")
        return
    }
    
    // parse b2log record from data
    records, buffer := BuffParse(data)
    
    if len(records) != 7 {
        t.Errorf("len(records) should be 7, but now it's %d", len(records))
    }
    if len(buffer) != 0 {
        t.Errorf("len(buffer) should be 0, but now it's %d", len(buffer))
    }
}

// test of BuffParse(), case 3
// compress_len of first record is not zero
func Test_BuffParse_3(t *testing.T) {
    // read testing data from file
    data, err := ioutil.ReadFile("test_data/pb_access_3.log")
    if err != nil {
        t.Error("fail to open file for testing data")
        return
    }
    
    // parse b2log record from data
    records, buffer := BuffParse(data)
    
    if len(records) != 7 {
        t.Errorf("len(records) should be 7, but now it's %d", len(records))
    }
    if len(buffer) != 0 {
        t.Errorf("len(buffer) should be 0, but now it's %d", len(buffer))
    }
}

// test of BuffParse(), case 4
// compress_len of first record is not zero, and in the first read, read only 32 bytes
func Test_BuffParse_4(t *testing.T) {
    // read testing data from file
    data, err := ioutil.ReadFile("test_data/pb_access_3.log")
    if err != nil {
        t.Error("fail to open file for testing data")
        return
    }
    
    // parse b2log record from data[0:32]
    records, buffer := BuffParse(data[0:32])
    
    if len(records) != 0 {
        t.Errorf("len(records) should be 0, but now it's %d", len(records))
    }
    if len(buffer) != 32 {
        t.Errorf("len(buffer) should be 32, but now it's %d", len(buffer))
    }
    
    // parse b2log record from data
    records, buffer = BuffParse(data)
    
    if len(records) != 7 {
        t.Errorf("len(records) should be 7, but now it's %d", len(records))
    }
    if len(buffer) != 0 {
        t.Errorf("len(buffer) should be 0, but now it's %d", len(buffer))
    }    
}

// test of BuffParse(), case 5
// compress_len of first record is larger than 100K
func Test_BuffParse_5(t *testing.T) {
    // read testing data from file
    data, err := ioutil.ReadFile("test_data/pb_access_4.log")
    if err != nil {
        t.Error("fail to open file for testing data")
        return
    }
    
    // parse b2log record from data
    records, buffer := BuffParse(data)
    
    if len(records) != 507 {
        t.Errorf("len(records) should be 507, but now it's %d", len(records))
    }
    if len(buffer) != 0 {
        t.Errorf("len(buffer) should be 0, but now it's %d", len(buffer))
    }
}

// test of BuffParse(), case 6
// uncompress_len of first record is larger than 100K
func Test_BuffParse_6(t *testing.T) {
    // read testing data from file
    data, err := ioutil.ReadFile("test_data/pb_access_5.log")
    if err != nil {
        t.Error("fail to open file for testing data")
        return
    }
    
    // parse b2log record from data
    records, buffer := BuffParse(data)
    
    if len(records) != 508 {
        t.Errorf("len(records) should be 508, but now it's %d", len(records))
    }
    if len(buffer) != 0 {
        t.Errorf("len(buffer) should be 0, but now it's %d", len(buffer))
    }
}