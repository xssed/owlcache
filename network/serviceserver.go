package network

import (
	"fmt"

	"github.com/xssed/owlcache/cache"
)

//创建一个全局的缓存DB
var BaseCacheDB *cache.BaseCache

//创建一个全局的身份认证缓存
var BaseAuth *cache.BaseCache

func BaseCacheDBInit() {
	//执行步骤信息
	fmt.Println("owlcache  DB running...")
	BaseCacheDB = cache.NewCache("owlcache") //创建DB

	//加载之前缓存本地的数据
	BaseCacheDB.LoadFromFile("./owlcache.db")

	//身份认证缓存,所有身份认证都在这里有效期30分钟
	//存储内容: key tonken  value "uuid"
	BaseAuth = cache.NewCache("Auth")

	fmt.Println("owlcache  TCPServer running...")
	go stratTCP()
	fmt.Println("owlcache  HTTPServer running...")
	go stratHTTP()
}
