/* http_util.go - some actions of http client, e.g., get, post  */
/*
modification history
--------------------
2015/4/13, by Zhang Miao, create
2016/1/8, by Taochunhua, add http headers support
*/
/*
DESCRIPTION
*/
package http_util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

import (
	"www.baidu.com/golang-lib/log"
)

// type for specify http post Content-Type in header
type ContentType string

const CONTENT_FORM ContentType = "application/x-www-form-urlencoded"
const CONTENT_JSON ContentType = "application/json"
const CONTENT_XML ContentType = "application/xml"

/*
generate http request, with given (method, urlPath, ipAddr)

param:
    method:     http method, GET or POST
    urlPath:    e.g., http://www.baidu.com/index.html
    ipAddr:     dst ip address, e.g., 220.181.111.188. if it is "", ipAddr and port are ignored
    port:       dst port, e.g., 8080
    dataType:   CONTENT_FORM, CONTENT_JSON or CONTENT_XML
    data:       data for body in http request. Its type is specified by dataType.
    headers:    other http headers
return:
    (*http.Request, error)
*/
func requestGen(method, urlPath, ipAddr, port string, dataType ContentType,
	data []byte, headers map[string]string) (*http.Request, error) {
	var req *http.Request
	var err error

	// check method
	switch method {
	case "GET", "POST", "DELETE", "PUT", "HEAD":
	default:
		return nil, fmt.Errorf("invalid method:%s", method)
	}

	// create new request
	if data != nil && len(data) != 0 {
		req, err = http.NewRequest(method, urlPath, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, urlPath, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("http.NewRequest():%s", err.Error())
	}

	// insert http headers
	req.Header.Set("Content-Type", string(dataType))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// modify request.URL.Host, if necessary
	if ipAddr != "" {
		req.URL.Host = fmt.Sprintf("%s:%s", ipAddr, port)
	}

	return req, nil
}

/*
do http request for given (urlPath, ipAddr) within timeout secs

param:
    method:  GET, or POST
    urlPath: e.g., http://www.baidu.com/index.html
    timeout: in seconds
    ipAddr: manually set ipAddr
    port:   manually set port
    dataType: CONTENT_FORM, CONTENT_JSON or CONTENT_XML
    data:    data to post (this param is for post only). Its type is specified by dataType.
    headers: other http headers
return:
    (response, error)
*/
func Do(method, urlPath string, timeout int, ipAddr, port string,
	dataType ContentType, data []byte, headers map[string]string) ([]byte, error) {
	// prepare http request
	req, err := requestGen(method, urlPath, ipAddr, port, dataType, data, headers)
	if err != nil {
		return nil, fmt.Errorf("requestGen():%s", err.Error())
	}

	// prepare http client
	t := time.Duration(timeout) * time.Second
	client := NewTimeoutClient(t, t)

	// send out request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client.Do():%s", err.Error())
	}

	// important! if not close, connection/goroutine/memory leaks
	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code:%d", resp.StatusCode)
	}

	// read from response
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll():%s", err.Error())
	}

	return bytes, nil
}

/*
read from given urlPath within timeout secs

param:
    urlPath: e.g., http://www.baidu.com/index.html
    timeout: in seconds
    headers: other http headers

return:
    (response, error)
*/
func Read(urlPath string, timeout int, headers map[string]string) ([]byte, error) {
	// Content-Type is useless for GET
	return Do("GET", urlPath, timeout, "", "0", CONTENT_FORM, nil, headers)
}

/*
read from given url, with given host IP address, within timeout secs

param:
    urlPath: e.g., http://www.baidu.com/index.html
    timeout: in seconds
    ipAddr: manually set ipAddr
    port:   manually set port
    headers: other http headers

return:
    (response, error)
*/
func ReadByIPAddr(urlPath string, timeout int,
	ipAddr, port string, headers map[string]string) ([]byte, error) {
	// Content-Type is useless for GET
	return Do("GET", urlPath, timeout, ipAddr, port, CONTENT_FORM, nil, headers)
}

/*
post data to given url, within timeout secs

param:
    urlPath: e.g., http://www.baidu.com/index.html
    timeout: in seconds
    dataType: CONTENT_FORM, CONTENT_JSON or CONTENT_XML
    data:   data to post. Its type is specified by dataType.
    headers: other http headers

return:
    (response, error)
*/
func Post(urlPath string, timeout int, dataType ContentType,
	data []byte, headers map[string]string) ([]byte, error) {
	return Do("POST", urlPath, timeout, "", "0", dataType, data, headers)
}

