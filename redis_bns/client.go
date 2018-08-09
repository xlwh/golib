/* client.go - interface for redis client */
/*
modification history
--------------------
2016/07/14, by Zhang Jiyang create
*/
package redis_bns

import "github.com/garyburd/redigo/redis"

// Client: redis client interface
type Client interface {
	GetConn() redis.Conn
	Put(key string, value []byte, expire int) error
	Get(key string) (interface{}, bool)
	Incr(key string, expire int) (int64, error)
	Decr(key string) (int64, error)
	PIncr([]string, []int) ([]int64, error)
}
