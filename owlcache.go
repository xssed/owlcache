package main

import (
	"runtime"

	owlaegis "github.com/xssed/owlcache/aegis"
	owlconfig "github.com/xssed/owlcache/config"
	owljob "github.com/xssed/owlcache/job"
	owllog "github.com/xssed/owlcache/log"
	owlnetwork "github.com/xssed/owlcache/network"
	owlsystem "github.com/xssed/owlcache/system"
)

//                _                _
//   _____      _| | ___ __ _  ___| |__   ___
//  / _ \ \ /\ / / |/ __/ _` |/ __| '_ \ / _ \
// | (_) \ V  V /| | (_| (_| | (__| | | |  __/
//  \___/ \_/\_/ |_|\___\__,_|\___|_| |_|\___|
//
//If you have any questions,Please contact us: xsser@xsser.cc
//Project Home:https://github.com/xssed/owlcache

func main() {
	//使用多核cpu(Use multi-core cpu)
	runtime.GOMAXPROCS(runtime.NumCPU())
	//欢迎信息(Welcome message)
	owlsystem.DosSayHello()
	//初始化配置(Initial configuration)
	owlconfig.ConfigInit()
	//初始化日志记录(Initialize logging)
	owllog.LogInit()
	//定时任务服务(Scheduled Task Service)
	owljob.JobInit()
	//初始化数据库服务,核心组件(Initialize database services, core components)
	owlnetwork.BaseCacheDBInit()
	//守护包。用于保证程序的稳健、安全运行。(It is used to ensure the stable and safe operation of the program.)
	owlaegis.AegisInit()
}
