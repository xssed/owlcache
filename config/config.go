package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//创建一个全局配置变量
var OwlConfigModel *OwlConfig

//配置文件模型
type OwlConfig struct {
	Configfile               string //配置文件路径
	Logfile                  string //日志文件路径
	DBfile                   string //数据库文件路径
	Pass                     string //owlcache密钥
	Host                     string //主机地址
	ResponseHost             string //程序响应IP,在TCP、HTTP返回的响应结果Json字符串中使用
	Tcpport                  string //Tcp监听端口
	Httpport                 string //Http监听端口
	HttpClientRequestTimeout string //集群互相通信时的请求超时时间
	GroupWorkMode            string //集群方式:owlcache、gossip
	Task_DataBackup          string //自动备份DB数据的存储时间
	Task_DataAuthBackup      string //自动备份用户认证数据的存储时间
	Task_ClearExpireData     string //自动清理数据库中过期数据的时间
	Task_ServerListBackup    string //自动备份服务器集群信息数据的时间
}

//创建一个默认配置文件的实体
func NewDefaultOwlConfig() *OwlConfig {
	return &OwlConfig{
		Configfile:               "owlcache.conf",
		Logfile:                  "",
		DBfile:                   "",
		Pass:                     "shi!jie9he?ping6",
		Host:                     "0.0.0.0",
		ResponseHost:             "127.0.0.1",
		Tcpport:                  "7720",
		Httpport:                 "7721",
		HttpClientRequestTimeout: "2",
		GroupWorkMode:            "owlcache",
		Task_DataBackup:          "1",
		Task_DataAuthBackup:      "1",
		Task_ClearExpireData:     "1",
		Task_ServerListBackup:    "1",
	}
}

//缓存系统初始化加载配置
func ConfigInit() {

	//执行步骤信息
	//fmt.Println("owlcache  loading config...")

	//读取配置文件获取一个最终的配置模型
	var config = make(map[string]string)

	//创建一个默认初始化配置模型
	OwlConfigModel = NewDefaultOwlConfig()

	config_file, err := os.Open(OwlConfigModel.Configfile) //打开配置文件
	defer config_file.Close()
	if err != nil {
		fmt.Println(err)
		fmt.Print("Can not read configuration file. now exit\n")
		os.Exit(0)
	}
	buff := bufio.NewReader(config_file) //将内容读入缓冲区
	//读取配置文件
	for {
		line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil {
			break
		}
		rs := []rune(line)
		if string(rs[0:1]) == `#` || len(line) < 3 {
			continue
		}
		if string(rs[0:1]) == `[` || len(line) < 3 {
			continue
		}
		type_name := string(rs[0:strings.Index(line, " ")])
		type_value := string(rs[strings.Index(line, " ")+1 : len(rs)-2]) //-1
		config[type_name] = type_value
	}

	//fmt.Println(OwlConfigModel) //打印出默认配置信息
	//fmt.Println(config)         //打印出*.conf中的配置信息
	//将文本配置绑定到全局配置
	ConfigBind(config, OwlConfigModel)
	//最后检查参数
	OwlConfigModel = CmdParamInit(OwlConfigModel)
	//fmt.Println(OwlConfigModel) //打印出最终赋值后的配置信息

	//执行步骤信息
	fmt.Println("owlcache  configuration initialization is complete...")

}

//将文本配置绑定到全局配置
func ConfigBind(config map[string]string, param *OwlConfig) {

	if len(config["Host"]) > 3 {
		//!!!如果在命令行中启动服务时指定了Host值，而配置文件这里没有注释掉则会以配置文件为准
		param.Host = config["Host"]
	}
	if len(config["ResponseHost"]) > 3 {
		param.ResponseHost = config["ResponseHost"]
	}
	if len(config["Tcpport"]) > 1 {
		param.Tcpport = config["Tcpport"]
	}
	if len(config["Httpport"]) > 1 {
		param.Httpport = config["Httpport"]
	}
	if len(config["Pass"]) > 1 {
		param.Pass = config["Pass"]
	}
	if len(config["Logfile"]) > 3 {
		//!!!如果在命令行中启动服务时指定了Logfile值，而配置文件这里没有注释掉则会以配置文件为准
		param.Logfile = config["Logfile"]
	}
	if len(config["DBfile"]) > 3 {
		param.DBfile = config["DBfile"]
	}
	if len(config["HttpClientRequestTimeout"]) >= 1 {
		param.HttpClientRequestTimeout = config["HttpClientRequestTimeout"]
	}
	if len(config["GroupWorkMode"]) >= 1 {
		param.GroupWorkMode = config["GroupWorkMode"]
	}
	if len(config["Task_DataBackup"]) >= 1 {
		param.Task_DataBackup = config["Task_DataBackup"]
	}
	if len(config["Task_DataAuthBackup"]) >= 1 {
		param.Task_DataAuthBackup = config["Task_DataAuthBackup"]
	}
	if len(config["Task_ClearExpireData"]) >= 1 {
		param.Task_ClearExpireData = config["Task_ClearExpireData"]
	}
	if len(config["Task_ServerListBackup"]) >= 1 {
		param.Task_ServerListBackup = config["Task_ServerListBackup"]
	}
}
