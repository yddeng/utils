package dutil_test

import (
	"dtest/dutil"
	"fmt"
	"testing"
)

func TestBuffer(t *testing.T) {
	b := dutil.NewBuffer()
	fmt.Println(b.GetUsedBuff())

	b.AppendUint16(uint16(322))
	//b.Read(bt)
	fmt.Println(b.GetUsedBuff())
	copy(b.GetUsableBuff(), []byte{1, 1})
	fmt.Println(b.GetUsedBuff())
	n, err := b.GetUint16()
	fmt.Println(b.GetUsedBuff())
	fmt.Println(n, err)
}

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
	name := "w/write.txt"
	content := "hello!"
	err := dutil.WriteString(name, content)
	fmt.Println(err)
}

func TestQueue(t *testing.T) {
	queue := dutil.NewQueue()
	for i := 1; i <= 15; i++ {
		queue.Add(i)
	}
	queue.P()
	fmt.Println(queue.GetFront())
	queue.P()
	fmt.Println(queue.Add(8))
	queue.P()
	fmt.Println(queue.Add(9))
	queue.P()
	fmt.Println(queue.Add(10))
	queue.P()
	fmt.Println(queue.GetAll())
	queue.P()
	fmt.Println(queue.GetFront())
	queue.P()
	queue.Add(1)
	fmt.Println(queue.GetFront())
	fmt.Println(queue.GetFront())
	queue.P()
}
