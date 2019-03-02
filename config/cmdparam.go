package config

import (
	"flag"
	"fmt"
)

//Owlcache启动时对接收参数初始化
func CmdParamInit() *OwlConfig {

	//创建一个默认配置模型
	var OwlConfigBase *OwlConfig
	OwlConfigBase = NewDefaultOwlConfig()

	//将Owlcache启动时接收到的参数进行解析绑定
	cmdConfig := CmdParamExe(OwlConfigBase)

	//执行步骤信息
	fmt.Println("owlcache  starting... ")

	return cmdConfig

}

//将Owlcache启动时接收到的参数进行解析绑定
func CmdParamExe(param *OwlConfig) *OwlConfig {
	host := flag.String("host", param.Host, "binding local host ip adress.")
	configPath := flag.String("config", param.Configfile, "Owlcache config file path.[demo:/var/home/owl.conf]")
	logPath := flag.String("log", "", "Owlcache log file path.[demo:/var/log/]") //Owlcache auto generation
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

	return param

}
