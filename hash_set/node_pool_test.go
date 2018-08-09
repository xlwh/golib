/* hash_array_test.go - unit test file for hash_array.go */
/*
modification history
--------------------
2014/8/21, by zhangjiyang01@baidu.com, create
*/
/*
DESCRIPTION

*/
package hash_set

import (
	"bytes"
	"testing"
)

func TestAddItem(t *testing.T) {
	np := newNodePool(32, 32*32, false)

	Item := []byte("keyForTest")
	node, err := np.add(-1, Item)

	// normal case 1
	if err != nil {
		t.Errorf("err is not nil %s", err.Error())
	}
	if node != 0 {
		t.Errorf("node should be 0 %d", node)
	}
	if bytes.Compare(np.element(0), Item) != 0 {
		t.Error("element is wrong")
	}

	// normal case 2
	node, err = np.add(0, Item)
	if err != nil {
		t.Errorf("err is not nil %s", err.Error())
	}
	if node != 1 {
		t.Errorf("node should be 0 %d", node)
	}
	if bytes.Compare(np.element(1), Item) != 0 {
		t.Error("element is wrong")
	}
	if !np.exist(0, Item) {
		t.Error("should find in this list")
	}
	if np.compare(Item, 1) != 0 {
		t.Error("should find in this list")
	}

}

func TestDelItem(t *testing.T) {
	np := newNodePool(32, 32*32, false)

	Item := []byte("keyForTest")
	Item1 := []byte("keyForTest1")
	node, err := np.add(-1, Item)
	node, err = np.add(0, Item)
	node, err = np.add(1, Item1)
	_, _ = node, err
	if np.array[1].next != 0 {
		t.Error("1 should link to 0")
	}
	if np.array[0].next != -1 {
		t.Error("0 should link to -1")
	}

	// del at list
	if np.del(1, Item) != 0 {
		t.Error("del should return newhead")
	}

	// del head
	if np.del(0, Item1) != 0 {
		t.Error("del should return newhead")
	}
}

func TestGetFreeNode(t *testing.T) {
	np := newNodePool(32, 32*32, false)

	//case1
	node, err := np.getFreeNode()
	if node != 0 || err != nil {
		t.Error("get node error")
	}
	//case after recyleNode
	np.recyleNode(3)
	node, err = np.getFreeNode()
	if node != 3 || err != nil {
		t.Error("get node error")
	}

}
