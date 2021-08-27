package frameio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/valyala/bytebufferpool"
)

const BLOCK_SIZE = 127

// Write one frame to the writer.
func (w *FrameWriter) Write(b []byte) (err error) {
	var size [8]byte
	binary.BigEndian.PutUint64(size[:], uint64(len(b)))
	_, err = w.W.Write(size[:])
	if err != nil {
		return err
	}
	for len(b) > 0 {
		if len(b) <= BLOCK_SIZE {
			return w.writeBlock(b, true)
		}
		err = w.writeBlock(b[:BLOCK_SIZE], false)
		if err != nil {
			return err
		}
		b = b[BLOCK_SIZE:]
	}
	return nil
}

// Write one block to the writer.
func (w *FrameWriter) writeBlock(b []byte, isLast bool) (err error) {
	length := byte(len(b))
	if isLast {
		length |= 0x80
	}
	_, err = w.W.Write([]byte{length})
	if err != nil {
		return err
	}
	_, err = w.W.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// Sum of all block sizes is not equal to TotalSize.
var ErrSizeMismatch = errors.New("size mismatch")

// Frame Size Exceeds MaxSize.
var ErrFrameTooLarge = errors.New("frame too large")

// Read one frame from the reader and write it to the buffer.
func (r *FrameReader) readToBuffer(w io.Writer, max int) (err error) {
	var sizeBytes [8]byte
	var payloadSize int
	_, err = io.ReadFull(r.R, sizeBytes[:])
	if err != nil {
		return err
	}
	size := int(binary.BigEndian.Uint64(sizeBytes[:]))
	if size > max && max != -1 {
		return ErrFrameTooLarge
	}
	for size > 0 {
		n, isLast, err := r.readBlock(w)
		if err != nil {
			return err
		}
		payloadSize += n
		size -= n
		if isLast {
			break
		}
		if payloadSize > max && max != -1 {
			return ErrFrameTooLarge
		}
	}
	if size != 0 {
		return ErrSizeMismatch
	}
	return nil
}

// Read one frame from the reader and write it to the buffer.
func (r *FrameReader) ReadToBuffer(w io.Writer) (err error) {
	return r.readToBuffer(w, -1)
}

// Read one frame from the reader and write it to the writer.
//
// If the frame is larger than maxSize, ErrFrameTooLarge is returned.
func (r *FrameReader) ReadToBufferMax(w io.Writer, maxSize int) (err error) {
	return r.readToBuffer(w, maxSize)
}

// Read one frame from the reader and call the callback function.
//
// Returned byte slice is not valid after the callback function returns. (due to reusing the buffer)
func (r *FrameReader) ReadCallback(cb func([]byte)) (err error) {
	buffer := bytebufferpool.Get()
	defer bytebufferpool.Put(buffer)
	err = r.readToBuffer(buffer, -1)
	if err != nil {
		return err
	}
	cb(buffer.B)
	return nil
}

// Read one frame from the reader and call the callback function.
//
// Returned byte slice is not valid after the callback function returns. (due to reusing the buffer)
//
// If the frame is larger than maxSize, ErrFrameTooLarge is returned.
func (r *FrameReader) ReadCallbackMax(cb func([]byte), maxSize int) (err error) {
	buffer := bytebufferpool.Get()
	defer bytebufferpool.Put(buffer)
	err = r.readToBuffer(buffer, maxSize)
	if err != nil {
		return err
	}
	cb(buffer.B)
	return nil
}

// Read one frame from the reader and return it as a byte array. (allocates memory)
func (r *FrameReader) Read() (data []byte, err error) {
	var buf bytes.Buffer
	err = r.readToBuffer(&buf, -1)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Read one frame from the reader and return it as a byte array. (allocates memory)
//
// If the frame is larger than maxSize, ErrFrameTooLarge is returned.
func (r *FrameReader) ReadMax(maxSize int) (data []byte, err error) {
	var buf bytes.Buffer
	err = r.readToBuffer(&buf, maxSize)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// BlockSize is bigger than the maximum BLOCK_SIZE.
var ErrInvalidBlockSize = errors.New("invalid block size")

// Read one block from the reader.
func (r *FrameReader) readBlock(w io.Writer) (n int, isLast bool, err error) {
	var lenBytes [1]byte
	var buf [BLOCK_SIZE]byte
	_, err = io.ReadFull(r.R, lenBytes[:])
	if err != nil {
		return 0, isLast, err
	}
	length := lenBytes[0] & 0x7F
	isLast = lenBytes[0]&0x80 != 0
	if length > BLOCK_SIZE {
		return 0, isLast, ErrInvalidBlockSize
	}
	_, err = io.ReadFull(r.R, buf[:length])
	if err != nil {
		return 0, isLast, err
	}
	_, err = w.Write(buf[:length])
	if err != nil {
		return 0, isLast, err
	}
	return int(length), isLast, nil
}
