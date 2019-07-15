package logger

import (
	"log"
	"os"
)

/**
===================
Console处理子类
===================
**/

//定义ConsoleHander struct
type ConsoleHander struct {
	LogHandler
}

//定义全局consolehander
var Console = NewConsoleHandler()

//创建全局consolehander
func NewConsoleHandler() *ConsoleHander {
	l := log.New(os.Stderr, "", log.LstdFlags)
	return &ConsoleHander{LogHandler: LogHandler{l}}
}
