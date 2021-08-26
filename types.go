package frameio

import (
	"io"
)

type Frame struct {
	Payload []byte
}

type frameWriter struct {
	w io.Writer
}

type frameReader struct {
	r io.Reader
}
