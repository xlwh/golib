/* redis_bns.go - redis client with bns support*/
/*
modification history
--------------------
2016/03/21, by zhangjiyang01, port from tls_sessioncache.go
*/
/*
DESCRIPTION
    redis client with bns support

Usage:
    import "www.baidu.com/golang-lib/redis_bns"

    bnsName := "bfe-tc.bfe.tc" // bns of redis server
    maxIdle := 10              // max Idle connection
    connectTimeout := 10       // connection timeout in ms
    readTimeout    := 10       // read redis server timeout in ms
    writeTimeout   := 10       // write redis server timeout in ms
    redisClient :=  redis_bns.NewRedisClient(bnsName, maxIdle,
                                connectTimeout, readTimeout, writeTimeout, module_state2.State)

    // put/get/incr/decr
    redisClient.Put("key", "val", expireTime)
    redisClient.Get("key")
    redisClient.Incr("key", expireTime)
    redisClient.Decr("key")
*/

package redis_bns

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"
)

import (
	"github.com/garyburd/redigo/redis"
	"www.baidu.com/golang-lib/bns"
	"www.baidu.com/golang-lib/delay_counter"
	"www.baidu.com/golang-lib/log"
	"www.baidu.com/golang-lib/module_state2"
)

var (
	// default bns update interval 10s
	DF_BNS_UPDATE_INTERVAL = 60 * time.Second
)

// counters for redis
var (
	REDIS_CONN      = "REDIS_CONN"
	REDIS_CONN_FAIL = "REDIS_CONN_FAIL"

	REDIS_SET         = "REDIS_SET"
	REDIS_SET_FAIL    = "REDIS_SET_FAIL"
	REDIS_EXPIRE_FAIL = "REDIS_EXPIRE_FAIL"

	REDIS_GET      = "REDIS_GET"
	REDIS_GET_FAIL = "REDIS_GET_FAIL"
	REDIS_GET_MISS = "REDIS_GET_MISS"
	REDIS_GET_HIT  = "REDIS_GET_HIT"

	REDIS_INCR      = "REDIS_INCR"
	REDIS_INCR_FAIL = "REDIS_INCR_FAIL"

	REDIS_DECR      = "REDIS_DECR"
	REDIS_DECR_FAIL = "REDIS_DECR_FAIL"

	REDIS_GET_BNS_INSTANCE_ERR = "REDIS_GET_BNS_INSTANCE_ERR"
	REDIS_NO_BNS_INSTANCE      = "REDIS_NO_BNS_INSTANCE"
	REDIS_BNS_INSTANCE_CHANGED = "REDIS_BNS_INSTANCE_CHANGED"
)

type RedisClient struct {
	Servers     []string     // tcp address for redis servers
	serversLock sync.RWMutex // lock for servers
	bnsClient   *bns.Client  // bns client

	ConnectTimeout time.Duration // connect timeout (ms)
	ReadTimeout    time.Duration // read timeout (ms)
	WriteTimeout   time.Duration // write timeout (ms)

	pool      *redis.Pool  // connection pool to redis server
	poolLock  sync.RWMutex // lock for pool
	MaxIdle   int          // max idle conenctions in pool
	MaxActive int          // max active connections in pool
	Wait      bool         // if pool meet MaxActive limit, and Wait is true, wait for a connection return to pool

	state     *module_state2.State       // state for redis client
	delay     *delay_counter.DelayRecent // delay counter for reids
	connDelay *delay_counter.DelayRecent // delay counter for connect to redis
}

// NewRedisClient(): create a new redisClient with bns support
// Notice:
//    - if resolve bns error, c.Servers will be empty.
// Params:
//    - bnsName: string, bns name of redis server
//    - maxIdle: int, max idle connections in connection pool
//    - ct: int, connect redis server timeout, in ms
//    - rt: int, read redis server timeout, in ms
//    - wt: int, write redis server timeout, in ms
//    - state: *module_state2.State, module state
// Returns:
//    - *redisClient: a new redis client
func NewRedisClient(bnsName string, maxIdle int, ct, rt, wt int,
	state *module_state2.State) *RedisClient {

	// default maxActive 0, means no limit
	maxActive := 0
	wait := false
	return newRedisClient(bnsName, maxIdle, maxActive, wait, ct, rt, wt, state)
}

// NewRedisClient2(): create a new redisClient with bns support
// Notice:
//    - if resolve bns error, c.Servers will be empty.
// Params:
//    - bnsName: string, bns name of redis server
//    - maxIdle: int, max idle connections in connection pool
//    - maxActive: int, max active connections in connection pool
//    - wait: bool, if wait is true and pool at the maxActive limit,
//                  command waits for a connection return to the pool
//    - ct: int, connect redis server timeout, in ms
//    - rt: int, read redis server timeout, in ms
//    - wt: int, write redis server timeout, in ms
//    - state: *module_state2.State, module state
// Returns:
//    - *redisClient: a new redis client
func NewRedisClient2(bnsName string, maxIdle, maxActive int, wait bool, ct, rt, wt int,
	state *module_state2.State) *RedisClient {

	return newRedisClient(bnsName, maxIdle, maxActive, wait, ct, rt, wt, state)
}

