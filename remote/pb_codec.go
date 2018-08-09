/* pb_codec.go - pb codec implement interface ClientCodec    */
/*
modification history
--------------------
2014/7/1, by Weiwei02, create
2014/9/28, modified by weiwei, reuse buffer to reduce memory assumption
        accept message type Gogopb instead of go pb
*/
/*
DESCRIPTION
 Message = header + body

 default protobuf implementation:
    header is encode/decode using encoding/binary. BigEndian
    |MagicStr(4)|MessageType(2)|ReservedSize(2)|MessageSize(4)|Sequence(4)|

    body is encode/decode using ProtoBuf, so body must implement interface proto.Message.
    MessageSize = headerSize + bodySize.

*/

package remote

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"sync"
)

import (
	"code.google.com/p/gogoprotobuf/proto"
)

// binary header encoding + pb body encoding
type PbClientCodec struct {
	rwc io.ReadWriteCloser
	dec *PbDecoder
	enc *PbEncoder

	// buffer to the write side of the connection so the header
	// and payload are sent as a unit
	encBuf *bufio.Writer
}

func NewPbClientCodec(wrc io.ReadWriteCloser) ClientCodec {
	encBuf := bufio.NewWriter(wrc)
	return &PbClientCodec{wrc, NewPbDecoder(bufio.NewReader(wrc)), NewPbEncoder(encBuf), encBuf}
}

type PbEncoder struct {
	w io.Writer
}

func NewPbEncoder(w io.Writer) *PbEncoder {
	return &PbEncoder{w: w}
}

type PbDecoder struct {
	r        io.Reader
	BodySize uint32 // save body size parsed from header
}

func NewPbDecoder(r io.Reader) *PbDecoder {
	d := new(PbDecoder)
	d.r = r
	return d
}

var (
	ErrTooLarge = errors.New("required slice size exceed 64k")
)

var (
	// buffer to put generated pbmessage, most message should be less than 4k
	buf4kPool  sync.Pool
	buf16kPool sync.Pool
	buf64kPool sync.Pool
)

// get proper []byte from pool
// if size > 64K, return ErrTooLarge
func newBuffer(size int) ([]byte, error) {
	var pool *sync.Pool

	// return buffer size
	originSize := size

	if size <= 4096 {
		size = 4096
		pool = &buf4kPool
	} else if size <= 16*1024 {
		size = 16 * 1024
		pool = &buf16kPool
	} else if size <= 64*1024 {
		size = 64 * 1024
		pool = &buf64kPool
	} else {
		// if message is larger than 16K, return err
		return nil, ErrTooLarge
	}

	if v := pool.Get(); v != nil {
		return v.([]byte)[:originSize], nil
	}

	return make([]byte, size)[:originSize], nil
}

func putBuffer(b []byte) {
	b = b[:cap(b)]
	if cap(b) == 4096 {
		buf4kPool.Put(b)
	}
	if cap(b) == 16*1024 {
		buf16kPool.Put(b)
	}
	if cap(b) == 64*1024 {
		buf64kPool.Put(b)
	}
}

// GogoPB don't define a gogomarshal interface
// define it here
type GogoMarshaler interface {
	proto.Marshaler

	// 2 additional function implementation by GogoPB
	MarshalTo([]byte) (int, error)
	Size() int
}

func (enc *PbEncoder) Encode(r *Header, body interface{}) error {
	pb, ok := body.(GogoMarshaler)
	if !ok {
		return ErrCodecType
	}

	// get from pb recycle buffer
	size := pb.Size()
	buf, err := newBuffer(size)
	if err != nil {
		// if size is too large, alloc directly
		buf = make([]byte, size)
	} else {
		defer putBuffer(buf)
	}

	// marshal to buffer
	n, err := pb.MarshalTo(buf)
	if err != nil {
		return err
	}

	r.MessageSize = uint32(len(buf) + HEADER_LEN)
	if err := binary.Write(enc.w, binary.BigEndian, r); err != nil {
		return err
	}
	_, err = enc.w.Write(buf[:n])

	return err
}

func (dec *PbDecoder) DecodeHeader(r *Header) error {
	err := binary.Read(dec.r, binary.BigEndian, r)
	if err == nil {
		dec.BodySize = r.MessageSize - uint32(HEADER_LEN)
	}
	return err
}

func (dec *PbDecoder) Decode(body interface{}) error {
	buf, err := newBuffer(int(dec.BodySize))
	if err != nil {
		// message should be read
		// so if newBuffer failed, make a new slice
		buf = make([]byte, dec.BodySize)
	} else {
		defer putBuffer(buf)
	}

	_, err = io.ReadFull(dec.r, buf)
	if err == nil && body != nil {
		pb, ok := body.(proto.Message)
		if !ok {
			return ErrCodecType
		}
		err = proto.Unmarshal(buf, pb)
	}
	return err
}

func (c *PbClientCodec) WriteRequest(r *Header, body interface{}) error {
	if err := c.enc.Encode(r, body); err != nil {
		return err
	}
	return c.encBuf.Flush()
}

func (c *PbClientCodec) ReadResponseHeader(r *Header) error {
	return c.dec.DecodeHeader(r)
}

func (c *PbClientCodec) ReadResponseBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *PbClientCodec) Close() error {
	return c.rwc.Close()
}
