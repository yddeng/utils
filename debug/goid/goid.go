package goid

import (
	"bytes"
	"runtime"
	"strconv"
)

// Getgoid returns goroutine id
func Getgoid() int {
	buf := make([]byte, 20)
	buf = buf[:runtime.Stack(buf, false)]
	gid, _ := strconv.Atoi(string(bytes.Split(buf, []byte(" "))[1]))
	return gid
}
