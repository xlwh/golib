/* brotli_writer.go - brotli writer for golang */
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
	"io"
)

// BrotliWriter implements the io.Writer interface, compressing the stream
// to an output Writer using Brotli.
type BrotliWriter struct {
	encoder *brotliEncoder
	writer  io.Writer

	// amount of data already copied into ring buffer
	inRingBuffer int
}

// NewBrotliWriter instantiates a new BrotliWriter with the provided compression
// parameters and output Writer
func NewBrotliWriter(writer io.Writer, quality int) *BrotliWriter {
	return &BrotliWriter{
		encoder:      newBrotliEncoder(quality),
		writer:       writer,
		inRingBuffer: 0,
	}
}

func (w *BrotliWriter) Write(buffer []byte) (int, error) {
	comp := w.encoder
	blockSize := int(comp.getInputBlockSize())
	roomFor := blockSize - w.inRingBuffer
	copied := 0

	for len(buffer) >= roomFor {
		comp.copyInputToRingBuffer(buffer[:roomFor])
		copied += roomFor

		compressedData, err := comp.writeBrotliData(false, false)
		if err != nil {
			return copied, err
		}

		_, err = w.writer.Write(compressedData)
		if err != nil {
			return copied, err
		}

		w.inRingBuffer = 0
		buffer = buffer[roomFor:]
		roomFor = blockSize
	}

	remaining := len(buffer)
	if remaining > 0 {
		comp.copyInputToRingBuffer(buffer)
		w.inRingBuffer += remaining
		copied += remaining
	}

	return copied, nil
}

func (w *BrotliWriter) Flush() error {
	compressedData, err := w.encoder.writeBrotliData(false, true)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(compressedData)
	if err != nil {
		return err
	}

	return nil
}

// Close cleans up the resources used by the Brotli encoder for this
// stream. If the output buffer is an io.Closer, it will also be closed.
func (w *BrotliWriter) Close() error {
	compressedData, err := w.encoder.writeBrotliData(true, false)
	if err != nil {
		return err
	}
	w.encoder.close()

	_, err = w.writer.Write(compressedData)
	if err != nil {
		return err
	}

	if v, ok := w.writer.(io.Closer); ok {
		return v.Close()
	}

	return nil
}
