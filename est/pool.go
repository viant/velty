package est

import (
	"sync"
)

type Pool struct {
	pool       sync.Pool
	bufferSize int
}

func (p *Pool) Borrow() *Buffer {
	buffer := p.pool.Get().(*Buffer)
	buffer.index = 0
	if len(buffer.buf) > p.bufferSize {
		buffer.buf = buffer.buf[:p.bufferSize]
	}
	return buffer
}

func (p *Pool) Put(bs *Buffer) {
	if len(bs.buf) > p.bufferSize {
		return
	}
	p.pool.Put(bs)
}

func NewPool(bufferSize int) *Pool {
	return &Pool{
		bufferSize: bufferSize,
		pool: sync.Pool{
			New: func() interface{} {
				return &Buffer{buf: make([]byte, bufferSize)}
			},
		},
	}
}
