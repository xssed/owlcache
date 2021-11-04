package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

//动态切割File处理子类
//定义CutFileHandler struct
type CutFileHandler struct {
	LogHandler
	dir      string
	filename string
	indexNum int
	maxSize  int64
	logfile  *os.File
	mu       sync.Mutex
}

//创建CutFileHandler struct
func NewCutFileHandler(dir string, filename string, maxSize int64) *CutFileHandler {

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

	h := &CutFileHandler{
		LogHandler: LogHandler{l},
		dir:        dir,
		filename:   filename,
		indexNum:   0,
		maxSize:    maxSize,
	}

	//预检测日志文件大小，但是对于owlcache的日志来说创建时是当前的时间戳，所以这一段代码相对多余
	if h.checkSize() {
		h.renameFile()
	} else {
		h.mu.Lock()
		defer h.mu.Unlock()
		h.lg.SetOutput(logfile)
	}

	//每秒监控文件大小
	go func() {
		timer := time.NewTicker(1 * time.Second)
		//Ticker的另一种用法
		for {
			select {
			case <-timer.C:
				h.fileCheck()
			}
		}
	}()

	return h
}

//关闭日志文件资源
func (h *CutFileHandler) close() {
	if h.logfile != nil {
		h.logfile.Close()
	}
}

func (h *CutFileHandler) fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if h.checkSize() {
		h.renameFile()
	}
}

func (h *CutFileHandler) checkSize() bool {
	if fileSize(h.dir+"/"+h.filename) >= h.maxSize {
		return true
	}
	return false
}

//重命名备份文件日志
func (h *CutFileHandler) renameFile() {

	if h.logfile != nil {
		h.logfile.Close()
	}
	//去掉文件名中的.log字符串
	temp_filename := string(h.filename[:len(h.filename)-4])
	newpath := fmt.Sprintf("%s/%s.%d.log", h.dir, temp_filename, h.indexNum)
	//判断存在与否
	if isExist(newpath) {
		//删除这个文件
		os.Remove(newpath)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	filepath := h.dir + "/" + h.filename
	os.Rename(filepath, newpath) //重命名
	h.logfile, _ = os.Create(filepath)
	h.lg.SetOutput(h.logfile)

	h.indexNum += 1
}
