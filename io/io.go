package io

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

//读取json文件并反序列化
func DecodeJsonFromFile(i interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, i)
}

// 序列化对象并保存
func EncodeJsonToFile(i interface{}, filename string) error {
	data, err := json.Marshal(i)
	if err != nil {
		return err
	}
	_ = os.MkdirAll(path.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, data, os.ModePerm)
}

// csv
// 一种以逗号分割单元数据的文件，类似表格，但是很轻量。对于存储一些结构化的数据很有用。

// gob
// 无论纯文本还是csv文件的读写，所存储的数据文件是可以直接用文本工具打开的。
// 对于一些不希望被文件工具打开，需要将数据写成二进制。

func StoreGob(data interface{}, filename string) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}
	_ = os.MkdirAll(path.Dir(filename), os.ModePerm)
	return ioutil.WriteFile(filename, buffer.Bytes(), 0600)
}

func LoadGob(data interface{}, filename string) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	return dec.Decode(data)
}
