package strutil

import (
	"bytes"
	"regexp"
	"strings"
	"unsafe"
)

// 拼接字符串
func MergeString(args ...string) string {
	buffer := bytes.Buffer{}
	for _, str := range args {
		buffer.WriteString(str)
	}
	return buffer.String()
}

// 删除多余的分隔符
func DeleteExtraSep(s, sep string) string {
	return ""
}

// 将字符串中多个连续空格转换为一个空格，返回转换结果
func DeleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "	", " ", -1)       //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}

// 根据左边和右边的内容，获取指定字符串两者中间的内容，如果没找到则返回空字符串
func GetBetweenStr(str string, start string, end string) string {
	copy := str
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start) // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	if str == copy {
		return ""
	} else {
		return str
	}
}

// 检查字符串是否为空字符串，返回bool
func IsEmpty(str string) bool {
	if str == "" || len(str) == 0 {
		return true
	}
	return false
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 将字符串按照空格分割成字符串数组
func Str2Slice(s string) []string {
	words := strings.Split(DeleteExtraSpace(s), " ")
	ret := make([]string, 0, len(words))
	for _, w := range words {
		if w != "" {
			ret = append(ret, w)
		}
	}
	return ret
}
