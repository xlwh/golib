/* brotli_writer_test.go - unit test for brotli writer */
/*
modification history
--------------------
2016/10/26, by Sijie Yang, create
*/
/*
DESCRIPTION
*/
package brotli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

func testBrotliWriter(t *testing.T, quality int, orin []byte, comp []byte) {
	var buf bytes.Buffer
	brotliWriter := NewBrotliWriter(&buf, quality)

	_, err := brotliWriter.Write(orin)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	brotliWriter.Close()

	if !reflect.DeepEqual(comp, buf.Bytes()) {
		t.Errorf("compress error, expect: %v, actual: %v", comp, buf.Bytes())
	}
}

func TestBrotliWriter(t *testing.T) {
	orin, err := ioutil.ReadFile("testdata/quickfox")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	for quality := 0; quality < 12; quality++ {
		compFile := fmt.Sprintf("testdata/quickfox.c%d", quality)
		comp, err := ioutil.ReadFile(compFile)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		testBrotliWriter(t, quality, orin, comp)
	}
}
