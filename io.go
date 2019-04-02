package dutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

//const (
//	O_RDONLY int = syscall.O_RDONLY // 只读打开文件和os.Open()同义
//	O_WRONLY int = syscall.O_WRONLY // 只写打开文件
//	O_RDWR   int = syscall.O_RDWR   // 读写方式打开文件
//	O_APPEND int = syscall.O_APPEND // 当写的时候使用追加模式到文件末尾
//	O_CREATE int = syscall.O_CREAT  // 如果文件不存在，此案创建
//	O_EXCL   int = syscall.O_EXCL   // 和O_CREATE一起使用, 只有当文件不存在时才创建
//	O_SYNC   int = syscall.O_SYNC   // 以同步I/O方式打开文件，直接写入硬盘.
//	O_TRUNC  int = syscall.O_TRUNC  // 如果可以的话，当打开文件时先清空文件
//)

func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

//读取json文件并反序列化
func ReadJsonFileAndUnmarshal(filePath string, i interface{}) error {
	data, err := ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &i)
}

func writeFile(filePath string) (fileObj *os.File, err error) {
	fileObj, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("1", err)
			//os.MkdirAll(filePath, os.ModePerm)
			fileObj, err = os.Create(filePath)
			if err != nil {
				fmt.Println("2", err)
				return nil, err
			}
		} else {
			fmt.Println("3", err)
			return nil, err
		}
	}
	return
}

func WriteString(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), os.ModePerm)
	//fileObj, err := writeFile(filePath)
	//if err != nil {
	//	return err
	//}
	//defer fileObj.Close()
	//
	//_, err = fileObj.WriteString(content)
	//if err != nil {
	//	return err
	//} else {
	//	return nil
	//}
}

func WriteByte(filePath string, data []byte) error {
	fileObj, err := writeFile(filePath)
	if err != nil {
		return err
	}
	defer fileObj.Close()

	_, err = fileObj.Write(data)
	if err != nil {
		return err
	} else {
		return nil
	}
}
