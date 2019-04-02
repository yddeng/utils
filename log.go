package dutil

import (
	"log"
	"os"
	"path"
)

type LogType int

const (
	DebugfT  = LogType(3)
	DebuglnT = LogType(4)
	ErrorfT  = LogType(5)
	ErrorlnT = LogType(6)
)

func NewLogger(basePath string, fileName string) *log.Logger {

	//var logFile io.Writer
	if nil == os.MkdirAll(basePath, os.ModePerm) {
		// 打开日志文件
		// 第二个参数为打开文件的模式，可选如下：
		/*
		   O_RDONLY // 只读模式打开文件
		   O_WRONLY // 只写模式打开文件
		   O_RDWR   // 读写模式打开文件
		   O_APPEND // 写操作时将数据附加到文件尾部
		   O_CREATE // 如果不存在将创建一个新文件
		   O_EXCL   // 和O_CREATE配合使用，文件必须不存在
		   O_SYNC   // 打开文件用于同步I/O
		   O_TRUNC  // 如果可能，打开时清空文件
		*/
		// 第三个参数为文件权限，请参考linux文件权限，664在这里为八进制，代表：rw-rw-r--
		context := path.Join(basePath, fileName) + ".log"
		logFile, err := os.OpenFile(context, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		} else {
			debugLog := log.New(logFile, "[debug]", log.Ldate|log.Ltime|log.Llongfile)
			return debugLog
		}

	}

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
	return nil
}

func write(t LogType, format string, args ...interface{}) {

}
