package io_test

import (
	"fmt"
	"github.com/yddeng/dutil/io"
	"testing"
)

type Info struct {
	ID   int
	Name string
}

func TestDecodeJsonFile(t *testing.T) {
	var infos = []Info{}
	err := io.DecodeJsonFile("./file/config.json", &infos)
	fmt.Println(infos, err)

	err = io.WriteString("./", "out.txt", "out")
	fmt.Println(err)
}
