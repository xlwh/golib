/* zlib.go - compress/decompress for compress/zlib  */
/*
modification history
--------------------
2016/05/30, by Lei Hong, Create
*/
/*
DESCRIPTION
This encapsulates compress/zlib:
1. to provide advantage
2. to avoid some obscure problems

Usage:
    import "www.baidu.com/golang-lib/compress"

    data := "the data you need to compress"
    cData, err := compress.ZlibCompress(data)
    if err != nil {
        ...
    }

    newData, err := compress.ZlibDecompress(cData)
    if err != nil {
        ...
    }

    // newData == data
*/
package compress

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
)

// compress data with zlib
func ZlibCompress(data []byte) ([]byte, error) {
	// new buffer and writer
	var b bytes.Buffer
	w := zlib.NewWriter(&b)

	// write data
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("zlib.Write(): %s", err.Error())
	}

	// important, it must be closed
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("zlib.Writer.Close(): %s", err.Error())
	}

	return b.Bytes(), nil
}

// decompress data with zlib
func ZlibDecompress(data []byte) ([]byte, error) {
	// new reader
	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("zlib.NewReader(): %s", err.Error())
	}

	// important, use ioutil.ReadAll, r.Read may not read all with a max buff
	dData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll() from zlib.reader: %s", err.Error())
	}

	// must close
	r.Close()

	return dData, nil
}
