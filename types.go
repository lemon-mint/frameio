package frameio

import (
	"io"

	"github.com/lemon-mint/frameio/bufiopool"
)

var BufioPool *bufiopool.Pool = bufiopool.New(0, 0)

// FrameWriter write frame to io.Writer
// 		w: io.Writer  should be buffered
type FrameWriter struct {
	W io.Writer
}

// ReadFrame read frame from io.Reader
// 		r: io.Reader should be buffered
type FrameReader struct {
	R io.Reader
}

// Create a new FrameWriter
// 		w: io.Writer should be buffered
func NewFrameWriter(w io.Writer) FrameWriter {
	return FrameWriter{w}
}

// Create a new FrameReader
// 		r: io.Reader should be buffered
func NewFrameReader(r io.Reader) FrameReader {
	return FrameReader{r}
}
