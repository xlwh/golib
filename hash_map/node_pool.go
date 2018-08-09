/* node_pool.go - manage node pool for hash_map */
/*
modification history
--------------------
2016/10/11, by zhangjiyang01@baidu.com, modify
           - use uint32/int32 instead of int
*/
/*
DESCRIPTION

*/
package hash_map

import (
	"bytes"
	"fmt"
)

import (
	"www.baidu.com/golang-lib/byte_pool"
)

/* hash node */
type hashNode struct {
	next int32 // link to the next node
}

/* a list of hash node */
type nodePool struct {
	array []hashNode //node array

	freeNode int32 // manage the freeNode of nodePool
	capacity int   // capacity of nodePool
	length   int   // length of nodePool

	keyPool *byte_pool.BytePool // reference to []byte pool, store key
	valPool *byte_pool.BytePool // store value
}

/*
 * create a new nodePool
 *
 * PARAMS:
 *  - capacity: capacity
 *  - keySize: max size of each key
 *  - valSize: max size of each value
 *
 * RETURNS:
 *  - pointer to nodePool
 */
func newNodePool(capacity int, keySize, valSize int) *nodePool {
	np := new(nodePool)

	// make and init node array
	np.array = make([]hashNode, capacity)
	for i := 0; i < capacity-1; i += 1 {
		np.array[i].next = int32(i + 1) // link to the next node
	}
	np.array[capacity-1].next = -1 //intial value == -1, means end of the list

	np.freeNode = 0 //free node start from 0
	np.capacity = capacity
	np.length = 0

	np.keyPool = byte_pool.NewBytePool(capacity, keySize)
	np.valPool = byte_pool.NewBytePool(capacity, valSize)

	return np
}

/*
 * add
 *  - add key into the list starting from head
 *  - return the new headNode
 *
 * PARAMS:
 *  - head: first node of the list
 *  - key/val: []byte type
 *
 * RETURNS:
 *  - (newHead, nil), success, new headNode of the list
 *  - (-1, error), if fail
 */
func (np *nodePool) add(head int32, key, val []byte) (int32, error) {
	// get a bucket from freeNode List
	node, err := np.getFreeNode()
	if err != nil {
		return -1, err
	}

	np.array[node].next = head

	//set the node with key
	np.keyPool.Set(node, key)
	np.valPool.Set(node, val)

	np.length += 1
	return node, nil
}

/*
 * del
 *  - remove the key([]byte) in the given list
 *  - return the new head of the list
 *
 * PARAMS:
 *  - head: int, the first node of the list
 *  - key: []byte, the key need to be del
 *
 * RETURNS:
 *  - newHead int, the new head node of the list
 */
func (np *nodePool) del(head int32, key []byte) int32 {
	var newHead int32
	// check at the head of List
	if np.compare(key, head) == 0 {
		newHead = np.array[head].next
		np.recyleNode(head) //recyle the node
		return newHead
	}

	// check at the list
	pindex := head
	for {
		index := np.array[pindex].next
		if index == -1 {
			break
		}
		if np.compare(key, index) == 0 {
			np.array[pindex].next = np.array[index].next
			np.recyleNode(index) //recyle the node
			return head
		}
		pindex = index
	}
	return head
}

/* del the node, add the node into freeNode list */
func (np *nodePool) recyleNode(node int32) {
	index := np.freeNode
	np.freeNode = node
	np.array[node].next = index
	np.length -= 1
}

/* check if the key exist in the list */
func (np *nodePool) exist(head int32, key []byte) bool {
	for index := head; index != -1; index = np.array[index].next {
		if np.compare(key, index) == 0 {
			return true
		}
	}
	return false
}

func (np *nodePool) search(head int32, key []byte) ([]byte, bool) {
	for index := head; index != -1; index = np.array[index].next {
		if np.compare(key, index) == 0 {
			return np.elementVal(index), true
		}
	}
	return nil, false
}

/* get a free node from freeNode list */
func (np *nodePool) getFreeNode() (int32, error) {
	if np.freeNode == -1 {
		return -1, fmt.Errorf("NodePool: no more node to use")
	}

	// return freeNode and make freeNode = freeNode.next
	node := np.freeNode
	np.freeNode = np.array[node].next
	np.array[node].next = -1

	return node, nil
}

/* get node num in use of nodePool */
func (np *nodePool) elemNum() int {
	return np.length
}

/* check if the node Pool is full */
func (np *nodePool) full() bool {
	return np.length >= np.capacity
}

/* compare the given key with index node */
func (np *nodePool) compare(key []byte, i int32) int {
	item := np.elementKey(i)
	return bytes.Compare(key, item)
}

/* get the element key of the giving index*/
func (np *nodePool) elementKey(i int32) []byte {
	return np.keyPool.Get(i)
}

/* get the element val of the giving index*/
func (np *nodePool) elementVal(i int32) []byte {
	return np.valPool.Get(i)
}

/* get the space allocate for each element */
func (np *nodePool) keySize() int {
	return np.keyPool.MaxElemSize()
}

func (np *nodePool) valSize() int {
	return np.valPool.MaxElemSize()
}

/* check whether the key is legal for the map */
func (np *nodePool) validateKey(key []byte) error {
	if len(key) <= np.keySize() {
		return nil
	}
	return fmt.Errorf("element len[%d] > bucketSize[%d]", len(key), np.keySize())
}

/* check whether the key is legal for the map */
func (np *nodePool) validateVal(val []byte) error {
	if len(val) <= np.valSize() {
		return nil
	}
	return fmt.Errorf("element len[%d] > bucketSize[%d]", len(val), np.valSize())
}

/* check whether key/val is legal */
func (np *nodePool) validate(key, val []byte) error {
	err := np.validateKey(key)
	if err != nil {
		return err
	}

	return np.validateVal(val)
}
