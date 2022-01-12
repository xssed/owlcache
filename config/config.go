package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xssed/owlcache/tools"
)

//创建一个全局配置变量
var OwlConfigModel *OwlConfig

//配置文件模型
type OwlConfig struct {
	Configfile                               string //配置文件路径
	Logfile                                  string //日志文件路径
	DBfile                                   string //数据库文件路径
	Pass                                     string //owlcache密钥
	Tonken_expire_time                       string //为Pass命令产生的Tonken值设置一个过期时间
	Host                                     string //主机地址
	ResponseHost                             string //程序响应IP,在TCP、HTTP返回的响应结果Json字符串中使用
	Tcpport                                  string //Tcp监听端口
	Httpport                                 string //Http监听端口
	HttpClientRequestTimeout                 string //集群互相通信时的请求超时时间
	HttpClientRequestLocalCacheLifeTime      string //Http客户端请求优化设置。请求得到得数据将会短时间缓存在本地,避免同一个Key在高并发状态下反复对其他节点疯狂查询。适当的延迟查询有助于性能和效率得优化。单位毫秒，默认这个缓存生命周期为5000毫秒。
	GroupDataSync                            string //是否开启集群数据同步。0表示不开启。1表示开启。默认不开启。
	Gossipport                               string //启用Gossip服务该项才会生效。Gossip监听端口，默认值为0(系统自动监听一个端口并在启动信息输出该端口)。
	GossipDataSyncAuthKey                    string //启用Gossip服务该项才会生效。集群中通过Gossip协议交换数据的令牌。整个集群需要统一的令牌。
	GossipHttpClientRequestTimeout           string //集群同步数据互相通信时的请求超时时间
	Task_DataBackup                          string //自动备份DB数据的存储时间
	Task_DataAuthBackup                      string //自动备份用户认证数据的存储时间
	Task_ClearExpireData                     string //自动清理数据库中过期数据的时间
	Task_ServerListBackup                    string //自动备份服务器集群信息数据的时间
	Task_ServerGossipListBackup              string //自动备份Gossip集群信息数据的时间
	Task_MemoryInfoToLog                     string //自动输出Owl的内存信息的时间
	Task_OutputSuccessLog                    string //开启Task执行成功日志输出
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
	Open_Websocket_Server                    string //是否开启Websocket Server。值为0(关闭)、1(开启)。默认关闭。
	GroupData_Mode                           string //owlcache集群方式，分为“Http”(短链接)和“Websocket”(长链接)，两个选项，默认是“Http”短链接。
	Websocket_Client_WriteWait               string //客户端写超时。单位秒。默认10秒。
	Websocket_Client_MaxMessageSize          string //客户端支持接受的消息最大长度，单位字节。默认7000000字节。
	Websocket_Client_MinRecTime              string //客户端断开与服务端后最小重连时间间隔。单位秒。默认2秒。
	Websocket_Client_MaxRecTime              string //客户端断开与服务端后最大重连时间间隔。单位秒。默认60秒。
	Websocket_Client_MessageBufferSize       string //客户端消息发送缓冲池大小，单位字节。默认1024字节。
	CloseTcp                                 string //是否关闭Tcp服务(因为TCP模式下无密码认证)  值为"1"(开启)和"0"(关闭)。默认为1开启服务。
	Cors                                     string //是否开启跨域
	Access_Control_Allow_Origin              string //设置指定的域
	MemcacheClient_Request_Timeout_Sleeptime string //MemcacheClient客户端请求设置。请求的睡眠时间。单位是整数。
	MemcacheClient_Request_Max_Error_Number  string //MemcacheClient客户端Max_Error_Number超过限定值时，MemcacheClient请求将“暂停”Sleeptime值，来优化程序响应速度。
	RedisClient_Request_Timeout_Sleeptime    string //RedisClient客户端请求设置。请求的睡眠时间。单位是整数。
	RedisClient_Request_Max_Error_Number     string //RedisClient客户端Max_Error_Number超过限定值时，RedisClient请求将“暂停”Sleeptime值，来优化程序响应速度。
	Open_Urlcache                            string //URL代理访问后将得到的HTTP响应数据缓存到Owlcache中，值为"1"(开启服务)和"0"(关闭服务)。默认为0关闭服务。
	Urlcache_Filename                        string //开启URL缓存后的需要加载的配置文件的名称，格式为XML，它存在于DBfile配置项目录中，默认名称为sites.xml。
	Urlcache_Request_Easy                    string //开启Urlcache的快捷访问。不影响UrlCache的默认访问方式。值为"1"(开启)和"0"(关闭)。默认为0关闭。
}

