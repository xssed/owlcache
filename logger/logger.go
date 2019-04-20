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

type Logger struct {
	handlers []Handler
	level    Level
	mu       sync.Mutex
}

func New() *Logger {
	return &Logger{
		handlers: []Handler{
			Console,
		},
		level: INFO,
	}
}

func (logger *Logger) SetHandlers(handlers ...Handler) {
	logger.handlers = handlers
}

func (logger *Logger) SetFlags(flag int) {
	for i := range logger.handlers {
		logger.handlers[i].SetFlags(flag)
	}
}

func (logger *Logger) SetLevel(level Level) {
	logger.level = level
}

func (logger *Logger) Print(v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Print(v...)
	}
}

func (logger *Logger) Printf(format string, v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Printf(format, v...)
	}
}

func (logger *Logger) Println(v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Println(v...)
	}
}

func (logger *Logger) Fatal(v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Fatal(v...)
	}
	os.Exit(1)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Fatalf(format, v...)
	}
	os.Exit(1)
}

func (logger *Logger) Fatalln(v ...interface{}) {
	for i := range logger.handlers {
		logger.handlers[i].Fatalln(v...)
	}
	os.Exit(1)
}

func (logger *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	for i := range logger.handlers {
		logger.handlers[i].Output(2, s)
	}
	panic(s)
}

func (logger *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	for i := range logger.handlers {
		logger.handlers[i].Output(2, s)
	}
	panic(s)
}

func (logger *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	for i := range logger.handlers {
		logger.handlers[i].Output(2, s)
	}
	panic(s)
}

// func (logger *Logger) Debug(v ...interface{}) {
// 	if logger.level <= DEBUG {
// 		for i := range logger.handlers {
// 			logger.handlers[i].Debug(NewFormat(v))
// 		}
// 	}
// }

func (logger *Logger) Info(v ...interface{}) {
	if logger.level <= INFO {
		for i := range logger.handlers {

			logger.handlers[i].Info(NewFormat(v))
		}
	}
}

func (logger *Logger) Warn(v ...interface{}) {
	if logger.level <= WARN {
		for i := range logger.handlers {
			logger.handlers[i].Warn(NewFormat(v))
		}
	}
}

func (logger *Logger) Error(v ...interface{}) {
	if logger.level <= ERROR {
		for i := range logger.handlers {
			logger.handlers[i].Error(NewFormat(v))
		}
	}
}

func (logger *Logger) Close() {
	for i := range logger.handlers {
		logger.handlers[i].close()
	}
}
