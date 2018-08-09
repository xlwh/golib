/* internal_client.go - bfe remote client, use to talk to remote server in a sync/async manner   */
/*
modification history
--------------------
2014/7/1, by Weiwei, create
2016/3/10, by Sijie Yang, modify
    - merge bug fix for close from upstream source
      (see https://codereview.appspot.com/91230045)
*/
/*
DESCRIPTION
client example:
    // Dial(network, address, connectTimeout, fnCreateCodec, maxPendingNum)
    client := Dial("unix", "/tmp/waf.sock", 1000, nil, 100)

    // request, response, sync case:
    err := client.Call(request, response, 10*time.Millisecond)

	// request no response case
	err := client.GoNoReturn(request)
*/

package remote

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

import (
	"www.baidu.com/golang-lib/log"
)

const (
	MSG_TYPE_REQUEST             = 0
	MSG_TYPE_REQUEST_NO_RESPONSE = 1
	MSG_TYPE_RESPONSE            = 2

	DEFAULT_MAX_PENDING_CALLS = 10000 // default pending call(wait for send) queue length
)

var (
	MAGIC_STR = [4]byte{0xB0, 0xAE, 0xBE, 0xA7} // proto magic number
)

type IClient interface {
	Call(req interface{}, res interface{}, timeout time.Duration) error
	GoNoReturn(req interface{})
}

type Header struct {
	MagicStr    [4]byte // magic str
	MessageType uint16  // sequence number chosen by client(deprecated)
	Reserved    uint16
	MessageSize uint32 // message size, including header size
	Seq         uint32 // seq number, increase by 1
}

var ErrShutdown = errors.New("connection is shut down")
var ErrTimeout = errors.New("timeout error")
var ErrFull = errors.New("pending call full")

// Call represents an active call.
type Call struct {
	Req   interface{} // The request message.
	Res   interface{} // The reply from the server.
	Seq   uint32
	Error error // After completion, the error status.
	Type  uint16
	Done  chan *Call // Strobes when call is complete.
}

// InternalClient represents an Client talk to a remote server.
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
type InternalClient struct {
	codec   ClientCodec // codec for messages
	sending sync.Mutex  // mutex for send

	mutex   sync.Mutex // protects shutdown, closing, seq, pending
	seq     uint32
	pending map[uint32]*Call // register msgid with call
	calls   chan *Call       // buffered calls waiting to send

	closing    bool      // active close
	shutdown   bool      // shutdown when processing input message
	outputQuit chan bool // notify output channel to quit
}

func (client *InternalClient) send(call *Call) {
	client.sending.Lock()
	defer client.sending.Unlock()

	// check if client closed
	client.mutex.Lock()
	if client.shutdown || client.closing {
		call.Error = ErrShutdown
		client.mutex.Unlock()
		call.done()
		return
	}

	call.Seq = client.seq
	client.seq++

	// put request to call queue(implement by chan)
	select {
	case client.calls <- call:
		client.pending[call.Seq] = call
	default:
		// chan is full
		call.Error = ErrFull
	}

	client.mutex.Unlock()
}

// send message to server
func (client *InternalClient) output() {
	// sending header
	header := Header{MagicStr: MAGIC_STR}

loop:
	for {
		// wait for a pending or quit signal
		select {
		case <-client.outputQuit:
			break loop
		case call := <-client.calls:
			// Encode and send the request.
			seq := call.Seq
			header.Seq = seq
			header.MessageType = call.Type

			err := client.codec.WriteRequest(&header, call.Req)
			if err != nil {
				client.mutex.Lock()
				call = client.pending[seq]
				delete(client.pending, seq)
				client.mutex.Unlock()
				if call != nil {
					call.Error = err
					call.done()
				}
			}
		}
	}
}

