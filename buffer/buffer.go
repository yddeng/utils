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
const MinRead = 512

type Buffer struct {
	r, w  int
	buf   []byte
	cap   int
	empty bool
}

func (b *Buffer) Empty() bool {
	return b.empty
}

func (b *Buffer) Full() bool {
	if b.r == b.w && !b.empty {
		return true
	}
	return false
}

func (b *Buffer) Clear() {
	b.r = 0
	b.w = 0
	b.empty = true
}

func (b *Buffer) Len() int {
	if b.w > b.r {
		return b.w - b.r
	} else if b.w < b.r {
		return b.cap - b.r + b.w
	} else {
		if b.empty {
			return 0
		} else {
			return b.cap
		}
	}
}

func (b *Buffer) Cap() int  { return b.cap }
func (b *Buffer) Free() int { return b.Cap() - b.Len() }

func (b *Buffer) Grow(n int) {
	need := b.Len() + n
	newCap := 2 * b.Cap()
	for newCap < need {
		newCap *= 2
	}

	buf := make([]byte, newCap)
	var num int
	if !b.empty {
		if b.w > b.r {
			num = copy(buf, b.buf[b.r:b.w])
		} else {
			num = copy(buf, b.buf[b.r:])
			if b.w > 0 {
				num += copy(buf[num:], b.buf[:b.w])
			}
		}
	}
	b.cap = newCap
	b.buf = buf
	b.r = 0
	b.w = num % b.cap
}

func (b *Buffer) Read(buf []byte) (n int, err error) {
	plen := len(buf)
	if plen == 0 || b.empty {
		return
	}

	if b.w > b.r {
		n = copy(buf, b.buf[b.r:b.w])
		b.r = (b.r + n) % b.cap
	} else {
		n = copy(buf, b.buf[b.r:])
		b.r = (b.r + n) % b.cap
		if plen > n && b.w > 0 {
			b.r = copy(buf[n:], b.buf[:b.w])
			n += b.r
		}
	}

	if b.r == b.w {
		b.empty = true
	}
	return
}

func (b *Buffer) Reset(n int) {
	if n > b.Len() {
		n = b.Len()
	}

	if n > 0 {
		b.r = (b.r + n) % b.cap
		if b.r == b.w {
			b.empty = true
		}
	}
}

func (b *Buffer) Bytes() []byte {
	if !b.empty {
		buf := make([]byte, b.Len())
		if b.w > b.r {
			copy(buf, b.buf[b.r:b.w])
		} else {
			n := copy(buf, b.buf[b.r:])
			if b.w > 0 {
				copy(buf[n:], b.buf[:b.w])
			}
		}
		return buf
	}
	return nil
}

// 只读一次，512 字节
func (b *Buffer) ReadFrom(reader io.Reader) (int64, error) {
	free := b.Free()
	if free == 0 {
		return 0, nil
	}

	buf := make([]byte, free)
	n, err := reader.Read(buf)
	if n < 0 {
		panic(errNegativeRead)
	}

	n, _ = b.Write(buf[:n])
	return int64(n), err
}

func (b *Buffer) Write(buf []byte) (n int, err error) {
	plen := len(buf)
	if plen == 0 || (b.r == b.w && !b.empty) {
		return
	}

	if b.w < b.r {
		n = copy(b.buf[b.w:b.r], buf)
		b.w = (b.w + n) % b.cap
	} else {
		n = copy(b.buf[b.w:], buf)
		b.w = (b.w + n) % b.cap
		if plen > n && b.r > 0 {
			b.w = copy(b.buf[:b.r], buf[n:])
			n += b.w
		}
	}

	if n > 0 {
		b.empty = false
	}
	return
}

func (b *Buffer) WriteUint8BE(num uint8) (int, error) {
	if b.Free() < 1 {
		return 0, ErrNotEnough
	}
	return b.Write([]byte{num})
}

func (b *Buffer) ReadUint8BE() (uint8, error) {
	num, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (b *Buffer) WriteUint16BE(num uint16) (int, error) {
	if b.Free() < 2 {
		return 0, ErrNotEnough
	}
	var bt = make([]byte, 2)
	binary.BigEndian.PutUint16(bt, num)
	return b.Write(bt)
}

func (b *Buffer) ReadUint16BE() (uint16, error) {
	num, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(num), nil
}

func (b *Buffer) WriteUint32BE(num uint32) (int, error) {
	if b.Free() < 4 {
		return 0, ErrNotEnough
	}
	var bt = make([]byte, 4)
	binary.BigEndian.PutUint32(bt, num)
	return b.Write(bt)
}

func (b *Buffer) ReadUint32BE() (uint32, error) {
	num, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(num), nil
}

func (b *Buffer) WriteUint64BE(num uint64) (int, error) {
	if b.Free() < 8 {
		return 0, ErrNotEnough
	}
	var bt = make([]byte, 8)
	binary.BigEndian.PutUint64(bt, num)
	return b.Write(bt)
}

func (b *Buffer) ReadUint64BE() (uint64, error) {
	num, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(num), nil
}

func (b *Buffer) WriteBytes(data []byte) (int, error) {
	return b.Write(data)
}

//获取 n 长度的数据
func (b *Buffer) ReadBytes(n int) ([]byte, error) {
	m := b.Len()
	if m < n {
		return nil, ErrNotEnough
	}

	data := make([]byte, n)
	n, _ = b.Read(data)
	return data, nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.Empty() {
		return 0, ErrNotEnough
	}

	data := make([]byte, 1)
	b.Read(data)
	return data[0], nil
}

func (b *Buffer) WriteByte(c byte) (int, error) {
	return b.Write([]byte{c})
}

func (b *Buffer) WriteString(str string) (int, error) {
	data := *(*[]byte)(unsafe.Pointer(&str))
	return b.Write(data)
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

func NewBufferWithCap(size int) *Buffer {
	return &Buffer{buf: make([]byte, size), cap: size, empty: true}
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{buf: buf, cap: len(buf), empty: false}
}
