/* http_util_test.go - test of http_util.go */
/*
modification history
--------------------
2015/4/14, by Zhang Miao, create
2016/1/8, by Taochunhua, add http headers support
*/
/*
DESCRIPTION
*/
package http_util

import (
    "bytes"
    //"fmt"
    "net"
    //"strings"
    "testing"
    //"time"
)

// test for requestGen(), without given (ipaddr, port)
func Test_requestGen_Case1(t *testing.T) {
    // invoke requestGen()
    req, err := requestGen("GET", "http://www.baidu.com/index.html?a=1#pos1",
                           "", "80", CONTENT_FORM, nil, nil)
    if err != nil {
        t.Errorf("requestGen():%s", err.Error())
        return
    }
    
    // check area of req
    if req.Method != "GET" {
        t.Errorf("Method should be GET, now it's %s", req.Method)
    }

    if req.Host != "www.baidu.com" {
        t.Errorf("Host should be www.baidu.com, now it's %s", req.Host)
    }

    if req.URL.Host != "www.baidu.com" {
        t.Errorf("URL.Host should be www.baidu.com, now it's %s", req.URL.Host)
    }
    
    if req.URL.Path != "/index.html" {
        t.Errorf("URL.Path should be /index.html, now it's %s", req.URL.Path)
    }

    if req.URL.RawQuery != "a=1" {
        t.Errorf("URL.RawQuery should be a=1, now it's %s", req.URL.RawQuery)
    }

    if req.URL.Fragment != "pos1" {
        t.Errorf("URL.Fragment should be pos1, now it's %s", req.URL.Fragment)
    }    
}

// test for requestGen(), with given (ipaddr, port)
func Test_requestGen_Case2(t *testing.T) {
    // invoke requestGen()
    req, err := requestGen("GET", "http://www.baidu.com",
                            "220.181.112.244", "80", CONTENT_FORM, nil, nil)
    if err != nil {
        t.Errorf("requestGen():%s", err.Error())
        return
    }
    
    if req.URL.Host != "220.181.112.244:80" {
        t.Errorf("URL.Fragment should be 220.181.112.244:80, now it's %s", req.URL.Host)
    }
}

// test for requestGen(), with post data
func Test_requestGen_Case3(t *testing.T) {
    var data = []byte(`{"a":1}`)

    // invoke requestGen()
    req, err := requestGen("POST", "http://www.baidu.com", "", "0", CONTENT_FORM, data, nil)
    if err != nil {
        t.Errorf("requestGen():%s", err.Error())
        return
    }
    
    if req.ContentLength != int64(len(data)) {
        t.Errorf("req.ContentLength = %d", req.ContentLength)
    }
}

// test for Read()
func Test_Read_Case1(t *testing.T) {
    // invoke Read()
    data, err := Read("http://bfe.baidu.com", 60, nil)
    if err != nil {
        t.Errorf("Read():%s", err.Error())
        return
    }
    
    if ! bytes.HasPrefix(data, []byte("<!doctype html>")) {
        t.Errorf("err in Read()")
    }
}

// test for ReadByIPAddr()
func Test_ReadByIPAddr_Case1(t *testing.T) {
    // resolve hostname to ip
    addrs, err := net.LookupHost("bfe.baidu.com")
    if err != nil {
        t.Errorf("err in net.LookupHost(bfe.baidu.com)")
        return
    }

    // invoke ReadByIPAddr()
    data, err := ReadByIPAddr("http://bfe.baidu.com", 60, addrs[0], "80", nil)
    if err != nil {
        t.Errorf("ReadByIPAddr():%s", err.Error())
        return
    }
    
    if ! bytes.HasPrefix(data, []byte("<!doctype html>")) {
        t.Errorf("err in Read()")
    }
}

// test for PostByIPAddr()
func Test_PostByIPAddr_Case1(t *testing.T) {
    // resolve hostname to ip
    addrs, err := net.LookupHost("gtc.baidu.com")
    if err != nil {
        t.Errorf("err in net.LookupHost(bfe.baidu.com)")
        return
    }

    // invoke PostByIPAddr()
    _, err = PostByIPAddr("http://gtc.baidu.com", 60, addrs[0], "80", CONTENT_JSON, []byte("{}"), nil)
    if err != nil {
        t.Errorf("PostByIPAddr():%s", err.Error())
    }
}

// test for urlResolve()
func Test_urlResolve_Case1(t *testing.T) {
    addrs, port, err := urlResolve("http://bfe.baidu.com/indexl.html")
    if err != nil {
        t.Errorf("urlResolve(): return error:%s", err.Error())
        return
    }
    if len(addrs) == 0 {
        t.Errorf("urlResolve(): len(addrs) should >= 1")
    }
    if port != "80" {
        t.Errorf("urlResolve(): port should be 80")
    }
}

// test for ReadRR()
func Test_ReadRR_Case1(t *testing.T) {
    // invoke ReadRR()
    data, err := ReadRR("http://bfe.baidu.com", 60, nil)
    if err != nil {
        t.Errorf("ReadRR():%s", err.Error())
        return
    }
    
    if ! bytes.HasPrefix(data, []byte("<!doctype html>")) {
        t.Errorf("err in ReadRR()")
    }
}

// test for PostRR()
func Test_PostRR_Case1(t *testing.T) {
    // invoke PostRR()
    _, err := PostRR("http://gtc.baidu.com", 60, CONTENT_JSON, []byte("{}"), nil)
    if err != nil {
        t.Errorf("PostRR():%s", err.Error())
        return
    }
}
