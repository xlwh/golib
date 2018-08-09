/* Copyright 2017 Baidu Inc. All Rights Reserved. */
/* decrypt_baiduid.go - decrypt baiduid cookie key
/*
modification history
--------------------
2017/5/7, by Li Bingyi, create
*/
/*
DESCRIPTION

The descrypt algorithm refers http://gitlab.baidu.com/nginx/baidu-usertrack-module.

*/

package baidu_id

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
)

const DECODED_SIZE = BAIDUID_STR_LEN / 2

var newFunc = func() interface{} {
	return make([]byte, DECODED_SIZE)
}
var pool = sync.Pool{New: newFunc}

// Init cipher config, key must be 8 bytes len.
func Init(key string) (cipher.Block, error) {
	return des.NewCipher([]byte(key))
}

/*
* DecryptBaidudid- decrypt baiduid and get the clientip and unix time when generating it.
*
* PARAMS:
*   - block: cipher block initialized by des key.
*   - baiduidStr: baiduid cookie key, excluding the ":FG=1" tail.
*
* RETURNS:
*   uint32: the clientip addr when generating the baiduid, as big-endian byte order.
*   uint32: the unix time when generating the baiduid.
*   error: err info.
 */
func DecryptBaiduid(block cipher.Block, baiduidStr string) (uint32, uint32, error) {
	// check cipher len
	if len(baiduidStr) != BAIDUID_STR_LEN {
		return 0, 0, fmt.Errorf("baiduidStr len err")
	}

	// hex decode cipher
	decoded, err := hex.DecodeString(baiduidStr)
	if err != nil {
		return 0, 0, fmt.Errorf("hex.DecodeString failed")
	}

	// decrypt cipher text to plain text
	bs := block.BlockSize()
	s := pool.Get().([]byte)
	defer pool.Put(s)
	plain := s
	for len(decoded) > 0 {
		block.Decrypt(s, decoded[:bs])
		s = s[bs:]
		decoded = decoded[bs:]
	}

	// verify the checksum.
	// plain text is composed by:
	// 0-3 bytes: clientip when generating baiduid, as big-endian byte order.
	// 4-7 bytes: unix time when generating baiduid.
	// 8-11 bytes: random number.
	// 12-15 bytes: :checksum of 0-12 bytes.
	sum := uint32(0)
	for i := 0; i < 12; i++ {
		sum += uint32(plain[i])
	}
	realChecksum := sum ^ 0xFFFF

	expectedChecksum := binary.LittleEndian.Uint32(plain[12:16])
	if realChecksum != expectedChecksum {
		return 0, 0, fmt.Errorf("verify checksum failed")
	}

	clientip := binary.BigEndian.Uint32(plain[0:4])
	unixTime := binary.LittleEndian.Uint32(plain[4:8])
	return clientip, unixTime, nil
}