//创建一个默认配置文件的实体
func NewDefaultOwlConfig() *OwlConfig {
	return &OwlConfig{
		Configfile:                               "owlcache.conf",
		Logfile:                                  "",
		DBfile:                                   "",
		Pass:                                     "",
		Tonken_expire_time:                       "0",
		Host:                                     "0.0.0.0",
		ResponseHost:                             "127.0.0.1",
		Tcpport:                                  "7720",
		Httpport:                                 "7721",
		HttpClientRequestTimeout:                 "2700",
		HttpClientRequestLocalCacheLifeTime:      "5000",
		GroupDataSync:                            "0",
		Gossipport:                               "0",
		GossipDataSyncAuthKey:                    "",
		GossipHttpClientRequestTimeout:           "5000",
		Task_DataBackup:                          "1",
		Task_DataAuthBackup:                      "1",
		Task_ClearExpireData:                     "1",
		Task_ServerListBackup:                    "1",
		Task_ServerGossipListBackup:              "1",
		Task_MemoryInfoToLog:                     "1",
		Task_OutputSuccessLog:                    "1",
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
		Open_Websocket_Server:                    "0",
		GroupData_Mode:                           "Http",
		Websocket_Client_WriteWait:               "10",
		Websocket_Client_MaxMessageSize:          "7000000",
		Websocket_Client_MinRecTime:              "2",
		Websocket_Client_MaxRecTime:              "60",
		Websocket_Client_MessageBufferSize:       "1024",
		CloseTcp:                                 "1",
		Cors:                                     "0",
		Access_Control_Allow_Origin:              "*",
		MemcacheClient_Request_Timeout_Sleeptime: "2",
		MemcacheClient_Request_Max_Error_Number:  "2",
		RedisClient_Request_Timeout_Sleeptime:    "2",
		RedisClient_Request_Max_Error_Number:     "2",
		Open_Urlcache:                            "0",
		Urlcache_Filename:                        "sites.xml",
		Urlcache_Request_Easy:                    "0",
	}
}

//缓存系统初始化加载配置
func ConfigInit() {

	//读取配置文件获取一个最终的配置模型
	var config = make(map[string]string)

	//创建一个默认初始化配置模型
	OwlConfigModel = NewDefaultOwlConfig()

	//打开配置文件
	config_file, err := os.Open(OwlConfigModel.Configfile)
	defer config_file.Close()
	//监控错误
	if err != nil {
		fmt.Println(err)
		fmt.Print("Can not read configuration file. now exit\n")
		os.Exit(0)
	}
	//判断文件是否为UTF-8编码
	if !tools.ValidUTF8(OwlConfigModel.Configfile) {
		fmt.Print("The configuration file is not UTF-8 encode. now exit\n")
		os.Exit(0)
	}
	//将内容读入缓冲区
	buff := bufio.NewReader(config_file)
	//读取配置文件
	for {
		line, _, err := buff.ReadLine()
		if err != nil {
			break
		}
		rs := []rune(string(line))
		if string(rs[0:1]) == `#` || len(string(line)) < 3 {
			continue
		}
		if string(rs[0:1]) == `[` || len(string(line)) < 3 {
			continue
		}
		config_string_arr := strings.Split(string(line), " ")

		type_name := config_string_arr[0]
		type_value := config_string_arr[1]

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
	fmt.Println("owlcache  system configuration initialization is complete...")

	//检查是否启动URL缓存功能
	if OwlConfigModel.Open_Urlcache == "1" {
		//URL缓存初始化加载配置
		OwlUCConfigInit()
	}

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
	if len(config["Tonken_expire_time"]) > 1 {
		param.Tonken_expire_time = config["Tonken_expire_time"]
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
	if len(config["HttpClientRequestLocalCacheLifeTime"]) >= 1 {
		param.HttpClientRequestLocalCacheLifeTime = config["HttpClientRequestLocalCacheLifeTime"]
	}
	if len(config["GroupDataSync"]) >= 1 {
		param.GroupDataSync = config["GroupDataSync"]
	}
	if len(config["Gossipport"]) >= 1 {
		param.Gossipport = config["Gossipport"]
	}
	if len(config["GossipDataSyncAuthKey"]) >= 1 {
		param.GossipDataSyncAuthKey = config["GossipDataSyncAuthKey"]
	}
	if len(config["GossipHttpClientRequestTimeout"]) >= 1 {
		param.GossipHttpClientRequestTimeout = config["GossipHttpClientRequestTimeout"]
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
	if len(config["Task_ServerGossipListBackup"]) >= 1 {
		param.Task_ServerGossipListBackup = config["Task_ServerGossipListBackup"]
	}
	if len(config["Task_MemoryInfoToLog"]) >= 1 {
		param.Task_MemoryInfoToLog = config["Task_MemoryInfoToLog"]
	}
	if len(config["Task_OutputSuccessLog"]) >= 1 {
		param.Task_OutputSuccessLog = config["Task_OutputSuccessLog"]
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
	if len(config["Open_Websocket_Server"]) >= 1 {
		param.Open_Websocket_Server = config["Open_Websocket_Server"]
	}
	if len(config["GroupData_Mode"]) >= 1 {
		param.GroupData_Mode = config["GroupData_Mode"]
	}
	if len(config["Websocket_Client_WriteWait"]) >= 1 {
		param.Websocket_Client_WriteWait = config["Websocket_Client_WriteWait"]
	}
	if len(config["Websocket_Client_MaxMessageSize"]) >= 1 {
		param.Websocket_Client_MaxMessageSize = config["Websocket_Client_MaxMessageSize"]
	}
	if len(config["Websocket_Client_MinRecTime"]) >= 1 {
		param.Websocket_Client_MinRecTime = config["Websocket_Client_MinRecTime"]
	}
	if len(config["Websocket_Client_MaxRecTime"]) >= 1 {
		param.Websocket_Client_MaxRecTime = config["Websocket_Client_MaxRecTime"]
	}
	if len(config["Websocket_Client_MessageBufferSize"]) >= 1 {
		param.Websocket_Client_MessageBufferSize = config["Websocket_Client_MessageBufferSize"]
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
	if len(config["Open_Urlcache"]) >= 1 {
		param.Open_Urlcache = config["Open_Urlcache"]
	}
	if len(config["Urlcache_Filename"]) >= 1 {
		param.Urlcache_Filename = config["Urlcache_Filename"]
	}
	if len(config["Urlcache_Request_Easy"]) >= 1 {
		param.Urlcache_Request_Easy = config["Urlcache_Request_Easy"]
	}

}
