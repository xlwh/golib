/* byte_pool.go - manage byte slice pool for hash_set */
/*
modification history
--------------------
2014/8/21, by zhangjiyang01@baidu.com, create
2014/9/15, by zhangjiyang01@baidu.com, modify
           - move error handle from FixedBytePool to hashSet
2016/10/12, by zhangjiyang01@baidu.com, move from hashSet
*/
/*
DESCRIPTION implement byte slice pool for hash_set

*/
package byte_pool

import "fmt"

type FixedBytePool struct {
	buf        []byte
	elemSize   int // element length
	maxElemNum int // max element num
}

/*
* NewFixedBytePool - create a new FixedBytePool
*
* PARAMS:
*   - elemNum: int, the max element num of FixedBytePool
*   - elemSize: int, the max length of each element
*
* RETURNS:
*   - a pointer point to the FixedBytePool
 */
func NewFixedBytePool(elemNum int, elemSize int) *FixedBytePool {
	pool := new(FixedBytePool)
	pool.buf = make([]byte, elemNum*elemSize)
	pool.elemSize = elemSize
	pool.maxElemNum = elemNum

	return pool
}

/*
* Set - set the index node of FixedBytePool with key
*
* PARAMS:
*   - index: index of the byte Pool
*   - key: []byte key
 */
func (pool *FixedBytePool) Set(index int32, key []byte) error {
	if int(index) >= pool.maxElemNum {
		return fmt.Errorf("index out of range %d %d", index, pool.maxElemNum)
	}

	if len(key) != pool.elemSize {
		return fmt.Errorf("length must be %d while %d", pool.elemSize, len(key))
	}
	start := int(index) * pool.elemSize
	copy(pool.buf[start:], key)

	return nil
}

/*
* Get the byte slice of giving index and length
*
* PARAMS:
*   - index: int, index of the FixedBytePool
*
* RETURNS:
*   - key: []byte type store in the FixedBytePool
 */
func (pool *FixedBytePool) Get(index int32) []byte {
	start := int(index) * pool.elemSize
	end := start + pool.elemSize

	return pool.buf[start:end]
}

/* get the space allocate for each element */
func (pool *FixedBytePool) MaxElemSize() int {
	return pool.elemSize
}
