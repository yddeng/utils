package bitmap

import "bytes"

type Bitmap struct {
	bits   []byte
	length int
}

func New(size uint16) *Bitmap {
	n := size / 8
	if size%8 != 0 {
		n++
	}
	return &Bitmap{bits: make([]byte, n)}
}

func (bm *Bitmap) Set(num uint64) bool {
	idx, bit := num/8, num%8
	if idx < uint64(len(bm.bits)) && (bm.bits[idx]&(1<<bit)) == 0 {
		bm.bits[idx] |= 1 << bit
		bm.length++
		return true
	}
	return false
}

func (bm *Bitmap) Clear(num uint64) bool {
	idx, bit := num/8, num%8
	if idx < uint64(len(bm.bits)) && (bm.bits[idx]&(1<<bit)) != 0 {
		bm.bits[idx] &^= 1 << bit
		bm.length--
		return true
	}
	return false
}

func (bm *Bitmap) Has(num uint64) bool {
	idx, bit := num/8, num%8
	if idx < uint64(len(bm.bits)) && (bm.bits[idx]&(1<<bit)) != 0 {
		return true
	}
	return false
}

func (bm *Bitmap) Len() int {
	return bm.length
}

func (bm *Bitmap) Cap() uint64 {
	return uint64(len(bm.bits) * 8)
}

func (bm *Bitmap) String() string {
	buffer := bytes.Buffer{}
	for _, bit := range bm.bits {
		for i := 0; i < 8; i++ {
			if bit&(1<<i) != 0 {
				buffer.WriteString("1")
			} else {
				buffer.WriteString("0")
			}
		}
	}
	return buffer.String()
}

func (bm *Bitmap) Copy(b *Bitmap) {
	min := len(b.bits)
	if len(bm.bits) < min {
		min = len(bm.bits)
	}

	for i := 0; i < min; i++ {
		bm.bits[i] = bm.bits[i] | b.bits[i]
	}
}
