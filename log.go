package dutil

import (
	"fmt"
	"log"
	"os"
	"time"
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
	InfoLog LogType = iota
	DebugLog
	ErrorLog
)

var LogLevel = [...]string{
	"[INFO]",
	"[DEBUG]",
	"[ERROR]",
}

type Logger struct {
	fPath  string
	logger *log.Logger
}

func NewLogger(basePath, fileName string) *Logger {
	return newLogger(basePath, fileName)
}

func newLogger(basePath, fileName string) *Logger {
	//var logFile io.Writer
	year, month, day := time.Now().Date()
	hour, min, sec := time.Now().Clock()

	dir := fmt.Sprintf("%s/%04d-%02d-%02d", basePath, year, month, day)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fPath := fmt.Sprintf("%s/%s.%02d.%02d.%02d.log", dir, fileName, hour, min, sec)

	// 第三个参数为文件权限，请参考linux文件权限，664在这里为八进制，代表：rw-rw-r--
	logFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	logger := log.New(logFile, "[DEBUG]", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		fPath:  fPath,
		logger: logger,
	}
}

func (l *Logger) print(t LogType, format string, v ...interface{}) {
	prefix := LogLevel[t]

	if l.logger.Prefix() != prefix {
		l.logger.SetPrefix(prefix)
	}

	if format == "" {
		l.logger.Println(v...)
	} else {
		l.logger.Printf(format, v...)
	}

}

func (l *Logger) Infoln(v ...interface{}) {
	l.print(InfoLog, "", v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.print(InfoLog, format, v...)
}

func (l *Logger) Debugln(v ...interface{}) {
	l.print(DebugLog, "", v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.print(DebugLog, format, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.print(ErrorLog, "", v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.print(ErrorLog, format, v...)
}
