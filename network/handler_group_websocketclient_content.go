package network

import (
	//"sync"
	"time"

	//"github.com/xssed/owlcache/cache"
	owlconfig "github.com/xssed/owlcache/config"
	owltools "github.com/xssed/owlcache/tools"
)

//创建一个客户端传输模型
type WSCContent struct {
	Uuid             string
	Key              string
	Handshake_string string
	Content          OwlResponse
}

//函数:创建WSCContent实体
func NewWSCContent(key, address string) WSCContent {

	uuid := owltools.GetUUIDString()                          //生成UUID
	temp_handshake := owltools.JoinString(uuid, "_", address) //uuid_address   前面加key+"@"就是发送给客户端时的自定义字符串

	return WSCContent{Uuid: uuid, Key: key, Handshake_string: temp_handshake, Content: OwlResponse{}}

}

//将服务端发送给客户端的数据放进临时存储数据库BaseWSCGroupCache
func OwlResponseToWSCGroupCache(text string, owlresponse OwlResponse) {

	//缓存生命周期
	exptime, _ := time.ParseDuration(owltools.JoinString(owlconfig.OwlConfigModel.HttpClientRequestLocalCacheLifeTime, "ms"))
	//Key此时传进来的格式例:hello@f8d8ab78-6b13-4595-9b16-67dc3ff612ab_ws://127.0.0.1:7721/ws
	BaseWSCGroupCache.Set(owlresponse.Key, []byte(text), exptime)

}

//BaseWSCGroupCache的数据阅后即焚
func WSCGroupCacheBurnAfterReading(key string) ([]byte, bool) {

	//查找是否存在,如果存在删除这个缓存中的临时数据，并把数据返回获取者
	if v, found := BaseWSCGroupCache.Get(key); found {
		BaseWSCGroupCache.Delete(key)
		return v, true
	}
	return nil, false

}
