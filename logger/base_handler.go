package logger

import (
	"fmt"
	"io"
	"log"
)

//定义封装接口
type Handler interface {
	SetOutput(w io.Writer)
	Output(calldepth int, s string) error
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})

	//Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})

	Flags() int
	SetFlags(flag int)
	Prefix() string
	SetPrefix(prefix string)
	close()
}

/**
===================
LogHandler struct是对标准库log.Logger的封装
===================
**/

type LogHandler struct {
	lg *log.Logger
}

func (l *LogHandler) SetOutput(w io.Writer) {
	l.lg.SetOutput(w)
}

func (l *LogHandler) Output(calldepth int, s string) error {
	return l.lg.Output(calldepth, s)
}

func (l *LogHandler) Printf(format string, v ...interface{}) {
	l.lg.Printf(format, v...)
}

func (l *LogHandler) Print(v ...interface{}) {
	l.lg.Print(v...)
}

func (l *LogHandler) Println(v ...interface{}) {
	l.lg.Println(v...)
}

func (l *LogHandler) Fatal(v ...interface{}) {
	l.lg.Output(2, fmt.Sprint(v...))
}

func (l *LogHandler) Fatalf(format string, v ...interface{}) {
	l.lg.Output(2, fmt.Sprintf(format, v...))
}

func (l *LogHandler) Fatalln(v ...interface{}) {
	l.lg.Output(2, fmt.Sprintln(v...))
}

func (l *LogHandler) Flags() int {
	return l.lg.Flags()
}

func (l *LogHandler) SetFlags(flag int) {
	l.lg.SetFlags(flag)
}

func (l *LogHandler) Prefix() string {
	return l.lg.Prefix()
}

func (l *LogHandler) SetPrefix(prefix string) {
	l.lg.SetPrefix(prefix)
}

// func (l *LogHandler) Debug(v ...interface{}) {
// 	l.lg.Output(2, fmt.Sprintln(v)) //fmt.Sprintln("debug", v)
// }

func (l *LogHandler) Info(v ...interface{}) {
	l.lg.Output(2, fmt.Sprintln(v))
}

func (l *LogHandler) Warn(v ...interface{}) {
	l.lg.Output(2, fmt.Sprintln(v))
}

func (l *LogHandler) Error(v ...interface{}) {
	l.lg.Output(2, fmt.Sprintln(v))
}

func (l *LogHandler) close() {

}
