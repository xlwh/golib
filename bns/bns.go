/* bns.go - bns client for golang */
/*
modification history
--------------------
2015/5/28, by Sijie Yang, added
*/
/*
DESCRIPTION
    This package is forked from http://git.baidu.com/ksarch/gobns
*/
package bns

import (
	"time"
)

import (
    "code.google.com/p/goprotobuf/proto"
)

//go:generate protoc --go_out=. naming.proto naminglib.proto service.proto

type MsgType int

const (
	ReqService         MsgType = 1
	ReqAuthService     MsgType = 2
	ResService         MsgType = 3
	ResAuthService     MsgType = 4
	ReqServiceList     MsgType = 9
	ReqAuthServiceList MsgType = 10
	ResServiceList     MsgType = 11
	ResAuthServiceList MsgType = 12
	ReqServiceConf     MsgType = 15
	ResServiceConf     MsgType = 16
	ReqRealService     MsgType = 30
)

var (
	LocalAddr  = "localhost:793"
	RemoteAddr = "unamed.noah.baidu.com:793"
)

const (
	defaultRetryTimes = 2
)

type Client struct {
	// socket read/write timeout.
	// if zero, DefaultTimeout is used
	Timeout time.Duration

	local  *client
	remote *client
}

func NewClient() *Client {
	return &Client{
		local:  &client{Addr: LocalAddr},
		remote: &client{Addr: RemoteAddr},
	}
}

func (c *Client) Call(args proto.Message, reply proto.Message) error {
	var reqType MsgType
	switch args.(type) {
	case *LocalNamingRequest:
		reqType = ReqService
	case *LocalNamingAuthRequest:
		reqType = ReqAuthService
	case *LocalServiceConfRequest:
		reqType = ReqServiceConf
	case *LocalNamingListRequest:
		reqType = ReqServiceList
	case *LocalNamingAuthListRequest:
		reqType = ReqAuthServiceList
	}

	content, err := proto.Marshal(args)
	if err != nil {
		return err
	}
	req := newRequest(reqType, content)

	var resp *response
	resp, err = c.local.do(req)
	if err != nil {
		resp, err = c.remote.doWithRetry(req, defaultRetryTimes)
		if err != nil {
			return err
		}
	}

	return proto.Unmarshal(resp.Body, reply)
}
