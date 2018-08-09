/* baiduid.go - process baidu-id */
/*
modification history
--------------------
2014/9/15, by Zhang Miao, create
2014/9/29, by Zhang Miao, move from go-bfe to golang-lib
2014/11/14, by zhangjiyang, add BaiduIDHexToStr method
*/
/*
DESCRIPTION
*/
package baidu_id

import (
	"encoding/hex"
	"errors"
	"strings"
)

const (
	BAIDUID_STR_LEN = 32 // length of baiduID string
	BAIDUID_HEX_LEN = 16 // length of baiduID Hex
)

var (
	ERR_BAIDUID_INVALID = errors.New("BAIDUID_INVALID") // fail to convert to hex
)

/* BaiduIDStrToHex - convert baiduID from string to hex
 *
 * Params:
 *      - idStr: baiduID in string,
 *               e.g., "2D1B68287C7A612BA8E8E22AF31703CB:FG=1",
 *                     "2D1B68287C7A612BA8E8E22AF31703CB",
 *
 * Returns:
 *      (idHex, hasFlag, nil), if succeed
 *        - idHex: baiduID in hex (16 bytes)
 *        - hasFlag: whether have ":FG=1" in baiduID
 *
 *      (nil, hasFlag, error), if fail
 */
func BaiduIDStrToHex(idStr string) ([]byte, bool, error) {
	var hasFlag bool

	if strings.HasSuffix(idStr, ":FG=1") {
		hasFlag = true

		// remove suffix of ":FG=1"
		idStr = BaiduIDTrim(idStr)
	}

	// check length of baiduID string
	if len(idStr) != BAIDUID_STR_LEN {
		return nil, hasFlag, ERR_BAIDUID_INVALID
	}

	// convert to hex
	idHex, err := hex.DecodeString(idStr)
	if err != nil {
		return nil, hasFlag, ERR_BAIDUID_INVALID
	}

	return idHex, hasFlag, nil
}

/* BaiduIDHexToStr - convert baiduID from hex format to string
 *
 * Params:
 *      - idHex: baiduID in hex format
 *
 * Returns:
 *      (idStr, nil), if succeed
 *        - idStr: baiduID in String format (32 bytes)
 *      (nil, error), if fail
 */
func BaiduIDHexToStr(idHex []byte) (string, error) {
	// check length of baiduID hex
	if len(idHex) != BAIDUID_HEX_LEN {
		return "", ERR_BAIDUID_INVALID
	}

	// convert to String format
	idStr := hex.EncodeToString(idHex)

	// check length of BAIDUID_STR
	if len(idStr) != BAIDUID_STR_LEN {
		return "", ERR_BAIDUID_INVALID
	}
	return strings.ToUpper(idStr), nil
}

// Trim baiduid: remove :FG=1 suffix
func BaiduIDTrim(idStr string) string {
	return strings.TrimSuffix(idStr, ":FG=1")
}
