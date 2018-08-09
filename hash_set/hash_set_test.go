/* hash_set_test.go - unit test for hash_set */
/*
modification history
--------------------
2016/10/25, by zhangjiyang01@baidu.com, create
*/
/*

*/

package hash_set

import "testing"

const TEST_COUNT = 400

func Hash([]byte) uint64 {
	num := 1
	return uint64(num)
}

func TestHashSet(t *testing.T) {
	table, err := NewHashSet(-1, -1, false, nil)
	if table != nil || err == nil {
		t.Error("wrong param, err should be nil")
	}

	table, err = NewHashSet(TEST_COUNT, 32, false, Hash)
	if err != nil {
		t.Error("err")
	}
	item := []byte("5728B5B85A6B1865E43F36DBB5F995EF")
	item1 := []byte("5728B5B85A6B1865E43F36DBB5F995EE")

	// test add item
	if table.Add(item) != nil {
		t.Error("add item should success")
	}

	if table.Len() != 1 {
		t.Error("length of hashTable should be 1")
	}
	if !table.Exist(item) {
		t.Error("should exist")
	}

	if table.Exist(item1) {
		t.Error("should exist")
	}

	table.Add(item1)
	if !table.Exist(item) {
		t.Error("should exist")
	}

	if table.Len() != 2 {
		t.Error("length of hashTable should be 2")
	}

	if !table.Exist(item1) {
		t.Error("should exist")
	}

	// test remove item
	err = table.Remove(item)
	if err != nil {
		t.Error("should remove success")
	}
	if table.Len() != 1 {
		t.Error("length of hashTable should be 1")
	}
	if table.Exist(item) {
		t.Error("should not exist")
	}

	if !table.Exist(item1) {
		t.Error("should exist")
	}

	// test remove wrong case
	wrongItem := []byte("5728B5B85A6B1865E43F36DBB5F995EFFFFFFFFF")
	err = table.Remove(wrongItem)
	if err == nil {
		t.Error("err should not be nil")
	}

	if table.Len() != 1 {
		t.Error("length of hashTable should be 1")
	}

}

func TestHashSetWithFixedSize(t *testing.T) {
	table, err := NewHashSet(-1, -1, true, nil)
	if table != nil || err == nil {
		t.Error("wrong param, err should be nil")
	}
	table, err = NewHashSet(TEST_COUNT, 32, true, Hash)
	if err != nil {
		t.Error("err")
	}
	item := []byte("5728B5B85A6B1865E43F36DBB5F995EF")
	item1 := []byte("5728B5B85A6B1865E43F36DBB5F995EE")

	// test add item
	if table.Add(item) != nil {
		t.Error("add item should success")
	}

	if table.Len() != 1 {
		t.Error("length of hashTable should be 1")
	}

	if !table.Exist(item) {
		t.Error("should exist")
	}

	if table.Exist(item1) {
		t.Error("should exist")
	}

	table.Add(item1)
	if !table.Exist(item) {
		t.Error("should exist")
	}

	if table.Len() != 2 {
		t.Error("length of hashTable should be 2")
	}

	if !table.Exist(item1) {
		t.Error("should exist")
	}

	// test remove item
	err = table.Remove(item)
	if err != nil {
		t.Error("should remove success")
	}
	if table.Exist(item) {
		t.Error("should not exist")
	}

	if table.Len() != 1 {
		t.Error("length of hashTable should be 1")
	}

	if !table.Exist(item1) {
		t.Error("should exist")
	}

	// test remove wrong case
	wrongItem := []byte("5728B5B85A6B1865E43F36DBB5F995EFFFFFFFFF")
	err = table.Remove(wrongItem)
	if err == nil {
		t.Error("err should not be nil")
	}

}
