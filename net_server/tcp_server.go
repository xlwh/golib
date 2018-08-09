/* tcp_server.go - process request from client  */
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
    "encoding/json"
    "errors"
    "fmt"
    "net"
)

import (
    "www.baidu.com/golang-lib/log"
    "www.baidu.com/golang-lib/module_state2"
)

/* size of receiving buffer */
const RECV_BUF_LEN = 200 * 1024

// command in SendMsg
const (
    SEND_CMD_QUIT = 0   // ask sender to quit
    SEND_CMD_SEND = 1   // ask sender to send msg
)

type SendMsg struct {
    cmd     int         // command, QUIT, SEND
    data    interface{} // data to send
}

// callback functions used in TcpServer
type CallBacks interface {
    // for process recved msg
    RecvMsgProc(header MsgHeader, body []byte, tcpConn TcpConn, 
                tcpState *module_state2.State) error
    // for prepare sending msg
    SendMsgMake(data interface{}) ([]byte, error)
}

type TcpServer struct {
    Listener    net.Listener        // for listen
    
    tcpState    module_state2.State // for collecting state
    tcpTable    TcpConnTable        // for maintain tcp conn
    
    callBacks   CallBacks           // callback functions    
}

func NewTcpServer(callBacks CallBacks) *TcpServer {
    srv := new(TcpServer)
    
    // initialize tcpTable
    srv.tcpTable.Init()
    
    // initialize tcpState
    srv.tcpState.Init()
    
    // set callBacks
    srv.callBacks = callBacks
    
    return srv
}

// handleListen - listen to tcp port
func (srv *TcpServer) handleListen() {
    for {
        conn, err := srv.Listener.Accept()
        if err != nil {
            log.Logger.Error("handleListen():error in accept(): %s", err.Error())
            srv.tcpState.Inc("TCP_ACCEPT_ERR", 1)
            continue
        }
        srv.tcpState.Inc("TCP_ACCEPT_SUCC", 1)

        /* add to tcp connection table  */
        tcpConn := srv.tcpTable.Add(conn)
        /* start two go-routines for this connection    */
        srv.startTcpConn(conn, tcpConn)
        
    }
}

/*
* ListenAndServe - tcp listen and start a new go-routine to serve
*
* PARAMS:
*   - port: listen port, e.g., "80"
*
* RETURNS: 
*   nil, if succeed
*   error, if fail
*/
func (srv *TcpServer) ListenAndServe(port string) error {
    port = ":" + port
    
    /* try to listen    */
    var err error
    srv.Listener, err = net.Listen("tcp", port)
    if err != nil {
        errStr := fmt.Sprintf("err in listen(): %s", err.Error())
        return errors.New(errStr)
    }

    /* use seperate go-routine to handle listen */
    go srv.handleListen()
    
    return nil
}

/*
* ListenAndServe2 - tcp listen and start a new go-routine to serve
*
* PARAMS:
*   - netType: "tcp", "tcp4", "tcp6", "unix" or "unixpacket"
*   - addr: address, e.g., ":80"(for "tcp"), or "/tmp/echo.sock" for "unix"
*     visit http://golang.org/pkg/net/#Listen for more information
*
* RETURNS: 
*   nil, if succeed
*   error, if fail
*/
func (srv *TcpServer) ListenAndServe2(netType, addr string) error {   
    /* try to listen    */
    var err error
    srv.Listener, err = net.Listen(netType, addr)
    if err != nil {
        errStr := fmt.Sprintf("err in listen(): %s", err.Error())
        return errors.New(errStr)
    }

    /* use seperate go-routine to handle listen */
    go srv.handleListen()
    
    return nil
}


/* get state counters of tcp server */
func (srv *TcpServer) TcpStateGet() *module_state2.State {
    return &srv.tcpState
}

/* get tcp conn table of tcp server in json */
func (srv *TcpServer) TcpTableGetJson() ([]byte, error) {
    output := srv.tcpTable.GetState()

    /* convert to json  */
    buff, err := json.Marshal(output)
    
    return buff, err    
}

/* start go-routines for new tcp connection  */
func (srv *TcpServer) startTcpConn(conn net.Conn, tcpConn TcpConn) {
    /* start goroutine for receiving data   */
    go srv.tcpRecver(conn, tcpConn)
    /* start goroutine for sending data   */
    go srv.tcpSender(conn, tcpConn)
}

/* read waf-body from connection  */
func (srv *TcpServer) readBody(conn net.Conn, buff []byte, len int) ([]byte, error) {
    /* read from socket */
    readBuf, err := ReadWithLen(conn, buff, len)
    if err != nil {
        srv.tcpState.Inc("TCP_READ_BODY_ERR", 1)
        return readBuf, err
    }

    return readBuf, nil
}

