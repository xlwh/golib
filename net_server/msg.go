/* msg.go - process msg between client and net-server */
/*
modification history
--------------------
2014/3/12, by Zhang Miao, create
2014/8/6, by Zhang Miao, move from waf_server
*/
/*
DESCRIPTION
*/
package net_server

import (
    "bytes"
    "encoding/binary"
)

var MAGIC_STR [4]byte = [4]byte{0xB0, 0xAE, 0xBE, 0xA7}

/* type of msg  */
const (
    MSG_TYPE_REQUEST                = 0
    MSG_TYPE_REQUEST_NO_RESPONSE    = 1
    MSG_TYPE_RESPONSE               = 2
)

type MsgHeader struct {
    MagicStr        [4]byte     // for validate start of header
    MsgType         uint16      // see define of msg types
    Reserved        uint16      // for 4-byte alignment
    MsgSize         uint32      // full size of msg, i.e., len(header+body)
    ReqId           uint32      // ID of request
}

// modify the global MAGIC_STR
func MagicStrSet(magicStr [4]byte) {
    MAGIC_STR = magicStr
}

/* get length of MsgHeader   */
func MsgHeaderLen() int {
    var header MsgHeader

    length := binary.Size(header)
    
    return length
}

var MSG_HEADER_LEN = MsgHeaderLen()

/* convert MsgHeader to binary buffer   */
func MsgHeaderEncode(header MsgHeader) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, &header)

    if err != nil {
        return nil, err
    } else {
        return buf.Bytes(), nil
    }
}

/* make header of request msg    */
func RequestHeaderMake(reqId uint32, bodyLen int, needResponse bool) ([]byte, error) {
    var header MsgHeader
    
    header.MagicStr = MAGIC_STR
    
    if needResponse {
        header.MsgType = MSG_TYPE_REQUEST
    } else {
        header.MsgType = MSG_TYPE_REQUEST_NO_RESPONSE
    }
    
    header.ReqId = reqId
    header.MsgSize = uint32(MSG_HEADER_LEN +  bodyLen)

    buff, err := MsgHeaderEncode(header)
    
    return buff, err
}

/* convert from binary buffer to MsgHeader  */
func MsgHeaderDecode(buf []byte) (MsgHeader, error) {
	var header MsgHeader
    
    bufReader := bytes.NewReader(buf)
	err := binary.Read(bufReader, binary.BigEndian, &header)

    return header, err
}

/* make header of response msg    */
func ResponseHeaderMake(reqId uint32, bodyLen int) ([]byte, error) {
    var header MsgHeader
    header.MagicStr = MAGIC_STR
    header.MsgType = MSG_TYPE_RESPONSE
    header.ReqId = reqId
    header.MsgSize = uint32(MSG_HEADER_LEN +  bodyLen)

    buff, err := MsgHeaderEncode(header)
    
    return buff, err
}
