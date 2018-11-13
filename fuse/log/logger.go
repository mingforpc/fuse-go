package log

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 致命错误
)

func init() {
	Trace = log.New(os.Stdout, "TRACE: ", log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "Info: ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "Warning: ", log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr), "Error", log.Ltime|log.Lshortfile)
}
