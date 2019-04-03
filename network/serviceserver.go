package network

import (
	"fmt"

	"github.com/xssed/owlcache/cache"
	owlconfig "github.com/xssed/owlcache/config"
	"github.com/xssed/owlcache/group"
	"github.com/xssed/owlcache/network/httpclient"
)

//创建一个全局的缓存DB
var BaseCacheDB *cache.BaseCache

//创建一个全局的身份认证缓存
var BaseAuth *cache.BaseCache

//创建一个全局的服务器集群信息存储列表
var ServerGroupList *group.Servergroup

//创建一个全局的HttpClient客户端
var HttpClient *httpclient.OwlClient

func BaseCacheDBInit() {

	//执行步骤信息
	fmt.Println("owlcache  database running...")

	//创建DB
	BaseCacheDB = cache.NewCache("owlcache")

	//加载之前缓存本地的DB数据
	BaseCacheDB.LoadFromFile(owlconfig.OwlConfigModel.DBfile + "owlcache.db")

	//身份认证数据,所有客户端身份认证都在这里有效期60分钟
	//存储内容: key:tonken  value:"uuid"
	BaseAuth = cache.NewCache("Auth")

	//加载之前缓存本地的DB数据
	BaseAuth.LoadFromFile(owlconfig.OwlConfigModel.DBfile + "auth.db")

	//初始化服务器集群信息存储列表
	ServerGroupList = group.NewServergroup()

	//加载之前缓存本地的服务器集群信息
	ServerGroupList.LoadFromFile(owlconfig.OwlConfigModel.DBfile, "server_group_config.json")

	//初始化HttpClient客户端
	HttpClient = httpclient.NewOwlClient()

	//启动tcp服务
	fmt.Println("owlcache  tcp server running...")
	go stratTCP()

	//启动http服务
	fmt.Println("owlcache  http server running...")
	go stratHTTP()

}