func (srv *TcpServer) removeTcpConn(conn net.Conn) {
    // remove from tcpTable
    tcpConn, err := srv.tcpTable.Remove(conn)    
    if err == nil {
        // close tcp connection
        conn.Close()
        
        // notify sender go-routine to quit    
        pMsg := &SendMsg{cmd:SEND_CMD_QUIT}
        tcpConn.SendQueue.Append(pMsg)
    }
}

/*  tcpRecver - handle incoming msgs    */
func (srv *TcpServer) tcpRecver(conn net.Conn, tcpConn TcpConn) {
    log.Logger.Debug("tcpRecver start")
    var buff [RECV_BUF_LEN]byte
            
    for {
        /* read header  */
        header, err := ReadHeader(conn, buff[:], MSG_HEADER_LEN, &srv.tcpState)
        
        if err != nil && err != ErrNoDataClose {
            log.Logger.Warn("ReadHeader(): %s %s", conn.RemoteAddr().String(), err.Error())            
            break
        }

        if err == ErrNoDataClose {
            log.Logger.Info("ReadHeader(): %s %s", conn.RemoteAddr().String(), err.Error())
            break
        }
        
        /* check msg type   */
        if header.MsgType != MSG_TYPE_REQUEST &&
                header.MsgType != MSG_TYPE_REQUEST_NO_RESPONSE {
            log.Logger.Warn("tcpRecver():error: msg type is not request[%d]",
                            header.MsgType)
            srv.tcpState.Inc("TCP_READ_MSGTYPE_ERR", 1)
            break
        }        
        
        /* calc and check body length   */
        bodyLen := int(header.MsgSize) - MSG_HEADER_LEN
        if bodyLen <= 0 {
            log.Logger.Warn("tcpRecver():err: bodyLen = %d", bodyLen)
            srv.tcpState.Inc("TCP_BODY_ZERO_ERR", 1)
            break            
        }        
        if bodyLen > RECV_BUF_LEN {            
            log.Logger.Warn("tcpRecver():err: bodyLen = %d", bodyLen)
            srv.tcpState.Inc("TCP_BODY_TOOLONG_ERR", 1)
            // bypass msg with bodyLen > RECV_BUF_LEN
            err = ReadBypass(conn, buff[:], bodyLen)
            if err != nil {
                log.Logger.Warn("ReadBypass():", err.Error())
                srv.tcpState.Inc("TCP_READ_BYPASS_ERR", 1)
                break
            } else {
                // to read next header
                continue
            }
        }

        /* read body    */
        var body []byte
        body, err = srv.readBody(conn, buff[:], bodyLen)        
        if err != nil {
            log.Logger.Warn("tcpRecver():err in readBody():", err.Error())
            break
        }

        // process received msg
        srv.callBacks.RecvMsgProc(header, body, tcpConn, &srv.tcpState)
    }
    srv.removeTcpConn(conn)
}

/*  tcpSender - handle outgoing msgs    */
func (srv *TcpServer) tcpSender(conn net.Conn, tcpConn TcpConn) {
    log.Logger.Debug("tcpSender start")
    for {
        // read SendMsg from sending queue
        sendMsg := tcpConn.SendQueue.Remove().(*SendMsg)        
        if sendMsg.cmd == SEND_CMD_QUIT {
            // quit go-routine of tcpSender
            break
        }
        
        // check cmd
        if sendMsg.cmd != SEND_CMD_SEND {
            log.Logger.Warn("cmd should be SEND_CMD_SEND, now:%d", sendMsg.cmd)
            continue
        }
        
        // prepare msg ([]byte) to send
        msg, err := srv.callBacks.SendMsgMake(sendMsg.data)
        if err != nil {
            log.Logger.Warn("tcpSender():err in SendMsgMake():%s", err.Error())
            srv.tcpState.Inc("RESPONSE_MAKE_ERR", 1)
            continue
        }

        // send out the msg
        err = WriteMsg(conn, msg)        
        if err != nil {
            log.Logger.Warn("tcpSender:err in WriteMsg(): %s", err.Error())
            srv.tcpState.Inc("RESPONSE_SEND_ERR", 1)
            
            /* remove the connection    */
            srv.removeTcpConn(conn)
            /* quit the loop    */
            break
        }
                
        log.Logger.Debug("tcpSender:succ in write()")
        srv.tcpState.Inc("RESPONSE_SEND_OK", 1)
    }
    log.Logger.Debug("tcpSender quit")
}
