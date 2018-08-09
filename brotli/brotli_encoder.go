/* brotli_encoder.go - golang brotli c-lib binding */
/*
modification history
--------------------
2016/10/26, by Sijie Yang, create
*/
/*
DESCRIPTION
*/
package brotli

/*
#include <string.h>

#include "brotli/encode.h"
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

const (
	DefualtOutputBufferSize = 64 * 1024 // 64 KB
)

var copyPool sync.Pool

func newByteBuf() []byte {
	if v := copyPool.Get(); v != nil {
		return v.([]byte)
	}
	return make([]byte, DefualtOutputBufferSize)
}

func putByteBuf(buf []byte) {
	copyPool.Put(buf)
}

// Errors which may be returned when encoding
var (
	errInputLargerThanBlockSize = errors.New("data copied to ring buffer larger than brotli encoder block size")
	errBrotliCompression        = errors.New("brotli compression error")
)

// Cgo utility function for brotli lib
func CBool(flag bool) C.BROTLI_BOOL {
	if flag {
		return C.BROTLI_TRUE
	} else {
		return C.BROTLI_FALSE
	}
}

func CArray(array []byte) *C.uint8_t {
	return (*C.uint8_t)(unsafe.Pointer(&array[0]))
}

// Brotli compressor
type brotliEncoder struct {
	c            *C.BrotliEncoderState
	outputBuffer []byte
}

// An instance can not be reused for multiple brotli streams.
func newBrotliEncoder(quality int) *brotliEncoder {
	cbp := C.BrotliEncoderCreateInstance(nil, nil, nil)
	C.BrotliEncoderSetParameter(cbp, C.BROTLI_PARAM_QUALITY, C.uint32_t(quality))
	bp := &brotliEncoder{c: cbp}
	bp.outputBuffer = newByteBuf()

	return bp
}

// The maximum input size that can be processed at once.
func (bp *brotliEncoder) getInputBlockSize() int {
	return int(C.BrotliEncoderInputBlockSize(bp.c))
}

// Copies the given input data to the internal ring buffer of the encoder.
// No processing of the data occurs at this time and this function can be
// called multiple times before calling WriteBrotliData() to process the
// accumulated input. At most getInputBlockSize() bytes of input data can be
// copied to the ring buffer, otherwise the next WriteBrotliData() will fail.
func (bp *brotliEncoder) copyInputToRingBuffer(input []byte) {
	C.BrotliEncoderCopyInputToRingBuffer(bp.c, C.size_t(len(input)), CArray(input))
}

// Processes the accumulated input data and returns the new output meta-block,
// or zero if no new output meta-block was created (in this case the processed
// input data is buffered internally).
// Returns ErrInputLargerThanBlockSize if more data was copied to the ring buffer
// than the block sized.
// If isLast or forceFlush is true, an output meta-block is always created
func (bp *brotliEncoder) writeBrotliData(isLast bool, forceFlush bool) ([]byte, error) {
	var outSize C.size_t
	var output *C.uint8_t

	success := C.BrotliEncoderWriteData(bp.c, CBool(isLast), CBool(forceFlush), &outSize, &output)
	if success == CBool(false) {
		return nil, errInputLargerThanBlockSize
	}

	// resize buffer if output is larger than we've anticipated
	if int(outSize) > cap(bp.outputBuffer) {
		bp.outputBuffer = make([]byte, int(outSize))
	}

	MemCopy(unsafe.Pointer(&bp.outputBuffer[0]), unsafe.Pointer(output), outSize)
	return bp.outputBuffer[:outSize], nil
}

func (bp *brotliEncoder) close() {
	if bp.c == nil {
		return
	}
	C.BrotliEncoderDestroyInstance(bp.c)
	putByteBuf(bp.outputBuffer)
	bp.c = nil
}
