package logger

import (
	"fmt"
	"os"
	"sync"
)

/**
===================
 logger
===================
**/

type Level int32

const (
	//DEBUG Level = iota
	INFO Level = iota
	WARN
	ERROR
)

type _Logger struct {
	handlers []Handler
	level    Level
	mu       sync.Mutex
}

var logger = &_Logger{
	handlers: []Handler{
		Console,
	},
	level: INFO,
}

func SetHandlers(handlers ...Handler) {
	logger.handlers = handlers
}

func SetFlags(flag int) {
	for i := range logger.handlers {
		logger.handlers[i].SetFlags(flag)
	}
}

func SetLevel(level Level) {
	logger.level = level
}

func Print(v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Print(v...)
	}
}

func Printf(format string, v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Printf(format, v...)
	}
}

func Println(v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Println(v...)
	}
}

func Fatal(v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Fatal(v...)
	}
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Fatalf(format, v...)
	}
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Fatalln(v...)
	}
	os.Exit(1)
}

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	for i := range logger.handlers {
		logger.handlers[i].Output(2, s)
	}
	panic(s)
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	for i := range logger.handlers {
		logger.handlers[i].Output(2, s)
	}
	panic(s)
}

func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	for i := range logger.handlers {
		logger.handlers[i].Output(2, s)
	}
	panic(s)
}

// func Debug(v ...interface{}) {
// 	if logger.level <= DEBUG {
// 		for i := range logger.handlers {
// 			logger.handlers[i].Debug(NewFormat(v))
// 		}
// 	}
// }

func Info(v ...interface{}) {
	if logger.level <= INFO {
		for i := range logger.handlers {

			logger.handlers[i].Info(NewFormat(v))
		}
	}
}

func Warn(v ...interface{}) {
	if logger.level <= WARN {
		for i := range logger.handlers {
			logger.handlers[i].Warn(NewFormat(v))
		}
	}
}

func Error(v ...interface{}) {
	if logger.level <= ERROR {
		for i := range logger.handlers {
			logger.handlers[i].Error(NewFormat(v))
		}
	}
}

func Close() {
	for i := range logger.handlers {
		logger.handlers[i].close()
	}
}
