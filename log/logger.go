package log

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/comail/colog"
	config "github.com/xssed/owlcache/config"
)

////创建一个全局日志监听器
//var OwlLoggerModel *OwlLogger

func LogInit() {
	//执行步骤信息
	fmt.Println("owlcache  logger running...")
	LogRegister(config.OwlConfigModel)
}

//注册日志记录器
func LogRegister(owlconfig *config.OwlConfig) {

	logFilePath := owlconfig.Logfile

	//日志存储目录校验
	temp_last := logFilePath[len(logFilePath)-1:]
	if temp_last != "/" {
		logFilePath = logFilePath + "/"
	}
	//fmt.Println(logFilePath) //输出日志存放目录

	formatLogFileName := "owlcache_" + time.Now().Format("20060102_150405") + ".log"
	if logFilePath != "" {
		os.MkdirAll(logFilePath, os.ModePerm)
	}
	logfile := logFilePath + formatLogFileName

	file, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		//panic(err)
		fmt.Println(err)
		fmt.Print("Failure to create a log file.Please check whether 'Logfile' is set correctly. now exit\n")
		os.Exit(0)
	}

	colog.Register()
	colog.SetOutput(file)
	colog.ParseFields(true)
	colog.SetFormatter(&colog.JSONFormatter{
		TimeFormat: time.RFC3339Nano,
		//Flag:       log.Lshortfile,
	})

}

func Debug(str string) {
	log.Print("debug: " + str)
}

func Info(str string) {
	//{"level":"info","time":"2015-08-16T13:26:07+02:00","file":"logger.go","line":24,"message":"logging this to json"}
	log.Print("info: " + str)
}

func Warning(str string) {
	log.Print("warning: " + str)
}

func Error(err error) {
	log.Print("error: " + err.Error())
}

func ErrorMsg(str string, err error) {
	log.Print("error: " + err.Error() + " errorMsg: " + str)
}
