package buffer

import (
	"encoding/binary"
	"errors"
	"io"
	"unsafe"
)

var ErrNotEnough = errors.New("buffer.Buffer: read length but not enough")
var errNegativeRead = errors.New("buffer.Buffer: reader returned negative count from Read")

const smallBufferSize = 64
const maxInt = int(^uint(0) >> 1)

type Buffer struct {
	buf        []byte
	roff, woff int //读、写偏移
}

// 扩容 n ,如果剩余空间足够，则不扩容
func (b *Buffer) grow(n int) {
	if b.buf == nil || cap(b.buf) == 0 {
		b.buf = make([]byte, smallBufferSize)
	}

	if b.FreeLen() >= n {
		if b.Cap()-b.woff < n {
			b.Reset()
		}
		return
	}

	need := b.Len() + n
	newCap := 2 * b.Cap()
	for newCap < need {
		newCap *= 2
	}
	buf := make([]byte, newCap)
	b.woff = copy(buf, b.buf[b.roff:b.woff])
	b.roff = 0
	b.buf = buf

}

func (b *Buffer) Grow(n int) {
	if n < 0 {
		return
	}
	b.grow(n)
}

func (b *Buffer) empty() bool  { return b.Len() == 0 }
func (b *Buffer) Len() int     { return b.woff - b.roff }
func (b *Buffer) Cap() int     { return cap(b.buf) }
func (b *Buffer) FreeLen() int { return b.Cap() - b.Len() }

// 清空
func (b *Buffer) Clear() {
	b.buf = b.buf[:0]
	b.roff = 0
	b.woff = 0
}

// 清理已读数据，保留未读
func (b *Buffer) Reset() {
	if b.roff != 0 {
		copy(b.buf, b.buf[b.roff:b.woff])
		b.woff = b.woff - b.roff
		b.roff = 0
	}
}

func (b *Buffer) Bytes() []byte {
	return b.buf[b.roff:b.woff]
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.empty() {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}

	n = copy(p, b.buf[b.roff:b.woff])
	b.roff += n
	return n, nil
}

const MinRead = 512

// 只读一次，512 字节
func (b *Buffer) ReadFrom(reader io.Reader) (int64, error) {
	b.Grow(MinRead)
	n, e := reader.Read(b.buf[b.woff:])
	if n < 0 {
		panic(errNegativeRead)
	}

	b.woff += n
	return int64(n), e
}

func (b *Buffer) ReadAllFrom(reader io.Reader) (n int64, e error) {
	for {
		m, err := b.ReadFrom(reader)
		if m > 0 {
			n += m
		}
		if err == io.EOF {
			return n, nil
		}
		if err != nil {
			return n, err
		}
	}
}

// 写入缓冲区
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.Grow(len(p))
	n = copy(b.buf[b.woff:], p)
	b.woff += n
	return
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	if nBytes := b.Len(); nBytes > 0 {
		n, err := w.Write(b.buf[b.roff:b.woff])
		if n > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		b.roff += n
		if err != nil {
			return int64(n), err
		}

		if n != nBytes {
			err = io.ErrShortWrite
			return int64(n), err
		}
	}
	return 0, nil
}

func (b *Buffer) WriteUint8BE(num uint8) {
	b.Write([]byte{num})
}

func (b *Buffer) ReadUint8BE() (uint8, error) {
	num, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (b *Buffer) WriteUint16BE(num uint16) {
	var bt = make([]byte, 2)
	binary.BigEndian.PutUint16(bt, num)
	b.Write(bt)
}

func (b *Buffer) ReadUint16BE() (uint16, error) {
	num, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(num), nil
}

func (b *Buffer) WriteUint32BE(num uint32) {
	var bt = make([]byte, 4)
	binary.BigEndian.PutUint32(bt, num)
	b.Write(bt)
}

func (b *Buffer) ReadUint32BE() (uint32, error) {
	num, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(num), nil
}

func (b *Buffer) WriteUint64BE(num uint64) {
	var bt = make([]byte, 8)
	binary.BigEndian.PutUint64(bt, num)
	b.Write(bt)
}

func (b *Buffer) ReadUint64BE() (uint64, error) {
	num, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(num), nil
}

func (b *Buffer) WriteBytes(data []byte) {
	b.Write(data)
}

//获取 n 长度的数据
func (b *Buffer) ReadBytes(n int) ([]byte, error) {
	m := b.Len()
	if m < n {
		return nil, ErrNotEnough
	}

	var data = make([]byte, n)
	copy(data, b.buf[b.roff:b.roff+n])
	b.roff += n

	return data, nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.Len() == 0 {
		return 0, ErrNotEnough
	}

	ret := b.buf[b.roff]
	b.roff++

	return ret, nil
}

func (b *Buffer) WriteByte(c byte) {
	b.Write([]byte{c})
}

func (b *Buffer) WriteString(str string) {
	data := *(*[]byte)(unsafe.Pointer(&str))
	b.Write(data)
}

//获取len长度的数据
func (b *Buffer) ReadString(n int) (string, error) {
	bytes, err := b.ReadBytes(n)
	if err != nil {
		return "", err
	}

	//不用拷贝
	ret := *(*string)(unsafe.Pointer(&bytes))
	return ret, nil
}

func NewBuffer(buf []byte) *Buffer     { return &Buffer{buf: buf, woff: len(buf)} }
func NewBufferWithCap(cap int) *Buffer { return &Buffer{buf: make([]byte, cap)} }
