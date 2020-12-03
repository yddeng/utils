package download

import (
	"fmt"
	"testing"
)

func TestHttpDownload(t *testing.T) {
	err := HttpDownload("https://xact02.baidupcs.com/file/a8c3ba9e6g6c6808b7d781cf9528afd2?bkt=en-1df98aeb722e81f571f813a79fcde5916400b97102fc95ce98e596901faaef291c35efaa005e9698&fid=898957251-250528-571572761257512&time=1606980275&sign=FDTAXUGERLQlBHSKfWqi-DCb740ccc5511e5e8fedcff06b081203-uX9VH7kJWDA7QkL6cz7zFHeugaA%3D&to=126&size=17822916&sta_dx=17822916&sta_cs=2&sta_ft=&sta_ct=4&sta_mt=4&fm2=MH%2CXian%2CAnywhere%2C%2Csichuan%2Cct&ctime=1605669941&mtime=1605669941&resv0=-1&resv1=0&resv2=rlim&resv3=5&resv4=17822916&vuk=898957251&iv=0&htype=&randtype=&newver=1&newfm=1&secfm=1&flow_ver=3&pkey=en-2bf37c85ec420e54cfca863b3077bd128ffaefbb7f1a956c4363af1dc1e8daf22a08acd4aba3f6b9&sl=76480590&expires=8h&rt=pr&r=402647875&vbdid=2733177546&fin=node_rank&rtype=1&dp-logid=7830872843946205686&dp-callid=0.1&hps=1&tsl=80&csl=80&fsl=-1&csign=W7WOF2KiOTKFIJODrYtGJDDh4y8%3D&so=0&ut=6&uter=4&serv=0&uc=3761371749&ti=75a34b3e95b22363e652326fe56329eda9f93f3c0a1455b6&hflag=30&adg=c_a5c234541b13b1999485c72418276942&reqlabel=250528_f_8602c82a76ff6e9bd62c9f792ed02809_-1_17d39c8ba7655d79f3634c805c333c53&by=themis",
		2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok")
}
