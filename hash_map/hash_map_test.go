package hash_map

import (
	"bytes"
	"testing"
)

const TEST_COUNT = 400

func Hash([]byte) uint64 {
	num := 1
	return uint64(num)
}

func TestHashSet(t *testing.T) {
	table, err := NewHashMap(-1, -1, 1, nil)
	if table != nil || err == nil {
		t.Error("wrong param, err should be nil")
	}

	table, err = NewHashMap(TEST_COUNT, 32, 32, Hash)
	if err != nil {
		t.Error("err")
	}
	item := []byte("5728B5B85A6B1865E43F36DBB5F995EF")
	item1 := []byte("5728B5B85A6B1865E43F36DBB5F995EE")

	// test add item
	if table.Add(item, item1) != nil {
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

	if val, ok := table.Search(item); !ok || bytes.Compare(val, item1) != 0 {
		t.Error("should get val")
	}

	table.Add(item1, item)
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
