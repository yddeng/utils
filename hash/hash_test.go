package hash

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	h := New()
	h.Add("abcdefg")
	h.Add("hijklmn")
	h.Add("bfghdfg")

	a, err := h.Get("0abcdefg")
	fmt.Println(a, err)

	c, err := h.GetN("0abcdefg", 3)
	fmt.Println(c, err)

}
