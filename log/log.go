package log

import (
	"fmt"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
)

//创建一个全局运行日志
var OwlLogRun *OwlLog

//创建一个全局TCP日志
var OwlLogTcp *OwlLog

//创建一个全局HTTP-Single-Node日志
var OwlLogHttpSN *OwlLog

//创建一个全局HTTP-Group日志
var OwlLogHttpG *OwlLog

//创建一个全局UrlCache日志
var OwlLogUC *OwlLog

//创建一个全局Gossip日志
var OwlLogGossip *OwlLog

//创建一个全局Task日志
var OwlLogTask *OwlLog

//创建一个全局系统资源日志
var OwlLogSystemResource *OwlLog

func LogInit() {
	//执行步骤信息
	fmt.Println("owlcache  logger running...")
	//日志目录
	logFilePath := owlconfig.OwlConfigModel.Logfile
	//注册全局运行日志
	OwlLogRunRegister(logFilePath)
	//注册全局TCP日志
	OwlLogTcpRegister(logFilePath)
	//注册全局HTTP单节点日志
	OwlLogHttpSingleNodeRegister(logFilePath)
	//注册全局HTTP集群日志
	OwlLogHttpGroupRegister(logFilePath)
	//判断是否开启Url Cache
	if owlconfig.OwlConfigModel.Open_Urlcache == "1" {
		//注册全局UrlCache日志
		OwlLogUCRegister(logFilePath)
	}
	//判断是否开启Gossip数据同步
	if owlconfig.OwlConfigModel.GroupDataSync == "1" {
		//注册一个全局Gossip日志
		OwlLogGossipRegister(logFilePath)
	}
	//注册一个全局Task日志
	OwlLogTaskRegister(logFilePath)
	//注册一个全局系统资源日志
	OwlLogSystemResourceRegister(logFilePath)
}

//注册全局运行日志
func OwlLogRunRegister(logFilePath string) {
	//日志文件
	formatLogFileName := "owl_run_" + time.Now().Format("150405") + ".log"
	//创建资源
	OwlLogRun = NewOwlLog(logFilePath, formatLogFileName)
}

//注册全局TCP日志
func OwlLogTcpRegister(logFilePath string) {
	//日志文件
	formatLogFileName := "owl_tcp_" + time.Now().Format("150405") + ".log"
	//创建资源
	OwlLogTcp = NewOwlLog(logFilePath, formatLogFileName)
}

//注册全局HTTP单节点日志
func OwlLogHttpSingleNodeRegister(logFilePath string) {
	//日志文件
	formatLogFileName := "owl_http_single_node_" + time.Now().Format("150405") + ".log"
	//创建资源
	OwlLogHttpSN = NewOwlLog(logFilePath, formatLogFileName)
}

//注册全局HTTP集群日志
func OwlLogHttpGroupRegister(logFilePath string) {
	//日志文件
	formatLogFileName := "owl_http_group_" + time.Now().Format("150405") + ".log"
	//创建资源
	OwlLogHttpG = NewOwlLog(logFilePath, formatLogFileName)
}

//注册全局UrlCache日志
func OwlLogUCRegister(logFilePath string) {
	//日志文件
	formatLogFileName := "owl_uc_" + time.Now().Format("150405") + ".log"
	//创建资源
	OwlLogUC = NewOwlLog(logFilePath, formatLogFileName)
}

//注册一个全局Gossip日志
func OwlLogGossipRegister(logFilePath string) {
	//日志文件
	formatLogFileName := "owl_gossip_" + time.Now().Format("150405") + ".log"
	//创建资源
	OwlLogGossip = NewOwlLog(logFilePath, formatLogFileName)
}

//注册一个全局Task日志
func OwlLogTaskRegister(logFilePath string) {
	//日志文件
	formatLogFileName := "owl_task_" + time.Now().Format("150405") + ".log"
	//创建资源
	OwlLogTask = NewOwlLog(logFilePath, formatLogFileName)
}

//注册一个全局系统资源日志
func OwlLogSystemResourceRegister(logFilePath string) {
	//日志文件
	formatLogFileName := "owl_system_resource_" + time.Now().Format("150405") + ".log"
	//创建资源
	OwlLogSystemResource = NewOwlLog(logFilePath, formatLogFileName)
}
