package config

import (
	"flag"
	"fmt"
)

//Owlcache启动时对接收参数初始化
func CmdParamInit(param *OwlConfig) *OwlConfig {

	//创建一个默认配置模型
	//	var OwlConfigBase *OwlConfig
	//	OwlConfigBase = NewDefaultOwlConfig()

	//将Owlcache启动时接收到的参数进行解析绑定
	//cmdConfig := CmdParamExe(OwlConfigBase)
	cmdConfig := CmdParamExe(param)

	//执行步骤信息
	fmt.Println("owlcache  paraminit running... ")

	return cmdConfig

}

//将Owlcache启动时接收到的参数进行解析绑定
func CmdParamExe(param *OwlConfig) *OwlConfig {
	host := flag.String("host", param.Host, "binding local host ip adress.")
	configPath := flag.String("config", param.Configfile, "owlcache config file path.[demo:/var/home/owl.conf]")
	logPath := flag.String("log", param.Logfile, "owlcache log file path.[demo:/var/log/]") //owlcache auto generation
	pass := flag.String("pass", param.Pass, "owlcache Http connection password.")
	flag.Parse()

	if len(*host) > 0 {
		param.Host = *host
	}
	if len(*configPath) > 0 {
		param.Configfile = *configPath
	}
	if len(*logPath) > 0 {
		param.Logfile = *logPath
	}
	if len(*pass) > 0 {
		param.Pass = *pass
	}

	return param

}
