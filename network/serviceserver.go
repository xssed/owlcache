package network

import (
	"fmt"

	"github.com/xssed/owlcache/cache"
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
	BaseCacheDB = cache.NewCache("owlcache") //创建DB

	//加载之前缓存本地的数据
	BaseCacheDB.LoadFromFile("./owlcache.db")

	//身份认证缓存,所有身份认证都在这里有效期30分钟
	//存储内容: key:tonken  value:"uuid"
	BaseAuth = cache.NewCache("Auth")

	//初始化服务器集群信息存储列表
	ServerGroupList = group.NewServergroup()

	//加载之前缓存本地的服务器集群信息
	ServerGroupList.LoadFromFile("./servergroup.db")

	//初始化HttpClient客户端
	HttpClient = httpclient.NewOwlClient()

	fmt.Println("owlcache  tcp server running...")
	go stratTCP()
	fmt.Println("owlcache  http server running...")
	go stratHTTP()
}