// newRedisClient(): create a new redisClient with bns support
// Notice:
//    - if resolve bns error, c.Servers will be empty.
// Params:
//    - bnsName: string, bns name of redis server
//    - maxIdle: int, max idle connections in connection pool
//    - maxActive: int, max active connections in connection pool
//    - wait: bool, if wait is true and pool at the maxActive limit,
//                  command waits for a connection return to the pool
//    - ct: int, connect redis server timeout, in ms
//    - rt: int, read redis server timeout, in ms
//    - wt: int, write redis server timeout, in ms
//    - state: *module_state2.State, module state
// Returns:
//    - *redisClient: a new redis client
func newRedisClient(bnsName string, maxIdle, maxActive int, wait bool, ct, rt, wt int,
	state *module_state2.State) *RedisClient {
	var err error
	// create RedisClient
	c := &RedisClient{
		bnsClient: bns.NewClient(),

		// timeout in ms
		ConnectTimeout: time.Duration(ct) * time.Millisecond,
		ReadTimeout:    time.Duration(rt) * time.Millisecond,
		WriteTimeout:   time.Duration(wt) * time.Millisecond,

		// max idle connection
		MaxIdle: maxIdle,

		// max active connection
		MaxActive: maxActive,
		Wait:      wait,

		// module state
		state: state,
	}

	// if resolve bns error, c.Servers will be empty
	if c.Servers, err = bns.GetAddr(c.bnsClient, bnsName); err != nil {
		log.Logger.Warn("get instance for %s err %s", bnsName, err.Error())
	}
	go c.checkServerInstance(bnsName) // goroutine to update bns

	// set connection pool
	c.pool = &redis.Pool{
		MaxIdle:   c.MaxIdle,
		MaxActive: c.MaxActive,
		Wait:      c.Wait,
		Dial:      c.dial,
	}

	return c
}

// ActiveConnNum returns the num of active connextions
func (c *RedisClient) ActiveConnNum() int {
	// get connection pool
	c.poolLock.RLock()
	pool := c.pool
	c.poolLock.RUnlock()

	return pool.ActiveCount()
}

// set delay counter for redisClient
func (c *RedisClient) SetDelay(delayCounter *delay_counter.DelayRecent) {
	c.delay = delayCounter
}

// set conn delay counter for redisClient
func (c *RedisClient) SetConnDelay(delayCounter *delay_counter.DelayRecent) {
	c.connDelay = delayCounter
}

// dial choose a random server from c.Servers and connect
func (c *RedisClient) dial() (redis.Conn, error) {
	c.state.Inc(REDIS_CONN, 1)

	// choose a random server
	c.serversLock.RLock()
	if len(c.Servers) == 0 {
		c.serversLock.RUnlock()
		return nil, fmt.Errorf("no available connnection in pool")
	}
	server := c.Servers[rand.Intn(len(c.Servers))]
	c.serversLock.RUnlock()

	// create connection to server
	conn, err := redis.DialTimeout("tcp", server, c.ConnectTimeout, c.ReadTimeout, c.WriteTimeout)
	if err != nil {
		// add conn failed counter
		c.state.Inc(REDIS_CONN_FAIL, 1)

		return nil, err
	}

	return conn, nil
}

// Put(): save key:value to redis server, and set expire time
// Params:
//    - key: string
//    - value: []byte
//    - expire: int, expire time in second
// Returns:
//    - nil, if success, otherwise return error
//save sessionState to session cache
func (c *RedisClient) Put(key string, value []byte, expire int) (err error) {
	c.state.Inc(REDIS_SET, 1)

	// get a connection
	conn := c.GetConn()
	defer conn.Close()

	procStart := time.Now()
	// send set & expire cmd
	conn.Send("SET", key, value)
	conn.Send("EXPIRE", key, expire)
	conn.Flush()
	if _, err = conn.Receive(); err != nil {
		c.state.Inc(REDIS_SET_FAIL, 1)
		return err
	}
	if _, err = conn.Receive(); err != nil {
		c.state.Inc(REDIS_SET_FAIL, 1)
		return err
	}

	// monitor delay
	if c.delay != nil {
		c.delay.AddBySub(procStart, time.Now())
	}

	return nil
}

// get value from redis
func (c *RedisClient) Get(key string) (interface{}, bool) {
	c.state.Inc(REDIS_GET, 1)

	// get connection from pool
	conn := c.GetConn()
	defer conn.Close()

	procStart := time.Now()
	// get session state from redis
	value, err := conn.Do("GET", key)
	if err != nil {
		if err != redis.ErrNil {
			c.state.Inc(REDIS_GET_FAIL, 1)
		} else {
			c.state.Inc(REDIS_GET_MISS, 1)
		}
		return nil, false
	}

	// monitor delay
	if c.delay != nil {
		c.delay.AddBySub(procStart, time.Now())
	}

	c.state.Inc(REDIS_GET_HIT, 1)
	return value, true
}

