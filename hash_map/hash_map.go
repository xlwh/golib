/* hash_map.go - a map implementation */
/*
modification history
--------------------
2016/10/25, by zhangjiyang01@baidu.com, create
*/
/*
DESCRIPTION
    an implementation of map, use []byte pool in order to reduce time of gc
    for more details: http://wiki.babel.baidu.com/twiki/bin/view/Main/GO_BFE_DOCUMENT
Usage:
    import "www.baidu.com/golang-lib/hash_map"

    // elementNum: num of the element
    // elementSize: max size of each element
    hashMap := hash_map.NewHashMap(elementNum, maxKeySize, maxValSize, nil)

    hashMap.Add(key,val)

    val, ok := hashMap.Search(key)

*/
package hash_map

import (
	"fmt"
)

import (
	"github.com/murmur3"
)

/* in order to reduce the conflict of hash
 * hash array can be LOAD_FACTOR times larger than nodePool
 */
const (
	LOAD_FACTOR = 5
)

/* index table of hashhm */
type hashArray []int32

/* make a new hashArray and init it */
func newHashArray(indexSize int) hashArray {
	ha := make(hashArray, indexSize)
	for i := 0; i < indexSize; i += 1 {
		ha[i] = -1
	}

	return ha
}

type HashMap struct {
	ha     hashArray // hashArray, the index table for nodePool
	haSize int       // hashArray size

	np *nodePool // nodeMgr manage the elements of hashhm

	hashFunc func(key []byte) uint64 //function for hash
}

/*
* NewHashMap - create a newHashMap
*
* PARAMS:
*   - elemNum: max element num of hashhm
*   - keySize: maxSize of hashKey
*   - valSize: maxSize of val
*   - hashFunc: hash function
*
* RETURNS:
*  - (*HashMap, nil), if success
*  - (nil, error), if fail
 */
func NewHashMap(elemNum, keySize, valSize int, hashFunc func([]byte) uint64) (*HashMap, error) {
	if elemNum <= 0 || keySize <= 0 || valSize <= 0 {
		return nil, fmt.Errorf("elementNum/keySize/valSize must > 0")
	}

	hashMap := new(HashMap)

	/* hashArray is larger in order to reduce hash conflict */
	hashMap.haSize = elemNum * LOAD_FACTOR
	hashMap.ha = newHashArray(hashMap.haSize)

	hashMap.np = newNodePool(elemNum, keySize, valSize)

	/* if hashFunc is not given, use default murmur Hash */
	if hashFunc != nil {
		hashMap.hashFunc = hashFunc
	} else {
		hashMap.hashFunc = murmur3.Sum64
	}

	return hashMap, nil
}

/*
* Add - add an element into the map
*
* PARAMS:
*   - key: []byte, element of the map
*
* RETURNS:
*   - nil, if succeed
*   - error, if fail
 */
func (hm *HashMap) Add(key []byte, val []byte) error {
	// check the whether hashMap if full
	if hm.Full() {
		return fmt.Errorf("hashMap is full")
	}

	// validate hashKey
	if err := hm.np.validate(key, val); err != nil {
		return err
	}

	// 1. calculate the hash num
	hashNum := hm.hashFunc(key) % uint64(hm.haSize)

	// 2. check if the key slice exist
	if hm.exist(hashNum, key) {
		return nil
	}

	// 3. add the key into nodePool
	head := hm.ha[hashNum]
	newHead, err := hm.np.add(head, key, val)
	if err != nil {
		return err
	}

	// 4. point to the new list head node
	hm.ha[hashNum] = newHead

	return nil
}

/*
* Remove - remove an element from the hashMap
*
* PARAMS:
*   - key: []byte, element of the Map
*
* RETURNS:
*   - nil, if succeed
*   - error, if fail
 */
func (hm *HashMap) Remove(key []byte) error {
	// validate hashKey
	err := hm.np.validateKey(key)
	if err != nil {
		return err
	}
	//1. calculate the hash num
	hashNum := hm.hashFunc(key) % uint64(hm.haSize)

	//2. remove key from hashNode
	head := hm.ha[hashNum]
	if head == -1 {
		return nil
	}
	newHead := hm.np.del(head, key)

	//3. point to the new list head node
	hm.ha[hashNum] = newHead

	return nil
}

func (hm *HashMap) Search(key []byte) ([]byte, bool) {
	//validate hashKey
	err := hm.np.validateKey(key)
	if err != nil {
		return nil, false
	}

	hashNum := hm.hashFunc(key) % uint64(hm.haSize)
	return hm.search(hashNum, key)
}

func (hm *HashMap) search(hashNum uint64, key []byte) ([]byte, bool) {
	head := hm.ha[hashNum]
	return hm.np.search(head, key)
}

/* check if the element exist in HashMap*/
func (hm *HashMap) Exist(key []byte) bool {
	//validate hashKey
	err := hm.np.validateKey(key)
	if err != nil {
		return false
	}

	hashNum := hm.hashFunc(key) % uint64(hm.haSize)
	return hm.exist(hashNum, key)
}

/* check the []byte exist in the giving list head */
func (hm *HashMap) exist(hashNum uint64, key []byte) bool {
	head := hm.ha[hashNum]
	return hm.np.exist(head, key)
}

/* get elementNum of hashMap */
func (hm *HashMap) Len() int {
	return hm.np.elemNum()
}

/* check if the hashhm full or not */
func (hm *HashMap) Full() bool {
	return hm.np.full()
}
