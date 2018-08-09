/* byte_pool.go - manage byte slice pool for hash_set */
/*
modification history
--------------------
2014/8/21, by zhangjiyang01@baidu.com, create
2014/9/15, by zhangjiyang01@baidu.com, modify
           - move error handle from BytePool to hashSet
2016/10/12, by zhangjiyang01@baidu.com, move from hashSet
*/
/*
DESCRIPTION implement byte slice pool

*/
package byte_pool

import "fmt"

type BytePool struct {
	buf         []byte
	length      []uint32
	maxElemSize int // max length of element
	maxElemNum  int // max element num
}

/*
* NewBytePool - create a new BytePool
*
* PARAMS:
*   - elemNum: int, the max element num of BytePool
*   - maxElemSize: int, the max length of each element
*
* RETURNS:
*   - a pointer point to the BytePool
 */
func NewBytePool(elemNum int, maxElemSize int) *BytePool {
	pool := new(BytePool)
	pool.buf = make([]byte, elemNum*maxElemSize)
	pool.length = make([]uint32, elemNum)
	pool.maxElemSize = maxElemSize
	pool.maxElemNum = elemNum

	return pool
}

/*
* Set - set the index node of BytePool with key
*
* PARAMS:
*   - index: index of the byte Pool
*   - key: []byte key
 */
func (pool *BytePool) Set(index int32, key []byte) error {
	if int(index) >= pool.maxElemNum {
		return fmt.Errorf("index out of range %d %d", index, pool.maxElemNum)
	}

	if len(key) > pool.maxElemSize {
		return fmt.Errorf("elemSize large than maxSize %d %d", len(key), pool.maxElemSize)
	}

	start := int(index) * pool.maxElemSize
	copy(pool.buf[start:], key)

	pool.length[index] = uint32(len(key))

	return nil
}

/*
* Get the byte slice
*
* PARAMS:
*   - index: int, index of the BytePool
*
* RETURNS:
*   - key: []byte type store in the BytePool
 */
func (pool *BytePool) Get(index int32) []byte {
	start := int(index) * pool.maxElemSize
	end := start + int(pool.length[index])

	return pool.buf[start:end]
}

/* get the space allocate for each element */
func (pool *BytePool) MaxElemSize() int {
	return pool.maxElemSize
}
