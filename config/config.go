package config

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

//创建一个全局配置变量
var OwlConfigModel *OwlConfig

//配置文件模型
type OwlConfig struct {
	Configfile                               string //配置文件路径
	Logfile                                  string //日志文件路径
	DBfile                                   string //数据库文件路径
	Pass                                     string //owlcache密钥
	Host                                     string //主机地址
	ResponseHost                             string //程序响应IP,在TCP、HTTP返回的响应结果Json字符串中使用
	Tcpport                                  string //Tcp监听端口
	Httpport                                 string //Http监听端口
	HttpClientRequestTimeout                 string //集群互相通信时的请求超时时间
	GroupWorkMode                            string //集群方式:owlcache、gossip
	Gossipport                               string //启用Gossip服务该项才会生效。Gossip监听端口，默认值为0(系统自动监听一个端口并在启动信息输出该端口)。
	Task_DataBackup                          string //自动备份DB数据的存储时间
	Task_DataAuthBackup                      string //自动备份用户认证数据的存储时间
	Task_ClearExpireData                     string //自动清理数据库中过期数据的时间
	Task_ServerListBackup                    string //自动备份服务器集群信息数据的时间
	Get_data_from_memcache                   string //是否开启查询不存在的数据时,从memcache查询并存入本地数据库
	Memcache_list                            string //需要查询的memcache列表,不同节点之间需要使用"|"符号间隔。
	Get_memcache_data_set_expire_time        string //为从memcache存入本地数据库的Key设置一个过期时间。默认为0，永久不过期。单位是秒。
	Get_data_from_redis                      string //是否开启查询不存在的数据时,从redis查询并存入本地数据库。0表示不开启。
	Redis_Addr                               string //需要查询的Redis地址。
	Redis_Password                           string //需要查询的Redis密码。不能是空值。
	Redis_DB                                 string //需要查询Redis的数据库。默认为0。
	Get_redis_data_set_expire_time           string //为从Redis存入本地数据库的Key设置一个过期时间。默认为0，永久不过期。单位是秒。
	Open_Https                               string //是否开启HTTPS。值为0(关闭)、1(开启)。默认关闭。
	Https_CertFile                           string //Cert文件路径。
	Https_KeyFile                            string //Key文件路径。
	HttpsClient_InsecureSkipVerify           string //当开启HTTPS模式后，owlcache之间互相通讯时是否校验证书。值为0(关闭)、1(开启)。默认关闭。开启时不会校验证书合法性。
	CloseTcp                                 string //是否关闭Tcp服务(因为TCP模式下无密码认证)  值为"1"(开启)和"0"(关闭)。默认为1开启服务。
	Cors                                     string //是否开启跨域
	Access_Control_Allow_Origin              string //设置指定的域
	HttpClient_Request_Timeout_Sleeptime     string //http客户端请求设置。请求的睡眠时间。单位是整数。
	HttpClient_Request_Max_Error_Number      string //http客户端Max_Error_Number超过限定值时，http客户端请求将“暂停”Sleeptime值，来优化程序响应速度。
	MemcacheClient_Request_Timeout_Sleeptime string //MemcacheClient客户端请求设置。请求的睡眠时间。单位是整数。
	MemcacheClient_Request_Max_Error_Number  string //MemcacheClient客户端Max_Error_Number超过限定值时，MemcacheClient请求将“暂停”Sleeptime值，来优化程序响应速度。
	RedisClient_Request_Timeout_Sleeptime    string //RedisClient客户端请求设置。请求的睡眠时间。单位是整数。
	RedisClient_Request_Max_Error_Number     string //RedisClient客户端Max_Error_Number超过限定值时，RedisClient请求将“暂停”Sleeptime值，来优化程序响应速度。
}

//创建一个默认配置文件的实体
func NewDefaultOwlConfig() *OwlConfig {
	return &OwlConfig{
		Configfile:                               "owlcache.conf",
		Logfile:                                  "",
		DBfile:                                   "",
		Pass:                                     "",
		Host:                                     "0.0.0.0",
		ResponseHost:                             "127.0.0.1",
		Tcpport:                                  "7720",
		Httpport:                                 "7721",
		HttpClientRequestTimeout:                 "2",
		GroupWorkMode:                            "owlcache",
		Gossipport:                               "0",
		Task_DataBackup:                          "1",
		Task_DataAuthBackup:                      "1",
		Task_ClearExpireData:                     "1",
		Task_ServerListBackup:                    "1",
		Get_data_from_memcache:                   "0",
		Memcache_list:                            "",
		Get_memcache_data_set_expire_time:        "0",
		Get_data_from_redis:                      "0",
		Redis_Addr:                               "",
		Redis_Password:                           "",
		Redis_DB:                                 "0",
		Get_redis_data_set_expire_time:           "0",
		Open_Https:                               "0",
		Https_CertFile:                           "/www/server.crt",
		Https_KeyFile:                            "/www/server.key",
		HttpsClient_InsecureSkipVerify:           "0",
		CloseTcp:                                 "1",
		Cors:                                     "0",
		Access_Control_Allow_Origin:              "*",
		HttpClient_Request_Timeout_Sleeptime:     "2",
		HttpClient_Request_Max_Error_Number:      "2",
		MemcacheClient_Request_Timeout_Sleeptime: "2",
		MemcacheClient_Request_Max_Error_Number:  "2",
		RedisClient_Request_Timeout_Sleeptime:    "2",
		RedisClient_Request_Max_Error_Number:     "2",
	}
}

