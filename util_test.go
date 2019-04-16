package dutil_test

import (
	"fmt"
	"github.com/tagDong/dutil"
	"testing"
)

func TestReadJsonFileAndUnmarshal(t *testing.T) {
	type Info struct {
		ID   int32
		Name string
	}
	var infos []Info
	err := dutil.ReadJsonFileAndUnmarshal("file/config.json", &infos)
	fmt.Println(infos, err)
}

func TestWriteFile(t *testing.T) {
	filePath := "w"
	name := "write.txt"
	content := "hello!"
	err := dutil.WriteByte(filePath, name, []byte(content))
	fmt.Println(err)
}
