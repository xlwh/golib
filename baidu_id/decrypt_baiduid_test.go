/* Copyright 2017 Baidu Inc. All rights Reserved. */
/* decrypt_baiduid_test.go - test file of decrypt_baiduid.go
/*
modification history
--------------------
2017/5/7, by Li Bingyi, create
*/
/*
DESCRIPTION
*/

package baidu_id

import (
	"encoding/hex"
	"testing"
)

func TestInit(t *testing.T) {
	_, err := Init("test_key")
	if err != nil {
		t.Errorf("Init should succ")
	}

	_, err = Init("test_keytest_key")
	if err == nil {
		t.Errorf("Init should fail")
	}
}

func TestDecryptBaiduid(t *testing.T) {
	// init
	block, err := Init("test_key")
	if err != nil {
		t.Errorf("Init err: %s", err)
	}
	bs := block.BlockSize()

	// cipher len error
	str := "A354C8741EB92F894C6692FFCC8D20F00"
	_, _, err = DecryptBaiduid(block, str)
	if err == nil {
		t.Errorf("len str: %d != %d.", len(str), BAIDUID_STR_LEN)
	}

	// hex decode cipher error
	str = "A354C8741EB92F894C6692FFCC8D20F@"
	_, _, err = DecryptBaiduid(block, str)
	if err == nil {
		t.Errorf("hex.DecodeString() should fail")
	}

	// decrypt ok
	b := []byte{127, 0, 0, 1, 107, 160, 30, 89, 171, 155, 182, 5, 252, 251, 0, 0}
	s := make([]byte, len(b))
	cipher := s
	for len(b) > 0 {
		block.Encrypt(s, b[:bs])
		s = s[bs:]
		b = b[bs:]
	}

	str = hex.EncodeToString(cipher)
	_, _, err = DecryptBaiduid(block, str)
	if err != nil {
		t.Errorf("DecryptBaiduid should succ, err: %s", err)
	}

	// checksum err
	b = []byte{127, 0, 0, 1, 107, 160, 30, 89, 171, 155, 182, 5, 251, 251, 0, 0}
	s = make([]byte, len(b))
	cipher = s
	for len(b) > 0 {
		block.Encrypt(s, b[:bs])
		s = s[bs:]
		b = b[bs:]
	}

	str = hex.EncodeToString(cipher)
	_, _, err = DecryptBaiduid(block, str)
	if err == nil {
		t.Errorf("checksum different, DecryptBaiduid should fail")
	}
}
