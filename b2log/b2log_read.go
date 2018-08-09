/* b2log_read.go - read b2log record from file  */
/*
modification history
--------------------
2014/11/4, by Zhang Miao, create
*/
/*
DESCRIPTION
*/
package b2log

import (
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
)

var (
	ErrNoEnoughData = errors.New("No enough data")
	ErrCompressed   = errors.New("Compress is not support")
)

/* 
BuffParse - read b2log records from given buffer

Params:
- buffer: buffer with binary data

Returns:
    (records, buffer)
    - records: [record, record, ...]
    - buffer: after removed the decoded record, or pass the unsynchronized data
*/
func BuffParse(buffer []byte) ([]Record, []byte) {
    var hasNext bool
    var record Record
    var err error
    
    records := make([]Record, 0)
    
    for {
        hasNext, record, buffer, err = recordParse(buffer)    
        if err == nil  {
            records = append(records, record)
        } else {
            // TBD: record if it is compressed
        }
    
        if !hasNext {
            break
        }
    }
    
    return records, buffer
}

/*
recordParse - read one b2log record from given buffer

Params:
- buffer: buffer with binary data

Returns:
    (hasNext, record, buffer, error)
    - hasNext: True, if need next read to the buffer
    - record: one record read from the buffer. Only available when err == nil
    - buffer: after removed the record, or pass the unsynchronized data
*/
func recordParse(buffer []byte) (bool, Record, []byte, error) {
    var record Record
    var offset int
    var dataLen int
    
    if len(buffer) < HEADER_SIZE {
        // no enough data to get header
        return false, record, buffer, ErrNoEnoughData
    }
    
    // read loghead
    logHeader, buffer, err := logHeaderRead(buffer)    
    if err != nil {
        // fail to read Header from buffer
        return true, record, buffer, fmt.Errorf("read header:%s", err.Error())
    }

    // check whether it is compressed record
    if logHeader.CompressLen != 0 {
        dataLen = int(logHeader.CompressLen)
        if dataLen > MAX_RECORD_LEN {
            // maybe logHeader.CompressLen is not correct
            dataLen = MAX_RECORD_LEN;
        }
            
        // Compression is not supported now.
        // So if there is compressed record, just bypass it and report error        
        offset = HEADER_SIZE + dataLen
                 
        if len(buffer) >= offset {
            // bypass the record, only when record is completely in the buffer
            buffer = buffer[offset:]
            return true, record, buffer, ErrCompressed
        } else {
            return false, record, buffer, ErrCompressed
        }
    }

    dataLen = int(logHeader.UnCompressLen)
    if dataLen > MAX_RECORD_LEN {
        // maybe pHead->uncompress_len is not correct
        dataLen = MAX_RECORD_LEN
    }    
    
    // check whether record is completely in the buffer    
    offset = HEADER_SIZE + dataLen

    if len(buffer) < offset {
        // no enough data, wait for the next time
        return false, record, buffer, ErrNoEnoughData
    }
    
    // get record out of the buffer
    record = buffer[HEADER_SIZE:offset]
    buffer = buffer[offset:]
    
    return true, record, buffer, nil
}

/*
logHeaderRead - read one b2log header from given buffer

Params:
- buffer: buffer with binary data

Returns:
    (header, buffer, error)
    - header: b2log header
    - buffer: after removed the record, or pass the unsynchronized data
*/
func logHeaderRead(buffer []byte) (Header, []byte, error) {
	var header Header

    // get data buffer of header
    headerBuf := buffer[0:HEADER_SIZE]
        
    // try to unpack buffer to b2log header
    bufReader := bytes.NewReader(headerBuf)
	err := binary.Read(bufReader, binary.LittleEndian, &header)
    if err != nil {
        // this should not happen, bypass the header in buffer
        buffer = buffer[HEADER_SIZE:]
        buffer = tryFindNextStart(buffer)

        return header, buffer, fmt.Errorf("header decode:%s", err.Error())
    }

    // check magic number
    if header.MagicNumber != MAGIC_NUMBER {
        buffer = tryFindNextStart(buffer)
        return header, buffer, fmt.Errorf("invalid magic number:0x%x", header.MagicNumber)
    }
    
    return header, buffer, nil
}

/*
tryFindNextStart - try to find the next start of b2log record

Params:
- buffer: buffer with binary data

Returns:
    buffer, removed the unsynchronized bytes
*/
func tryFindNextStart(buffer []byte) []byte {
    // get length of buffer
    length := len(buffer)
    
    // find the next position of magic-number
    offset := 0
    for offset < (length - 3) {
        if bytes.Compare(buffer[offset:(offset + 4)], MAGIC_NUMBER_STR) == 0 {
            break
        }

        offset = offset + 1
    }
    
    // modify buffer
    buffer = buffer[offset:]
    
    return buffer          
}