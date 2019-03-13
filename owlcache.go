package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	owlconfig "github.com/xssed/owlcache/config"
	owljob "github.com/xssed/owlcache/job"
	owllog "github.com/xssed/owlcache/log"
	owlnetwork "github.com/xssed/owlcache/network"
	owlsystem "github.com/xssed/owlcache/system"
)

func main() {
	//使用多核cpu
	runtime.GOMAXPROCS(runtime.NumCPU())
	//欢迎信息
	fmt.Println("Welcome to use Owlcache. Version:" + owlsystem.VERSION + "\nIf you have any questions,Please contact us: xsser@xsser.cc \nProject Home:https://github.com/xssed/owlcache")
	//初始化配置
	owlconfig.ConfigInit()
	//初始化日志记录
	owllog.LogInit()
	//初始化数据库服务
	owlnetwork.BaseCacheDBInit()
	//定时任务服务
	owljob.JobInit()

	//捕获程序正常退出操作 ctrl+c
	OnExit()
}

//捕获程序正常退出操作 ctrl+c
func OnExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	owllog.Info("Owlcache is stoped") //日志记录
	fmt.Println("Owlcache is stoped \nBye!")
}
