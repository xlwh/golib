/* zlib_test.go - test for zlib.go  */
/*
modification history
--------------------
2016/05/30, by Lei Hong, Create
*/
/*
DESCRIPTION

*/
package compress

import (
	"testing"
)

func TestZlibOperates(t *testing.T) {
	data := "This is a test string, yeah~~~"

	// compress
	cData, err := ZlibCompress([]byte(data))
	if err != nil {
		t.Errorf("ZlibCompress() error: %s", err.Error())
		return
	}

	// decompress
	newData, err := ZlibDecompress(cData)
	if err != nil {
		t.Errorf("ZlibDecompress() error: %s", err.Error())
		return
	}

	// compare
	newDataStr := string(newData)
	if data != newDataStr {
		t.Errorf("newData should to bu '%s' instead of '%s'", data, newDataStr)
	}
}
