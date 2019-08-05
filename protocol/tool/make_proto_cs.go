package tool

import (
	"fmt"
	"os"
)

var message_template string = `syntax = \"proto2\";\npackage message;\n\nmessage %s_toS {}\n\nmessage %s_toC {}\n`

func GenProto(out_path string, name string) {

	fmt.Printf("gen_proto message ............\n")

	filename := fmt.Sprintf("%s/%s.proto", out_path, name)
	//检查文件是否存在，如果存在跳过不存在创建
	f, err := os.Open(filename)
	if nil != err && os.IsNotExist(err) {
		f, err = os.Create(filename)
		if nil == err {
			var content string
			content = fmt.Sprintf(message_template, name, name)
			_, err = f.WriteString(content)

			if nil != err {
				fmt.Printf("%s Write error:%s\n", name, err.Error())
			}

			f.Close()

		} else {
			fmt.Printf("%s Create error:%s\n", name, err.Error())
		}
	} else if nil != f {
		fmt.Printf("%s.proto exist skip\n", name)
		f.Close()
	}

}
