package buffer_test

import (
	"fmt"
	dbuffer "github.com/yddeng/dutil/buffer"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buffer := dbuffer.NewBuffer(10)
	buffer.Write([]byte{0, 5, 4, 5, 6, 1, 3, 1})
	fmt.Println(buffer.Peek(), buffer.Len())

	u16, err := buffer.ReadUint16BE()
	fmt.Println(u16, err)
	fmt.Println(buffer.Peek(), buffer.Len())

	buffer.WriteUint16BE(56)
	fmt.Println(buffer.Peek(), buffer.Len())

	test, err := buffer.ReadBytes(4)
	fmt.Println(test, err)
	fmt.Println(buffer.Peek(), buffer.Len())

	buffer.ReadUint16BE()
	fmt.Println(test, err)
	fmt.Println(buffer.Peek())

	bt := buffer.Peek()
	bt[0] = 255
	fmt.Println(buffer.Peek())

	c, err := buffer.ReadByte()
	fmt.Println(c, err)
	fmt.Println(buffer.Peek(), buffer.Len())

	buffer.WriteString("hello")
	fmt.Println(buffer.Peek(), buffer.Len())
	_, _ = buffer.ReadByte()
	str, _ := buffer.ReadString(5)
	fmt.Println(str, buffer.Len())

}
