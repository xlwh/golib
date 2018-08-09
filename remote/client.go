/* client.go - wrapper for internal client, support multi-connection, auto reconnect on fail*/
/*
modification history
--------------------
2014/7/01, by Weiwei, create
2015/3/03, by Sijie Yang, modify
    - support connections to different addresses
2015/5/27, by Sijie Yang, modify
    - support Weighted Round Robin load balancing
    - reload support for server address and weight
2015/5/29, by Sijie Yang, modify
    - support BNS for service instance discovery
*/
/*
DESCRIPTION
client example:
    // NewClient(network, address, connectTimeout, concurrency, fn, pendingNum)
    // create 8 connections to server address unix: /tmp/waf.sock
    // connect timeout is 1s
    // each connection has a process pending queue(waiting for send), size is 100
    client := NewClient("unix", "/tmp/waf.sock", 1*time.Second, 8, nil, 100)

    // request, response, sync case:
    // wait timeout 10ms
    err := client.Call(request, response, 10*time.Millisecond)

	// request no response case
	err := client.GoNoReturn(request)
*/

package remote

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"www.baidu.com/golang-lib/bns"
	"www.baidu.com/golang-lib/log"
)

// client connect to server status
const (
	UNCONNECTED int32 = iota // not connected
	CONNECTING               // client is connecting but failed or waiting for result
	CONNECTED                // client is connected to server
)

// client wrapper maitain connection status to make sure only 1 routine is connecting to server
type SClient struct {
	iclient  *InternalClient
	lock     sync.RWMutex // lock for InternalClient
	StopDial chan bool    // stop AsyncDial
	Status   int32        // atomic status

	Weight  int32  // initial weight
	Left    int32  // atomic weight left
	Address string // address of server connected to
}

func NewSClient(address string, weight int32, left int32) *SClient {
	client := new(SClient)
	client.StopDial = make(chan bool, 0)
	client.Address = address
	client.Weight = weight
	client.Left = left
	return client
}

// async dial to server
// try to keep connection until Close() called
func (s *SClient) AsyncDial(network, address string, timeout time.Duration, fn fnCreateCodec,
	pendingNum int) {
	// go routine will dial every second until success
	for s.shouldDial() {
		interval := 1 * time.Second
		conn, err := net.DialTimeout(network, address, timeout)
		if err == nil {
			client := NewInternalClient(conn, fn, pendingNum)
			s.SetInternal(client)
			atomic.StoreInt32(&s.Status, CONNECTED)

			// wait until client stop
			closing := client.WaitClose()
			atomic.StoreInt32(&s.Status, CONNECTING)

			// close internal client conn
			if !closing {
				client.Close()
			}
		} else {
			log.Logger.Debug("connect to %s %s failed %s, %v", network, address,
				err.Error(), timeout)
			s.wait(interval)
		}
	}
}

func (s *SClient) shouldDial() bool {
	select {
	case <-s.StopDial:
		return false
	default:
		return true
	}
}

func (s *SClient) wait(interval time.Duration) {
	select {
	case <-time.After(interval):
		break
	case <-s.StopDial:
		break
	}
}

func (s *SClient) CloseAll() {
	close(s.StopDial)
	c := s.GetInternal()
	if c != nil {
		c.Close()
	}
}

func (s *SClient) GetInternal() *InternalClient {
	s.lock.RLock()
	c := s.iclient
	s.lock.RUnlock()
	return c
}

func (s *SClient) SetInternal(c *InternalClient) {
	s.lock.RLock()
	s.iclient = c
	s.lock.RUnlock()
}

// Wapper for client
// 1. support multi-connection to the same address
// 2. support auto-reconnect on fail
type Client struct {
	network        string
	addrInfo       map[string]int32 // map for server info (address, weight)
	connectTimeout time.Duration    // connect timeout
	pendingNum     int              // max pending msg number waiting to send for each internal client
	fn             fnCreateCodec    // used to add create codec
	mutex          sync.Mutex       // protect clients and index
	clients        []*SClient       // actual clients
	index          int              // current client index
	concurrency    int              // number of clients per address
	closed         bool             // client closed
	bnsClient      *bns.Client      // bns client
}

