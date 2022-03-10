package est

type Buffer struct {
	buf      []byte
	index    int
	poolSize int
}

func (b *Buffer) AppendBytes(bs []byte) {
	bsLen := len(bs)
	if bsLen == 0 {
		return
	}
	if bsLen+b.index >= len(b.buf) {
		size := b.poolSize
		if size < bsLen {
			size = bsLen
		}
		b.buf = append(b.buf, make([]byte, size)...)
	}
	copy(b.buf[b.index:], bs)
	b.index += bsLen
}
func (b *Buffer) AppendByte(bs byte) {
	if b.index+1 >= len(b.buf) {
		newBuffer := make([]byte, len(b.buf)+b.poolSize)
		b.buf = append(b.buf, newBuffer...)
	}
	b.buf[b.index] = bs
	b.index++
}
func (b *Buffer) AppendString(s string) {
	sLen := len(s)
	if sLen == 0 {
		return
	}
	if sLen+b.index >= len(b.buf) {
		size := len(b.buf) + b.poolSize
		if size < sLen {
			size = sLen
		}
		newBuffer := make([]byte, size)
		b.buf = append(b.buf, newBuffer...)
	}
	copy(b.buf[b.index:], s)
	b.index += sLen
}

func (b *Buffer) Reset() {
	b.index = 0
}

func (b *Buffer) Bytes() []byte {
	return b.buf[:b.index]
}
