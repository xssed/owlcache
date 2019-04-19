package log-pass

import (
	"fmt"
	"log"
	"os"
	"time"

	config "github.com/xssed/owlcache/config"
	"github.com/xssed/owlcache/logger"
)

//创建一个全局日志监听器
var OwlLoggerModel *OwlLog

func LogInit() {
	//执行步骤信息
	fmt.Println("owlcache  logger running...")
	LogRegister(config.OwlConfigModel)
}

//注册日志记录器
func LogRegister(owlconfig *config.OwlConfig) {

	//日志目录
	logFilePath := owlconfig.Logfile
	//日志文件
	formatLogFileName := "owlcache_" + time.Now().Format("20060102_150405") + ".log"

	OwlLoggerModel := logger.NewCutFileHandler(logFilePath, formatLogFileName, 7*1024*1024) //7M

	logger.SetHandlers(logger.Console, OwlLoggerModel)

	defer logger.Close()

	logger.SetFlags(0) //log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	//logger.SetLevel(logger.INFO)

	logger.Info("Info", "")
	logger.Error("Error", "")
	logger.Warn("Warn", "")
	//logger.Fatal("Fatal","")

}
