package bufiopool

import (
	"bufio"
	"io"
	"sync"
)

const DefaultBufSize = 4096

type Pool struct {
	rpool sync.Pool
	wpool sync.Pool

	size int
}

type dummyRW struct {
}

func (d dummyRW) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (d dummyRW) Write(p []byte) (n int, err error) {
	return 0, nil
}

var dummyRWInstance = &dummyRW{}

func New(rsize, wsize int) *Pool {
	if rsize == 0 {
		rsize = DefaultBufSize
	}
	if wsize == 0 {
		wsize = DefaultBufSize
	}
	return &Pool{
		rpool: sync.Pool{
			New: func() interface{} {
				return bufio.NewReaderSize(dummyRWInstance, rsize)
			},
		},
		wpool: sync.Pool{
			New: func() interface{} {
				return bufio.NewWriterSize(dummyRWInstance, wsize)
			},
		},
		size: rsize,
	}
}

func (p *Pool) GetReader(r io.Reader) *bufio.Reader {
	br := p.rpool.Get().(*bufio.Reader)
	br.Reset(r)
	return br
}

func (p *Pool) PutReader(br *bufio.Reader) {
	br.Reset(dummyRWInstance)
	p.rpool.Put(br)
}

func (p *Pool) GetWriter(w io.Writer) *bufio.Writer {
	bw := p.wpool.Get().(*bufio.Writer)
	bw.Reset(w)
	return bw
}

func (p *Pool) PutWriter(bw *bufio.Writer) {
	bw.Reset(dummyRWInstance)
	p.wpool.Put(bw)
}