//缓存系统初始化加载配置
func ConfigInit() {

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
		var type_value string
		systype := runtime.GOOS
		if systype == "windows" {
			// windows系统
			type_value = string(rs[strings.Index(line, " ")+1 : len(rs)-2]) //-1
		} else {
			//systype == "linux" or other
			// LINUX系统
			type_value = string(rs[strings.Index(line, " ")+1 : len(rs)-1]) //-1
		}
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
	if len(config["Gossipport"]) >= 1 {
		param.Gossipport = config["Gossipport"]
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
	if len(config["Get_data_from_memcache"]) >= 1 {
		param.Get_data_from_memcache = config["Get_data_from_memcache"]
	}
	if len(config["Memcache_list"]) >= 1 {
		param.Memcache_list = config["Memcache_list"]
	}
	if len(config["Get_memcache_data_set_expire_time"]) >= 1 {
		param.Get_memcache_data_set_expire_time = config["Get_memcache_data_set_expire_time"]
	}
	if len(config["Get_data_from_redis"]) >= 1 {
		param.Get_data_from_redis = config["Get_data_from_redis"]
	}
	if len(config["Redis_Addr"]) >= 1 {
		param.Redis_Addr = config["Redis_Addr"]
	}
	if len(config["Redis_Password"]) >= 1 {
		param.Redis_Password = config["Redis_Password"]
	}
	if len(config["Redis_DB"]) >= 1 {
		param.Redis_DB = config["Redis_DB"]
	}
	if len(config["Get_redis_data_set_expire_time"]) >= 1 {
		param.Get_redis_data_set_expire_time = config["Get_redis_data_set_expire_time"]
	}
	if len(config["Open_Https"]) >= 1 {
		param.Open_Https = config["Open_Https"]
	}
	if len(config["Https_CertFile"]) >= 1 {
		param.Https_CertFile = config["Https_CertFile"]
	}
	if len(config["Https_KeyFile"]) >= 1 {
		param.Https_KeyFile = config["Https_KeyFile"]
	}
	if len(config["HttpsClient_InsecureSkipVerify"]) >= 1 {
		param.HttpsClient_InsecureSkipVerify = config["HttpsClient_InsecureSkipVerify"]
	}
	if len(config["CloseTcp"]) >= 1 {
		param.CloseTcp = config["CloseTcp"]
	}
	if len(config["Cors"]) >= 1 {
		param.Cors = config["Cors"]
	}
	if len(config["Access_Control_Allow_Origin"]) >= 1 {
		param.Access_Control_Allow_Origin = config["Access_Control_Allow_Origin"]
	}
	if len(config["HttpClient_Request_Timeout_Sleeptime"]) >= 1 {
		param.HttpClient_Request_Timeout_Sleeptime = config["HttpClient_Request_Timeout_Sleeptime"]
	}
	if len(config["HttpClient_Request_Max_Error_Number"]) >= 1 {
		param.HttpClient_Request_Max_Error_Number = config["HttpClient_Request_Max_Error_Number"]
	}
	if len(config["MemcacheClient_Request_Timeout_Sleeptime"]) >= 1 {
		param.MemcacheClient_Request_Timeout_Sleeptime = config["MemcacheClient_Request_Timeout_Sleeptime"]
	}
	if len(config["MemcacheClient_Request_Max_Error_Number"]) >= 1 {
		param.MemcacheClient_Request_Max_Error_Number = config["MemcacheClient_Request_Max_Error_Number"]
	}
	if len(config["RedisClient_Request_Timeout_Sleeptime"]) >= 1 {
		param.RedisClient_Request_Timeout_Sleeptime = config["RedisClient_Request_Timeout_Sleeptime"]
	}
	if len(config["RedisClient_Request_Max_Error_Number"]) >= 1 {
		param.RedisClient_Request_Max_Error_Number = config["RedisClient_Request_Max_Error_Number"]
	}

}
