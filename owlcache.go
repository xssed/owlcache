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
	//使用多核cpu
	runtime.GOMAXPROCS(runtime.NumCPU())
	//欢迎信息
	owlsystem.DosSayHello()
	//初始化配置
	owlconfig.ConfigInit()
	//初始化日志记录
	owllog.LogInit()
	//初始化数据库服务,核心组件
	owlnetwork.BaseCacheDBInit()
	//定时任务服务
	owljob.JobInit()
	//守护包。用于保证程序的稳健、安全运行
	owlaegis.AegisInit()
	//捕获程序正常退出操作 ctrl+c
	owlsystem.OnExit()
}
