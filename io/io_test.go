package io

import (
	"fmt"
	"testing"
)

type Info struct {
	ID   int
	Name string
}

func Test_Json(t *testing.T) {
	var ret = []Info{{ID: 1, Name: "11"}, {ID: 2, Name: "22"}}
	err := EncodeJsonToFile(&ret, "./file/center_config.json")
	fmt.Println(ret, err)

	var infos []Info
	err = DecodeJsonFromFile(&infos, "./file/center_config.json")
	fmt.Println(infos, err)
}

func Test_Gob(t *testing.T) {
	var ret = []Info{{ID: 1, Name: "11"}, {ID: 2, Name: "22"}}
	err := StoreGob(&ret, "./file/config.gob")
	fmt.Println(ret, err)

	var infos []Info
	err = LoadGob(&infos, "./file/config.gob")
	fmt.Println(infos, err)
}
