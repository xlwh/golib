package bns

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	// DefaultTimeout is the default socket read/write timeout.
	DefaultTimeout = 100 * time.Millisecond

	MaxIdleConnsPerAddr = 2
)

type client struct {
	Timeout time.Duration
	Addr    string
}

// debugclientConnections controls whether all client connections are
// wrapped with a verbose logging wrapper
var debugclientConnections = false

func (c *client) getConn(addr string) (cn *clientConn, err error) {
	nc, err := c.dial(addr)
	if err != nil {
		return nil, err
	}
	if debugclientConnections {
		nc = newLoggingConn("client", nc)
	}
	cn = &clientConn{
		nc:   nc,
		addr: addr,
		rw:   bufio.NewReadWriter(bufio.NewReader(nc), bufio.NewWriter(nc)),
		c:    c,
	}
	return cn, nil
}

func (c *client) netTimeout() time.Duration {
	if c.Timeout != 0 {
		return c.Timeout
	}
	return DefaultTimeout
}

type clientConn struct {
	nc   net.Conn
	rw   *bufio.ReadWriter
	addr string
	c    *client
}

func (cn *clientConn) close() error {
	return cn.nc.Close()
}

func (c *client) do(req *request) (resp *response, err error) {
	err = c.withConn(func(rw *bufio.ReadWriter) error {
		if _, err := req.write(rw); err != nil {
			return err
		}
		if err := rw.Flush(); err != nil {
			return err
		}
		rsp, err := readResponse(rw)
		if err != nil {
			return err
		}
		resp = rsp
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) doWithRetry(req *request, n int) (resp *response, err error) {
	for i := 1; i < n; i++ {
		resp, err = c.do(req)
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			continue
		}
		return resp, err
	}
	return resp, err
}

func (c *client) withConn(fn func(*bufio.ReadWriter) error) (err error) {
	cn, err := c.getConn(c.Addr)
	if err != nil {
		return err
	}
	defer cn.close()

	deadline := time.Now().Add(c.netTimeout())
	cn.nc.SetWriteDeadline(deadline)
	cn.nc.SetReadDeadline(deadline)
	return fn(cn.rw)
}

func (c *client) dial(addr string) (net.Conn, error) {
	nc, err := net.DialTimeout("tcp", addr, c.netTimeout())
	if err == nil {
		return nc, nil
	}
	return nil, err
}

var (
	uniqNameMu   sync.Mutex
	uniqNameNext = make(map[string]int)
)

func newLoggingConn(baseName string, c net.Conn) net.Conn {
	uniqNameMu.Lock()
	defer uniqNameMu.Unlock()
	uniqNameNext[baseName]++
	return &loggingConn{
		name: fmt.Sprintf("%s-%d", baseName, uniqNameNext[baseName]),
		Conn: c,
	}
}

type loggingConn struct {
	name string
	net.Conn
}

func (c *loggingConn) Write(p []byte) (n int, err error) {
	log.Printf("%s.Write(%d) %s -> %s", c.name, len(p), c.Conn.LocalAddr(), c.Conn.RemoteAddr())
	n, err = c.Conn.Write(p)
	log.Printf("%s.Write(%d) = %d, %v\n%v", c.name, len(p), n, err, p[:n])
	return
}

func (c *loggingConn) Read(p []byte) (n int, err error) {
	log.Printf("%s.Read(%d) %s <- %s", c.name, len(p), c.Conn.LocalAddr(), c.Conn.RemoteAddr())
	n, err = c.Conn.Read(p)
	log.Printf("%s.Read(%d) = %d, %v\n%v", c.name, len(p), n, err, p[:n])
	return
}

func (c *loggingConn) Close() (err error) {
	log.Printf("%s.Close() = ...", c.name)
	err = c.Conn.Close()
	log.Printf("%s.Close() = %v", c.name, err)
	return
}

type request struct {
	Header header
	Body   []byte
}

func (r *request) write(w io.Writer) (n int, err error) {
	n, err = r.Header.write(w)
	if err != nil {
		return 0, err
	}
	return w.Write(r.Body)
}

func newRequest(t MsgType, body []byte) *request {
	req := new(request)
	req.Header.MagicNum = headerMagicNum
	req.Header.Id = uint16(t)
	req.Header.BodyLen = uint32(len(body))
	req.Body = body
	return req
}

type response struct {
	Header header
	Body   []byte
}

func readResponse(r io.Reader) (resp *response, err error) {
	resp = new(response)
	_, err = resp.Header.read(r)
	if err != nil {
		return nil, err
	}
	if resp.Header.MagicNum != headerMagicNum {
		return nil, fmt.Errorf("invalid magic number %x", resp.Header.MagicNum)
	}

	resp.Body = make([]byte, int(resp.Header.BodyLen))
	_, err = io.ReadFull(r, resp.Body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