// incr key to redis
func (c *RedisClient) Incr(key string, expire int) (int64, error) {
	c.state.Inc(REDIS_INCR, 1)

	// get connection from pool
	conn := c.GetConn()
	defer conn.Close()

	procStart := time.Now()
	// send incr & expire cmd
	conn.Send("INCR", key)
	conn.Send("EXPIRE", key, expire)
	conn.Flush()
	// get result from incr cmd
	count, err := redis.Int64(conn.Receive())
	if err != nil {
		c.state.Inc(REDIS_INCR_FAIL, 1)

		return count, err
	}

	// get result from expire cmd
	if _, err = conn.Receive(); err != nil {
		c.state.Inc(REDIS_EXPIRE_FAIL, 1)

		return count, err
	}

	// monitor delay
	if c.delay != nil {
		c.delay.AddBySub(procStart, time.Now())
	}

	return count, nil
}

// PIncr incr keys in pipeline mode
func (c *RedisClient) PIncr(keyList []string, expireList []int) ([]int64, error) {
	var err error
	var count int64
	countList := make([]int64, len(keyList), len(keyList))

	// get connection from pool
	conn := c.GetConn()
	defer conn.Close()

	procStart := time.Now()

	// send by pipeline
	for i := range keyList {
		c.state.Inc(REDIS_INCR, 1)
		// send incr cmd
		if err = conn.Send("INCR", keyList[i]); err != nil {
			c.state.Inc("REDIS_SEND_FAIL", 1)
			return countList, err
		}

		// send expire cmd
		if err = conn.Send("EXPIRE", keyList[i], expireList[i]); err != nil {
			c.state.Inc("REDIS_SEND_FAIL", 1)
			return countList, err
		}
	}

	// flush
	if err = conn.Flush(); err != nil {
		c.state.Inc("REDIS_FLUSH_FAIL", 1)
		return countList, err
	}

	// receive values
	for i := range keyList {
		// get result from incr cmd
		if count, err = redis.Int64(conn.Receive()); err != nil {
			c.state.Inc(REDIS_INCR_FAIL, 1)
			return countList, err
		}

		// read expire cmd
		if _, err = conn.Receive(); err != nil {
			c.state.Inc(REDIS_EXPIRE_FAIL, 1)
			return countList, err
		}

		// append to countList
		countList[i] = count
	}

	// monitor delay
	if c.delay != nil {
		c.delay.AddBySub(procStart, time.Now())
	}

	return countList, err
}

// decr key to redis
func (c *RedisClient) Decr(key string) (int64, error) {
	c.state.Inc(REDIS_DECR, 1)

	// get connection from pool
	conn := c.GetConn()
	defer conn.Close()

	procStart := time.Now()
	// send decr cmd
	conn.Send("DECR", key)
	conn.Flush()
	// get result from decr cmd
	count, err := redis.Int64(conn.Receive())
	if err != nil {
		c.state.Inc(REDIS_DECR_FAIL, 1)

		return count, err
	}

	// monitor delay
	if c.delay != nil {
		c.delay.AddBySub(procStart, time.Now())
	}

	return count, nil
}

func (c *RedisClient) UpdateServers(servers []string) {
	c.serversLock.Lock()
	c.Servers = servers
	c.serversLock.Unlock()
}

func (c *RedisClient) UpdatePool(pool *redis.Pool) *redis.Pool {
	c.poolLock.RLock()
	oldPool := c.pool
	c.pool = pool
	c.poolLock.RUnlock()

	return oldPool
}

// get a connection from connection pool
func (c *RedisClient) GetConn() redis.Conn {
	procStart := time.Now()

	// get connection pool
	c.poolLock.RLock()
	pool := c.pool
	c.poolLock.RUnlock()

	// get connection from pool
	conn := pool.Get()

	// monitor delay
	if c.connDelay != nil {
		c.connDelay.AddBySub(procStart, time.Now())
	}

	return conn
}

// update bns
func (c *RedisClient) checkServerInstance(name string) {
	for {
		time.Sleep(DF_BNS_UPDATE_INTERVAL)

		// check addresses of redis servers
		servers, err := bns.GetAddr(c.bnsClient, name)
		if err != nil {
			c.state.Inc(REDIS_GET_BNS_INSTANCE_ERR, 1)
			continue
		}
		if len(servers) == 0 {
			c.state.Inc(REDIS_NO_BNS_INSTANCE, 1)
			continue
		}
		if reflect.DeepEqual(servers, c.Servers) {
			continue
		}

		// update addresses of redis servers
		c.UpdateServers(servers)

		// counter bns instance changed
		c.state.Inc(REDIS_BNS_INSTANCE_CHANGED, 1)

		// update connection pool
		pool := &redis.Pool{
			MaxIdle:   c.MaxIdle,
			MaxActive: c.MaxActive,
			Wait:      c.Wait,
			Dial:      c.dial,
		}
		oldPool := c.UpdatePool(pool)
		oldPool.Close()
	}
}
