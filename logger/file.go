package logger

import (
	"log"
	"os"
)

//File处理子类
//定义FileHandler struct
type FileHandler struct {
	LogHandler
	logfile *os.File //日志文件资源
}

//创建FileHandler struct
func NewFileHandler(dir string, filename string) *FileHandler {

	//日志存储目录校验
	temp_last := dir[len(dir)-1:]
	if temp_last != "/" {
		dir = dir + "/"
	}
	//创建目录
	if dir != "" {
		os.MkdirAll(dir, os.ModePerm)
	}
	//创建文件
	newfile := dir + filename
	logfile, _ := os.OpenFile(newfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	l := log.New(logfile, "", log.LstdFlags)
	return &FileHandler{
		LogHandler: LogHandler{l},
		logfile:    logfile,
	}
}

//关闭文件资源
func (h *FileHandler) close() {
	if h.logfile != nil {
		h.logfile.Close()
	}
}
