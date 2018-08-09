/* codec.go - codec provide ClientCodec interface definition for ser/deser message   */
/*
modification history
--------------------
2014/7/1, by Weiwei, create
*/
/*
DESCRIPTION
*/

package remote

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrCodecType = errors.New("not expecting codec type")

// A ClientCodec implements writing of requests and
// reading of responses for the client side of an session.
// The client calls WriteRequest to write a request to the connection
// and calls ReadResponseHeader and ReadResponseBody in pairs
// to read responses.  The client calls Close when finished with the
// connection. ReadResponseBody may be called with a nil
// argument to force the body of the response to be read and then
// discarded.
type ClientCodec interface {
	// WriteRequest must be safe for concurrent use by multiple goroutines.
	WriteRequest(*Header, interface{}) error
	ReadResponseHeader(*Header) error
	ReadResponseBody(interface{}) error

	Close() error
}

// codec create func type def
type fnCreateCodec func(wrc io.ReadWriteCloser) ClientCodec

// protocol header length
var HEADER_LEN = func() int { return binary.Size(Header{}) }()

var DefautCreateCodec = NewPbClientCodec
