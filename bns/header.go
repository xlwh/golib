package bns

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	headerSize     = 20
	headerMagicNum = 0xfb709394
)

type header struct {
	Id       uint16 `json:"id"`
	Version  uint16 `json:"version"`
	LogId    uint32 `json:"log_id"`
	MagicNum uint32 `json:"magic_num"`
	Reserved uint32 `json:"reserved"`
	BodyLen  uint32 `json:"bodylen"`
}

func (h *header) unmarshal(b []byte) error {
	if len(b) < headerSize {
		return errors.New("incomplete header")
	}
	h.Id = binary.LittleEndian.Uint16(b[0:2])
	h.Version = binary.LittleEndian.Uint16(b[2:4])
	h.LogId = binary.LittleEndian.Uint32(b[4:8])
	h.MagicNum = binary.LittleEndian.Uint32(b[8:12])
	h.Reserved = binary.LittleEndian.Uint32(b[12:16])
	h.BodyLen = binary.LittleEndian.Uint32(b[16:20])
	return nil
}

func (h *header) marshal(b []byte) error {
	if len(b) < headerSize {
		return errors.New("not enough buffer for header")
	}
	binary.LittleEndian.PutUint16(b[0:2], h.Id)
	binary.LittleEndian.PutUint16(b[2:4], h.Version)
	binary.LittleEndian.PutUint32(b[4:8], h.LogId)
	binary.LittleEndian.PutUint32(b[8:12], h.MagicNum)
	binary.LittleEndian.PutUint32(b[12:16], h.Reserved)
	binary.LittleEndian.PutUint32(b[16:20], h.BodyLen)
	return nil
}

func (h *header) write(w io.Writer) (n int, err error) {
	var buf [headerSize]byte
	if err = h.marshal(buf[:]); err != nil {
		return 0, err
	}
	return w.Write(buf[:])
}

func (h *header) read(r io.Reader) (n int, err error) {
	var buf [headerSize]byte
	if n, err = io.ReadFull(r, buf[:]); err != nil {
		return n, err
	}
	if err = h.unmarshal(buf[:]); err != nil {
		return n, err
	}
	return n, nil
}
