package buffer

import (
	"fmt"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	fmt.Println("start")

	buffer := &Buffer{}
	buffer.Write([]byte{0, 5, 4, 5, 6, 1, 3, 1})
	fmt.Println(buffer.Bytes(), buffer.Len())

	u16, err := buffer.ReadUint16BE()
	fmt.Println(u16, err)
	fmt.Println(buffer.Bytes(), buffer.Len())

	buffer.WriteUint16BE(56)
	fmt.Println(buffer.Bytes(), buffer.Len())

	test, err := buffer.ReadBytes(4)
	fmt.Println(test, err)
	fmt.Println(buffer.Bytes(), buffer.Len())

	buffer.ReadUint16BE()
	fmt.Println(buffer.Bytes())

	bt := buffer.Bytes()
	bt[0] = 255
	fmt.Println(buffer.Bytes())

	c, err := buffer.ReadByte()
	fmt.Println(c, err)
	fmt.Println(buffer.Bytes(), buffer.Len())

	buffer.WriteString("hello")
	fmt.Println(buffer.Bytes(), buffer.Len())
	_, _ = buffer.ReadByte()
	str, _ := buffer.ReadString(5)
	fmt.Println(str, buffer.Len())

}
