/* util_funcs.go - util functions for net_server   */
/*
modification history
--------------------
2014/3/10, by Zhang Miao, create
2014/8/6, by Zhang Miao, move code from waf_server
*/
/*
DESCRIPTION
*/
package net_server

import (
	"errors"
	"io"
	"net"
)

import (
	"www.baidu.com/golang-lib/module_state2"
)

var (
	ErrNoDataClose = errors.New("close without receiving any data")
)

/* read specific len of data from connection    */
func ReadWithLen(conn net.Conn, buff []byte, len int) ([]byte, error) {
	totalRecved := 0 // total recved from socket

	for {
		recved, err := conn.Read(buff[totalRecved:len])
		// check special error
		if err == io.EOF && totalRecved <= 0 {
			return nil, ErrNoDataClose
		}
		// check other read error
		if err != nil {
			return nil, err
		}

		// record total receive size
		totalRecved = totalRecved + recved
		if totalRecved >= len {
			break
		}
	}
	return buff[0:totalRecved], nil
}

/* read and bypass specific len of data from connection */
func ReadBypass(conn net.Conn, buff []byte, length int) error {
	var toRead int
	toReadAll := length // total recved from socket
	buffLen := len(buff)

	for {
		if buffLen <= toReadAll {
			toRead = buffLen
		} else {
			toRead = toReadAll
		}

		recved, err := conn.Read(buff[0:toRead])

		if err != nil {
			return err
		}

		toReadAll = toReadAll - recved
		if toReadAll <= 0 {
			break
		}
	}
	return nil
}

/* write msg to connection  */
func WriteMsg(conn net.Conn, msg []byte) error {
	msgLen := len(msg)
	totalSent := 0

	for {
		sent, err := conn.Write(msg[totalSent:msgLen])

		if err != nil {
			return err
		}

		totalSent = totalSent + sent
		if totalSent >= msgLen {
			break
		}
	}
	return nil
}

/* read waf-header from connection  */
func ReadHeader(conn net.Conn, buff []byte, len int,
	tcpState *module_state2.State) (MsgHeader, error) {
	var header MsgHeader

	/* read from socket */
	readBuf, err := ReadWithLen(conn, buff, len)
	if err != nil {
		if tcpState != nil {
			tcpState.Inc("TCP_READ_HEADER_ERR", 1)
		}
		return header, err
	}

	/* decode header    */
	header, err = MsgHeaderDecode(readBuf)

	if err != nil {
		if tcpState != nil {
			tcpState.Inc("TCP_DECODE_HEADER_ERR", 1)
		}
		return header, err
	}

	/* compare with magic string    */
	if header.MagicStr != MAGIC_STR {
		err = errors.New("magic number is wrong")
		if tcpState != nil {
			tcpState.Inc("TCP_HEADER_MAGIC_ERR", 1)
		}
		return header, err
	}

	return header, nil
}
