package est

import (
	"fmt"
	"github.com/viant/velty/utils"
	"html"
	"strconv"
)

type Buffer struct {
	buf        []byte
	index      int
	poolSize   int
	escapeHTML bool
}

func (b *Buffer) AppendByte(bs byte) {
	if b.index+1 >= len(b.buf) {
		newBuffer := make([]byte, len(b.buf)+b.poolSize)
		b.buf = append(b.buf, newBuffer...)
	}
	b.buf[b.index] = bs
	b.index++
}

func (b *Buffer) AppendInt(v int) {
	b.growIfNeeded(65) // 64 int size and sign if < 0
	b.index += utils.AppendInt(b.buf[b.index:], int64(v), 10)
}

func (b *Buffer) AppendBool(v bool) {
	s := strconv.FormatBool(v)
	b.AppendString(s)
}

func (b *Buffer) AppendFloat(v float64) {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	b.AppendString(s)
}

func (b *Buffer) AppendString(s string) {
	if !b.escapeHTML {
		b.AppendStringWithoutEscaping(s)
		return
	}

	b.AppendStringWithoutEscaping(html.EscapeString(s))
}

func (b *Buffer) AppendStringWithoutEscaping(s string) {
	sLen := len(s)
	if sLen == 0 {
		return
	}
	b.growIfNeeded(sLen)
	copy(b.buf[b.index:], s)
	b.index += sLen
}

func (b *Buffer) growIfNeeded(sLen int) {
	if sLen > 10*1024*1024 {
		panic(fmt.Sprintf("to large to big %v\n", sLen))
	}
	if sLen+b.index >= len(b.buf) {
		size := len(b.buf) + b.poolSize
		if size < sLen {
			size = sLen
		}
		newBuffer := make([]byte, size)
		b.buf = append(b.buf, newBuffer...)
	}
}

func (b *Buffer) Reset() {
	b.index = 0
}

func (b *Buffer) Bytes() []byte {
	return b.buf[:b.index]
}

func (b *Buffer) String() string {
	return string(b.buf[:b.index])
}

func NewBuffer(size int, escape bool) *Buffer {
	return &Buffer{
		buf:        make([]byte, size),
		escapeHTML: escape,
	}
}