//
// NewClient create a client to communicate with server
// network, address represent server address
// 	network could be any network streaming protocol, "tcp/unix"
// connectTimeout: connect server timeout
// concurrency: internal connections
// fn: function to create codec, if fn is nil, pbcodec is used as default
// pendingNum is pending queue(waiting for send message) max size
func NewClient(network, address string, connectTimeout time.Duration, concurrency int, fn fnCreateCodec,
	pendingNum int) *Client {
	addrInfo := map[string]int32{address: 1000}
	return NewClientByAddr(network, addrInfo, connectTimeout, concurrency, fn, pendingNum)
}

/* NewClientByAddr - Create a client to communicate with server
 *
 * Params:
 *     - network       : network type
 *     - addrInfo      : map for server info: <address, weight>
 *     - connectTimeout: connect server timeout
 *     - concurrency   : number of connections to per address
 *     - fn            : function to create codec (pbcodec default)
 *     - pendingNum    : max size of pending queue
 *
 * Return:
 *     - Client
 */
func NewClientByAddr(network string, addrInfo map[string]int32,
	connectTimeout time.Duration, concurrency int, fn fnCreateCodec,
	pendingNum int) *Client {
	if concurrency <= 0 {
		concurrency = 1
	}

	// try to establish connection once created
	clients := make([]*SClient, 0, len(addrInfo)*concurrency)
	for address, weight := range addrInfo {
		for i := 0; i < concurrency; i++ {
			client := NewSClient(address, weight, weight)
			clients = append(clients, client)
			go client.AsyncDial(network, address, connectTimeout, fn, pendingNum)
		}
	}

	// create and initialize Client
	rec := &Client{
		network:        network,
		addrInfo:       addrInfo,
		connectTimeout: connectTimeout,
		pendingNum:     pendingNum,
		fn:             fn,
		clients:        clients,
		concurrency:    concurrency,
	}

	return rec
}

/* NewClientByName - Create a client to communicate with server
 *
 * Params:
 *     - serviceName   : service name
 *     - connectTimeout: connect server timeout
 *     - concurrency   : number of connections to per address
 *     - fn            : function to create codec (pbcodec default)
 *     - pendingNum    : max size of pending queue
 *
 * Return:
 *     - client        : client instance
 *     - error         : error if fail
 */
func NewClientByName(serviceName string, connectTimeout time.Duration,
	concurrency int, fn fnCreateCodec, pendingNum int) *Client {
	bnsClient := bns.NewClient()
	addrInfo, err := bns.GetAddrAndWeight(bnsClient, serviceName)
	if err != nil {
		log.Logger.Info("remote: get instance info for %v error(%s)",
			serviceName, err.Error())
		addrInfo = make(map[string]int32)
	}

	client := NewClientByAddr("tcp", addrInfo, connectTimeout, concurrency, fn, pendingNum)
	client.bnsClient = bnsClient
	go client.checkServiceInstance(serviceName)

	return client
}

func (rec *Client) checkServiceInstance(serviceName string) {
	for !rec.closed {
		time.Sleep(10 * time.Second)

		// get service instances by name
		addrInfo, err := bns.GetAddrAndWeight(rec.bnsClient, serviceName)
		if err != nil {
			log.Logger.Info("remote: get instance info for %v error(%s)",
				serviceName, err.Error())
			continue
		}
		if len(addrInfo) == 0 {
			log.Logger.Error("no instance configured for %v", serviceName)
			continue
		}
		if reflect.DeepEqual(addrInfo, rec.addrInfo) {
			continue
		}

		// update server addr and weight
		rec.Update(addrInfo)
	}
}

