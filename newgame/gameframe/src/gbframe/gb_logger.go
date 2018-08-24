package gbframe

import (
	"io"
	"log"
	"os"
)

var (
	//	Trace   *log.Logger // 记录所有日志
	Logger_Info    *log.Logger // 重要的信息
	Logger_Warning *log.Logger // 需要注意的信息
	Logger_Error   *log.Logger // 致命错误
)
var LogFile string = "gbgame_log.log"

//const var Max int64 = 10

func init() {
	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	//	Trace = log.New(io.MultiWriter(file, os.Stderr), "[TRACE]: ", log.Ltime|log.Lshortfile)
	Logger_Info = log.New(io.MultiWriter(file, os.Stderr), "[Info]: ", log.Ltime|log.Lshortfile)
	Logger_Warning = log.New(io.MultiWriter(file, os.Stderr), "[Warning]: ", log.Ltime|log.Lshortfile)
	Logger_Error = log.New(io.MultiWriter(file, os.Stderr), "[Error]", log.Ltime|log.Lshortfile)
	//	fileInfo, _ := os.Stat(file)
	//	fileSize := fileInfo.Size()
}