/*
post data to given url, within timeout secs

param:
    urlPath: e.g., http://www.baidu.com/index.html
    timeout: in seconds
    ipAddr: manually set ipAddr
    port:   manually set port
    dataType: CONTENT_FORM, CONTENT_JSON or CONTENT_XML
    data:   data to post. Its type is specified by dataType.
    headers: other http headers

return:
    (response, error)
*/
func PostByIPAddr(urlPath string, timeout int, ipAddr, port string,
	dataType ContentType, data []byte, headers map[string]string) ([]byte, error) {
	return Do("POST", urlPath, timeout, ipAddr, port, dataType, data, headers)
}

/*
Resolve hostname in urlPath to (ipAddrs, port)

param:
    urlPath: e.g., http://www.baidu.com/index.html

return:
    (addrs, port, error), e.g., [['127.0.0.1', '127.0.0.2'], "80", nil]
*/
func urlResolve(urlPath string) ([]string, string, error) {
	var port string

	//  Parse a URL into 6 components:
	//       <scheme>://<netloc>/<path>;<params>?<query>#<fragment>
	rawUrl, err := url.Parse(urlPath)
	if err != nil {
		return nil, "0", fmt.Errorf("url.Parse(%s):%s", urlPath, err.Error())
	}

	// get host, port from raw host name
	host, port, err := net.SplitHostPort(rawUrl.Host)
	if err != nil {
		host = rawUrl.Host
		port = "80"
	}

	// get ip address from hostname
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil, "0", fmt.Errorf("net.LookupHost(%s):%s", host, err.Error())
	}

	return addrs, port, nil
}

/*
read from given urlPath by round robin, within timeout secs for each IP

it will try to read from all ipaddrs from dns round robin

param:
    urlPath: e.g., http://www.baidu.com/index.html
    timeout: in seconds
    headers: other http headers

return:
    (response, error)
*/
func ReadRR(urlPath string, timeout int, headers map[string]string) ([]byte, error) {
	var err error
	var data []byte

	// resolve ip addresses from urlPath
	addrs, port, err := urlResolve(urlPath)
	if err != nil {
		return nil, fmt.Errorf("urlResolve(%s):%s", urlPath, err.Error())
	}

	// try to read from each address
	for _, addr := range addrs {
		data, err = ReadByIPAddr(urlPath, timeout, addr, port, headers)
		if err != nil {
			log.Logger.Warn("ReadRR():ReadByIPAddr(%s, %s, %d):%s",
				urlPath, addr, port, err.Error())
		} else {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("fail to read:%s", err.Error())
	} else {
		return data, nil
	}
}

/*
post to given urlPath by round robin, within timeout secs for each IP

it will try to post to all ipaddrs from dns round robin

param:
    urlPath: e.g., http://www.baidu.com/index.html
    timeout: in seconds
    dataType: CONTENT_FORM, CONTENT_JSON or CONTENT_XML
    data: data to post. Its type is specified by dataType.
    headers: other http headers

return:
    (response, error)
*/
func PostRR(urlPath string, timeout int, dataType ContentType,
	data []byte, headers map[string]string) ([]byte, error) {
	var err error
	var res []byte

	// resolve ip addresses from urlPath
	addrs, port, err := urlResolve(urlPath)
	if err != nil {
		return nil, fmt.Errorf("urlResolve(%s):%s", urlPath, err.Error())
	}

	// try to read from each address
	for _, addr := range addrs {
		res, err = PostByIPAddr(urlPath, timeout, addr, port, dataType, data, headers)
		if err != nil {
			log.Logger.Warn("PostRR():PostByIPAddr(%s, %s, %d):%s",
				urlPath, addr, port, err.Error())
		} else {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("fail to post")
	} else {
		return res, nil
	}
}

/*
post data to given url, within timeout secs, with given retry times

param:
    urlPath: e.g., http://www.baidu.com/index.html
    timeout: in seconds
    dataType: CONTENT_FORM, CONTENT_JSON or CONTENT_XML
    data:   data to post. Its type is specified by dataType.
    headers: other http headers
    retry: number of max retry times
    interval: time interval between retries. in micro-seconds.

return:
    (response, error)
    error is the last post's error.
*/
func PostWithRetry(urlPath string, timeout int, dataType ContentType,
	data []byte, headers map[string]string, retry int, interval int) ([]byte, error) {
	var resp []byte
	var err error
	for i := 0; i < retry; i++ {
		resp, err = Post(urlPath, timeout, dataType, data, headers)
		if err == nil {
			break
		}
		time.Sleep(time.Microsecond * time.Duration(interval))
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}