// try to get an available client instance
func (rec *Client) Balance() (*SClient, error) {
	rec.mutex.Lock()
	defer rec.mutex.Unlock()

	client, err := rec.balance()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// select a client using weighted round robin balancing
// Note: just return SClient for unit test
func (rec *Client) balance() (*SClient, error) {
	if len(rec.clients) == 0 {
		return nil, fmt.Errorf("no clients in pool")
	}

	var client *SClient
	allBroken := true
	used := rec.index

	for {
		client = rec.clients[used]

		if atomic.LoadInt32(&client.Status) == CONNECTED && client.Left > 0 {
			break
		}

		if atomic.LoadInt32(&client.Status) == CONNECTED && client.Left <= 0 {
			allBroken = false
		}

		// next client to check
		used = (used + 1) % len(rec.clients)

		// if all clients have been checked
		if used == rec.index {
			if allBroken {
				return nil, fmt.Errorf("all connections to servers are broken")
			} else {
				rec.initWeight()
				rec.index = 0
				used = 0
			}
		}
	}
	client.Left -= 1
	rec.index = (used + 1) % len(rec.clients)

	return client, nil
}

func (rec *Client) initWeight() {
	for i := 0; i < len(rec.clients); i++ {
		rec.clients[i].Left = rec.clients[i].Weight
	}
}

/* Update - update server addr and weights
 *
 * Params:
 *     - addrInfo: map for server info <address, weight>
 */
func (rec *Client) Update(addrInfo map[string]int32) error {
	rec.mutex.Lock()
	defer rec.mutex.Unlock()

	// filter entries with weight <= 0
	for address, weight := range addrInfo {
		if weight <= 0 {
			delete(addrInfo, address)
		}
	}
	if len(addrInfo) == 0 {
		return fmt.Errorf("no available address")
	}

	// process new addresses
	clients := make([]*SClient, 0, len(addrInfo)*rec.concurrency)
	for address, weight := range addrInfo {
		if _, ok := rec.addrInfo[address]; !ok {
			for i := 0; i < rec.concurrency; i++ {
				client := NewSClient(address, weight, weight)
				clients = append(clients, client)
				go client.AsyncDial(rec.network, address, rec.connectTimeout,
					rec.fn, rec.pendingNum)
			}
		}
	}

	// process exsiting addresses
	for _, client := range rec.clients {
		if _, ok := addrInfo[client.Address]; ok {
			client.Weight = addrInfo[client.Address]
			client.Left = client.Weight
			clients = append(clients, client)
		} else {
			client.CloseAll()
		}
	}

	rec.addrInfo = addrInfo
	rec.clients = clients
	rec.index = 0
	return nil
}

// try to get an available client instance
func (rec *Client) getClient() (*InternalClient, error) {
	client, err := rec.Balance()
	if err != nil {
		return nil, err
	}

	return client.GetInternal(), nil
}

// just send requst ,no response.
func (rec *Client) GoNoReturn(req interface{}) {
	client, err := rec.getClient()
	if err != nil {
		log.Logger.Warn("no available connection now")
		return
	}

	call := client.Go(req, nil, nil, MSG_TYPE_REQUEST_NO_RESPONSE)
	if call.Error != nil {
		log.Logger.Warn("GoNoReturn send return error %s", call.Error)
	}
	client.RemovePendingOnCall(call)
}

// sync call, before the response is revceived, wait for timeout at most
func (rec *Client) Call(req interface{}, res interface{}, timeout time.Duration) error {
	client, err := rec.getClient()
	if err == nil {
		err = client.Call(req, res, timeout)
	}

	return err
}

func (rec *Client) Close() error {
	rec.mutex.Lock()
	defer rec.mutex.Unlock()

	rec.closed = true
	for i := 0; i < int(rec.concurrency); i++ {
		rec.clients[i].CloseAll()
	}

	return nil
}

// get pending call number of client
// return []int, each represent a internal client PendingCallNum
func (rec *Client) PendingCallNum() []int {
	rec.mutex.Lock()
	defer rec.mutex.Unlock()

	pendingNum := make([]int, rec.concurrency)

	for i := 0; i < rec.concurrency; i++ {
		client := rec.clients[i].GetInternal()
		if client != nil {
			pendingNum[i] = client.pendingCallNum()
		}
	}

	return pendingNum
}
