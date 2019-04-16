package dutil

import (
	"log"
	"os"
	"path"
)

//log.New(logFile, "[debug]", log.Ldate|log.Ltime|log.Llongfile)
// 第一个参数为输出io，可以是文件也可以是实现了该接口的对象，此处为日志文件；
// 第二个参数为自定义前缀；第三个参数为输出日志的格式选项，可多选组合
// 第三个参数可选如下：
/*
   Ldate         = 1             // 日期：2009/01/23
   Ltime         = 2             // 时间：01:23:23
   Lmicroseconds = 4             // 微秒分辨率：01:23:23.123123（用于增强Ltime位）
   Llongfile     = 8             // 文件全路径名+行号： /a/b/c/d.go:23
   Lshortfile    = 16            // 文件无路径名+行号：d.go:23（会覆盖掉Llongfile）
   LstdFlags     = Ldate | Ltime // 标准logger的初始值
*/

type LogType int

const (
	InfoLog  = LogType(1)
	DebugLog = LogType(2)
	ErrorLog = LogType(3)
)

type DLogger struct {
}

func NewLogger(basePath string, fileName string) *log.Logger {
	//var logFile io.Writer
	os.MkdirAll(basePath, os.ModePerm)

	// 第三个参数为文件权限，请参考linux文件权限，664在这里为八进制，代表：rw-rw-r--
	logFile, err := os.OpenFile(path.Join(basePath, fileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	} else {
		debugLog := log.New(logFile, "[debug]", log.Ldate|log.Ltime|log.Llongfile)
		return debugLog
	}

	return nil
}

func write(t LogType, format string, args ...interface{}) {

}
