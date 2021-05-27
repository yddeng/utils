package bitmap

import "bytes"

type Bitmap struct {
	bits   []byte
	length int
}

// NewBitmap return [0, max-1]
func New(size uint16) *Bitmap {
	n := size / 8
	if size%8 != 0 {
		n++
	}
	return &Bitmap{bits: make([]byte, n)}
}

func (this *Bitmap) Set(num uint64) bool {
	idx, bit := num/8, num%8
	if idx < uint64(len(this.bits)) && (this.bits[idx]&(1<<bit)) == 0 {
		this.bits[idx] |= 1 << bit
		this.length++
		return true
	}
	return false
}

func (this *Bitmap) Clear(num uint64) bool {
	idx, bit := num/8, num%8
	if idx < uint64(len(this.bits)) && (this.bits[idx]&(1<<bit)) != 0 {
		this.bits[idx] &^= 1 << bit
		this.length--
		return true
	}
	return false
}

func (this *Bitmap) Has(num uint64) bool {
	idx, bit := num/8, num%8
	if idx < uint64(len(this.bits)) && (this.bits[idx]&(1<<bit)) != 0 {
		return true
	}
	return false
}

func (this *Bitmap) Len() int {
	return this.length
}

func (this *Bitmap) Range() (uint64, uint64) {
	return 0, uint64(len(this.bits)*8) - 1
}

func (this *Bitmap) String() string {
	buffer := bytes.Buffer{}
	for _, bit := range this.bits {
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

func (this *Bitmap) Copy(b *Bitmap) {
	min := len(b.bits)
	if len(this.bits) < min {
		min = len(this.bits)
	}

	for i := 0; i < min; i++ {
		this.bits[i] = this.bits[i] | b.bits[i]
	}
}
