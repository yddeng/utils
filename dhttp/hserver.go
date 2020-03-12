package dhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpServer struct {
	handler    *http.ServeMux
	listenAddr string
}

func NewHttpServer(addr string) *HttpServer {
	s := new(HttpServer)
	s.handler = http.NewServeMux()
	s.listenAddr = addr

	return s
}

func (s *HttpServer) Register(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.handler.HandleFunc(pattern, handler)
}

func (s *HttpServer) Listen() error {
	return http.ListenAndServe(s.listenAddr, s.handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//fmt.Fprintf(w, "Hi, I love you %s", html.EscapeString(r.URL.Path[1:]))
		r.ParseForm()                    //解析参数，默认是不会解析的
		fmt.Println("method:", r.Method) //获取请求的方法

		fmt.Println("username", r.Form["username"])
		fmt.Println("password", r.Form["password"])

		for k, v := range r.Form {
			fmt.Print("key:", k, "; ")
			fmt.Println("val:", strings.Join(v, ""))
		}
	} else if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)

		//未知类型的推荐处理方法

		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		for k, v := range m {
			switch vv := v.(type) {
			case string:
				fmt.Println(k, "is string", vv)
			case int:
				fmt.Println(k, "is int", vv)
			case float64:
				fmt.Println(k, "is float64", vv)
			case []interface{}:
				fmt.Println(k, "is an array:")
				for i, u := range vv {
					fmt.Println(i, u)
				}
			default:
				fmt.Println(k, "is of a type I don't know how to handle")
			}
		}

		//结构已知，解析到结构体
		/*
			var s Serverslice
			json.Unmarshal([]byte(result), &s)

			fmt.Println(s.ServersID)

			for i := 0; i < len(s.Servers); i++ {
				fmt.Println(s.Servers[i].ServerName)
				fmt.Println(s.Servers[i].ServerIP)
			}
		*/
	}
}
