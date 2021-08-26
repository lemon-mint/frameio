package frameio

import (
	"encoding/binary"
)

const BLOCK_SIZE = 128

func (w *frameWriter) Write(b []byte) (err error) {
	var size [8]byte
	binary.BigEndian.PutUint64(size[:], uint64(len(b)))
	_, err = w.w.Write(size[:])
	if err != nil {
		return err
	}
	for len(b) > 0 {
		if len(b) < BLOCK_SIZE {
			return w.writeBlock(b)
		}
		err = w.writeBlock(b[:BLOCK_SIZE])
		if err != nil {
			return err
		}
		b = b[BLOCK_SIZE:]
	}
	return nil
}

func (w *frameWriter) writeBlock(b []byte) (err error) {
	length := byte(len(b))
	_, err = w.w.Write([]byte{length})
	if err != nil {
		return err
	}
	_, err = w.w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// var ErrSizeMismatch = errors.New("size mismatch")
//
// func (r *frameReader) ReadToBuffer(w io.Writer) (err error) {
// 	var sizeBytes [8]byte
// 	_, err = io.ReadFull(r.r, sizeBytes[:])
// 	if err != nil {
// 		return err
// 	}
// 	size := int(binary.BigEndian.Uint64(sizeBytes[:]))
// 	for size > 0 {
// 		n, err := r.readBlock(w)
// 		if err != nil {
// 			return err
// 		}
// 		if n != BLOCK_SIZE {
// 			break
// 		}
// 		size -= n
// 	}
//
// }
//
// var ErrInvalidBlockSize = errors.New("invalid block size")
//
// func (r *frameReader) readBlock(w io.Writer) (n int, err error) {
// 	var lenBytes [1]byte
// 	var buf [BLOCK_SIZE]byte
// 	_, err = io.ReadFull(r.r, lenBytes[:])
// 	if err != nil {
// 		return 0, err
// 	}
// 	if lenBytes[0] > BLOCK_SIZE {
// 		return 0, ErrInvalidBlockSize
// 	}
// 	_, err = io.ReadFull(r.r, buf[:lenBytes[0]])
// 	if err != nil {
// 		return 0, err
// 	}
// 	_, err = w.Write(buf[:lenBytes[0]])
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int(lenBytes[0]), nil
// }
//
