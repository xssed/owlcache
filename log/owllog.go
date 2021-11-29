package log

import (
	"time"

	"github.com/xssed/owlcache/logger"
)

type OwlLog struct {
	*logger.Logger
}

func NewOwlLog(logFilePath, formatLogFileName string) *OwlLog {

	log := logger.New()
	datestr := time.Now().Format("2006_01_02")
	logFilePath = logFilePath + "/" + datestr
	LoggerModel := logger.NewCutFileHandler(logFilePath, formatLogFileName, 5*1024*1024) //5M
	log.SetHandlers(logger.Console, LoggerModel)
	log.SetFlags(0) //log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	//log.SetLevel(logger.INFO)
	// defer log.Close()
	// log.Info("Info", "")
	// log.Error("Error", "")
	// log.Warn("Warn", "")
	// log.Fatal("Fatal", "")
	return &OwlLog{log}
}
