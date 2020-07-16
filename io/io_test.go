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

func Test_Json(t *testing.T) {
	var ret = []Info{{ID: 1, Name: "11"}, {ID: 2, Name: "22"}}
	err := io.EncodeJsonToFile(&ret, "./file/config.json")
	fmt.Println(ret, err)

	var infos []Info
	err = io.DecodeJsonFromFile(&infos, "./file/config.json")
	fmt.Println(infos, err)
}

func Test_Gob(t *testing.T) {
	var ret = []Info{{ID: 1, Name: "11"}, {ID: 2, Name: "22"}}
	err := io.StoreGob(&ret, "./file/config.gob")
	fmt.Println(ret, err)

	var infos []Info
	err = io.LoadGob(&infos, "./file/config.gob")
	fmt.Println(infos, err)
}
