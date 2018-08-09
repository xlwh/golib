/* baiduid_test.go - test for baiduid.go    */
/*
modification history
--------------------
2014/9/15, by Zhang Miao, create
2014/9/29, by Zhang Miao, move from go-bfe to golang-lib
*/
/*
DESCRIPTION
*/
package baidu_id

import (
    "bytes"
    "testing"
)

func TestBaiduIDStrToHex(t *testing.T) {
    hexOK := []byte{0x2d, 0x1b, 0x68, 0x28, 0x7c, 0x7a, 0x61, 0x2b,
                    0xa8, 0xe8, 0xe2, 0x2a, 0xf3, 0x17, 0x03, 0xcb}

    // case 1, succeed
    idStr := "2D1B68287C7A612BA8E8E22AF31703CB"
    idHex, hasFlag, err := BaiduIDStrToHex(idStr)
    if err != nil {
        t.Errorf("BaiduIDStrToHex() case 1: should not return err")
    }
    if hasFlag {
        t.Errorf("BaiduIDStrToHex() case 1: hasFlag should not be true")
    }
    if bytes.Compare(hexOK, idHex) != 0 {
        t.Errorf("BaiduIDStrToHex() case 1: idHex should be equal to hexOK")
    }

    // case 2, idStr is too short
    idStr = "2D1B68287C7A612BA8E8E22AF317"
    idHex, hasFlag, err = BaiduIDStrToHex(idStr)
    if err != ERR_BAIDUID_INVALID {
        t.Errorf("BaiduIDStrToHex() case 2: err should be ERR_BAIDUID_INVALID")
    }

    // case 3, invalid hex string
    idStr = "DU1B68287C7A612BA8E8E22AF31703CB"
    idHex, hasFlag, err = BaiduIDStrToHex(idStr)
    if err != ERR_BAIDUID_INVALID {
        t.Errorf("BaiduIDStrToHex() case 3: err should be ERR_BAIDUID_INVALID")
    }

    // case 4, invalid hex string
    idStr = "2D1B68287C7A612BA8E8E22AF31703CB:FG=1"
    idHex, hasFlag, err = BaiduIDStrToHex(idStr)
    if err != nil {
        t.Errorf("BaiduIDStrToHex() case 4: should not return err")
    }
    if !hasFlag {
        t.Errorf("BaiduIDStrToHex() case 4: hasFlag should be true")
    }
    if bytes.Compare(hexOK, idHex) != 0 {
        t.Errorf("BaiduIDStrToHex() case 4: idHex should be equal to hexOK")
    }

    // case 5, invalid hex string
    idStr = "2D1B68287C7A612BA8E8E22AF31703C1:FG=1"
    idHex, hasFlag, err = BaiduIDStrToHex(idStr)
    if err != nil {
        t.Errorf("BaiduIDStrToHex() case 4: should not return err")
    }
    if !hasFlag {
        t.Errorf("BaiduIDStrToHex() case 4: hasFlag should be true")
    }
}

func TestBaiduIDHexToStr(t *testing.T) {

    // normal case 1
    idStr := "2D1B68287C7A612BA8E8E22AF31703CB"
    idHex, hasFlag, err := BaiduIDStrToHex(idStr)
    if err != nil {
        t.Errorf("BaiduIDHexToStr() case 1: should not return err")
    }
    if hasFlag {
        t.Errorf("BaiduIDHexToStr() case 1: hasFlag should not be true")
    }

    retStr, err := BaiduIDHexToStr(idHex)
    if err != nil {
        t.Errorf("BaiduIDHexToStr(): case1: should not return err")
    }

    if retStr != idStr {
        t.Errorf("BaiduIDHexToStr(): case1: want[%s] && result is [%s]", idStr, retStr)
    }

    // case 2
    idHex = idHex[:1]
    retStr, err = BaiduIDHexToStr(idHex)
    if err == nil {
        t.Errorf("err shouldn't be nil")
    }
}

func TestBaiduIDTrim(t *testing.T) {
    id := "2D1B68287C7A612BA8E8E22AF31703CB"
    idStr := id + ":FG=1"

    if BaiduIDTrim(idStr) != id {
        t.Error("BaiduIDTrim() case 1: should equal")
    }
}