// recv message from server
func (client *InternalClient) input() {
	var err error
	var response Header

	for err == nil {
		response = Header{}
		err = client.codec.ReadResponseHeader(&response)
		if err != nil {
			break
		}

		seq := response.Seq
		client.mutex.Lock()
		call := client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()

		switch {
		case call == nil:
			// We've got no pending call. That usually means that
			// WriteRequest partially failed, and call was already
			// removed; response is a server telling us about an
			// error reading request body. We should still attempt
			// to read error body, but there's no one to give it to.
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		default:
			err = client.codec.ReadResponseBody(call.Res)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}

	// Terminate pending calls.
	client.sending.Lock()
	client.mutex.Lock()
	client.shutdown = true
	closing := client.closing

	if err == io.EOF {
		if closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}

	// connection closed, now call done() for pending messages
	for k, call := range client.pending {
		call.Error = err
		call.done()
		delete(client.pending, k)
	}

	// quit output routine
	close(client.outputQuit)
	client.mutex.Unlock()
	client.sending.Unlock()

	if err != io.EOF && !closing {
		log.Logger.Debug("client protocol error: %s", err)
	}
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here.  It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
	}
}

// NewInternalClient returns a new Client to handle requests to the
// set of services at the other end of the connection.
// It adds a buffer to the write side of the connection so
// the header and payload are sent as a unit.
func NewInternalClient(conn io.ReadWriteCloser, fn fnCreateCodec, pendingNum int) *InternalClient {
	if fn == nil {
		fn = DefautCreateCodec
	}
	codec := fn(conn)
	return NewInternalClientWithCodec(codec, pendingNum)
}

// NewInternalClientWithCodec is like NewInternalClient but uses the specified
// codec to encode requests and decode responses.
func NewInternalClientWithCodec(codec ClientCodec, pendingNum int) *InternalClient {
	if pendingNum <= 0 {
		pendingNum = DEFAULT_MAX_PENDING_CALLS
	}
	client := &InternalClient{
		codec:   codec,
		pending: make(map[uint32]*Call),

		calls:      make(chan *Call, pendingNum),
		outputQuit: make(chan bool),
		closing:    false,
		shutdown:   false,
	}
	go client.input()
	go client.output()
	return client
}

// Dial connects to an server at the specified network address.
func Dial(network, address string, connectTimeout time.Duration, fn fnCreateCodec,
	pendingNum int) (*InternalClient, error) {
	conn, err := net.DialTimeout(network, address, connectTimeout)
	if err != nil {
		return nil, err
	}
	return NewInternalClient(conn, fn, pendingNum), nil
}

func (client *InternalClient) Close() error {
	client.mutex.Lock()
	if client.closing { // Close() called already
		client.mutex.Unlock()
		return ErrShutdown
	}

	client.closing = true
	client.mutex.Unlock()
	return client.codec.Close()
}

func (client *InternalClient) WaitClose() bool {
	<-client.outputQuit

	client.mutex.Lock()
	defer client.mutex.Unlock()

	return client.closing
}

// call RemovePendingOncall to tell client response is not needed
func (client *InternalClient) RemovePendingOnCall(call *Call) {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	delete(client.pending, call.Seq)
}

// Go process asynchronously.  It returns the Call structure representing
// the invocation.  The done channel will signal when the call is complete by returning
// the same Call object.  If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (client *InternalClient) Go(req interface{}, res interface{}, done chan *Call, msgType uint16) *Call {
	call := new(Call)
	call.Type = msgType
	call.Req = req
	call.Res = res

	if done == nil {
		done = make(chan *Call, 1) // buffered.
	} else {
		// If caller passes done != nil, it must arrange that
		// done has enough buffer for the number of simultaneous
		// RPCs that will be using that channel.  If the channel
		// is totally unbuffered, it's best not to run at all.
		if cap(done) == 0 {
			log.Logger.Debug("rpc: done channel is unbuffered")
		}
	}
	call.Done = done

	client.send(call)
	return call
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (client *InternalClient) Call(req interface{}, res interface{}, timeout time.Duration) error {
	call := client.Go(req, res, make(chan *Call, 1), MSG_TYPE_REQUEST)
	select {
	case <-call.Done:
		return call.Error
	case <-time.After(timeout):
		client.RemovePendingOnCall(call)
		return ErrTimeout
	}
}

// return pending call num
func (client *InternalClient) pendingCallNum() int {
	if client == nil {
		return 0
	}
	return len(client.calls)
}
