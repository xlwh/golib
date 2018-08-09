/* timeout_client.go - http client with timeout control */
/*
modification history
--------------------
2015/4/13, by Zhang Miao, create
*/
/*
DESCRIPTION
see http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
*/
package http_util

import (
    "net"
    "net/http"
    "time"
)

func TimeoutDialer(cTimeout time.Duration, 
                   rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
    return func(netw, addr string) (net.Conn, error) {
        conn, err := net.DialTimeout(netw, addr, cTimeout)
        if err != nil {
            return nil, err
        }
        conn.SetDeadline(time.Now().Add(rwTimeout))
        return conn, nil
    }
}

func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
        },
    }
}