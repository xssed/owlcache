package network

import (
	"net/http"
	"strconv"
	"time"

	"github.com/xssed/owlcache/cache"
	"github.com/xssed/owlcache/network/gossip"
	"github.com/xssed/owlcache/network/memcacheclient"
	"github.com/xssed/owlcache/network/redisclient"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
	owltoken "github.com/xssed/owlcache/tools/token"
)

//执行K/V数据查询,本地内存数据库->Memcache(如果开启)->Redis(如果开启）
func (owlhandler *OwlHandler) baseget() {

	if v, found := BaseCacheDB.GetKvStore(owlhandler.owlrequest.Key); found {
		owlhandler.Transmit(SUCCESS)
		owlhandler.owlresponse.Data = v.(*cache.KvStore).Value
		owlhandler.owlresponse.KeyCreateTime = v.(*cache.KvStore).CreateTime
		return
	} else {
		//NOT_FOUND状态下是否从memcache中查询数据
		if owlconfig.OwlConfigModel.Get_data_from_memcache == "1" {
			//请求优化部分
			temp_mcrts_exptime := owltools.DoubleNumberStringSubToString(owlconfig.OwlConfigModel.MemcacheClient_Request_Timeout_Sleeptime, "1") //字符串相减
			mcrts_exptime, _ := time.ParseDuration(owltools.JoinString(temp_mcrts_exptime, "s"))                                                 //拼接字符串转化为时间，请求失败的睡眠时间
			mrmen_maxnum, _ := strconv.Atoi(owlconfig.OwlConfigModel.MemcacheClient_Request_Max_Error_Number)                                    //最大错误请求数，超过该数就进入睡眠
			k := MemcacheClientRequestErrorCounter.Exe(owlhandler.owlrequest.Key, int64(mrmen_maxnum-1), mcrts_exptime)
			if k > 0 {
				//请求数据
				//owllog.OwlLogRun.Info("memcacheclient:get key " + " key:" + owlhandler.owlrequest.Key)
				result, err := memcacheclient.Get(owlhandler.owlrequest.Key)
				if err == nil {

					//执行成功-1
					MemcacheClientRequestErrorCounter.Dec(owlhandler.owlrequest.Key)

					//找到数据了，存入owlcache中
					exptime, _ := time.ParseDuration(owlconfig.OwlConfigModel.Get_memcache_data_set_expire_time + "s")
					ok := BaseCacheDB.Set(string(result.Key), result.Value, exptime)
					//设置数据时出错
					if !ok {
						owllog.OwlLogRun.Info("Get_data_from_memcache:Store data to owlcache  error, " + " key:" + owlhandler.owlrequest.Key)
					} else {
						owlhandler.Transmit(SUCCESS)
						owlhandler.owlresponse.Data = result.Value
						owlhandler.owlresponse.KeyCreateTime = time.Now()
						return
					}
				}
			}
		}
		//NOT_FOUND状态下是否从redis中查询数据
		if owlconfig.OwlConfigModel.Get_data_from_redis == "1" {

			//请求优化部分
			temp_rcrts_exptime := owltools.DoubleNumberStringSubToString(owlconfig.OwlConfigModel.RedisClient_Request_Timeout_Sleeptime, "1") //字符串相减
			rcrts_exptime, _ := time.ParseDuration(owltools.JoinString(temp_rcrts_exptime, "s"))                                              //拼接字符串转化为时间，请求失败的睡眠时间
			rcrmen_maxnum, _ := strconv.Atoi(owlconfig.OwlConfigModel.RedisClient_Request_Max_Error_Number)                                   //最大错误请求数，超过该数就进入睡眠
			k := RedisClientRequestErrorCounter.Exe(owlhandler.owlrequest.Key, int64(rcrmen_maxnum-1), rcrts_exptime)
			if k > 0 {
				//请求数据
				//owllog.OwlLogRun.Info("redisclient:get key " + " key:" + owlhandler.owlrequest.Key)
				rcres, err := redisclient.Get(owlhandler.owlrequest.Key)
				if err == nil {

					//执行成功-1
					RedisClientRequestErrorCounter.Dec(owlhandler.owlrequest.Key)

					//找到数据了，存入owlcache中
					rcexptime, _ := time.ParseDuration(owlconfig.OwlConfigModel.Get_redis_data_set_expire_time + "s")
					ok := BaseCacheDB.Set(owlhandler.owlrequest.Key, []byte(rcres), rcexptime)
					//设置数据时出错
					if !ok {
						owllog.OwlLogRun.Info("Get_data_from_redis:Store data to owlcache error" + " key:" + owlhandler.owlrequest.Key)
					} else {
						owlhandler.Transmit(SUCCESS)
						owlhandler.owlresponse.Data = []byte(rcres)
						owlhandler.owlresponse.KeyCreateTime = time.Now()
						return
					}
				}
			}
		}

		owlhandler.Transmit(NOT_FOUND)
		return
	}

}

func (owlhandler *OwlHandler) Get() {

	//执行K/V数据查询,本地内存数据库->Memcache(如果开启)->Redis(如果开启）
	owlhandler.baseget()

}

func (owlhandler *OwlHandler) Exists() {

	ok := BaseCacheDB.Exists(owlhandler.owlrequest.Key)
	if ok {
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(NOT_FOUND)
	}
}

func (owlhandler *OwlHandler) Set() {

	ok := BaseCacheDB.Set(owlhandler.owlrequest.Key, owlhandler.owlrequest.Value, owlhandler.owlrequest.Expires)
	if ok {
		owlhandler.Transmit(SUCCESS)
		owlhandler.owlresponse.Data = []byte("")
		owlhandler.owlresponse.KeyCreateTime = time.Now()
	} else {
		owlhandler.Transmit(ERROR)
	}

	//判断一致性数据同步-设置
	if owlconfig.OwlConfigModel.GroupDataSync == "1" {
		//发送数据到集群
		prefix := "http://"
		if owlconfig.OwlConfigModel.Open_Https == "1" {
			prefix = "https://"
		}
		key_resource := owltools.JoinString(prefix, owlconfig.OwlConfigModel.ResponseHost, ":", owlconfig.OwlConfigModel.Httpport)
		gossip.Set(owlhandler.owlrequest.Key, key_resource, owlhandler.owlrequest.Expires)
	}

}

func (owlhandler *OwlHandler) Expire() {

	//判断要设置过期的Key是否存在
	exist := BaseCacheDB.Exists(owlhandler.owlrequest.Key)
	if exist {
		ok := BaseCacheDB.Expire(owlhandler.owlrequest.Key, owlhandler.owlrequest.Expires)
		if ok {
			owlhandler.Transmit(SUCCESS)
		} else {
			owlhandler.Transmit(ERROR)
		}

		//判断一致性数据同步-设置Key过期
		if owlconfig.OwlConfigModel.GroupDataSync == "1" {
			//发送数据到集群
			gossip.Expire(owlhandler.owlrequest.Key, owlhandler.owlrequest.Expires)
		}

	} else {
		owlhandler.Transmit(NOT_FOUND)
	}

}

func (owlhandler *OwlHandler) Delete() {

	//判断要删除的Key是否存在
	exist := BaseCacheDB.Exists(owlhandler.owlrequest.Key)
	if exist {
		//数据存在删除它
		ok := BaseCacheDB.Delete(owlhandler.owlrequest.Key)
		if !ok {
			owlhandler.Transmit(SUCCESS)
		} else {
			owlhandler.Transmit(ERROR)
		}

		//判断一致性数据同步-删除Key
		if owlconfig.OwlConfigModel.GroupDataSync == "1" {
			//发送数据到集群
			gossip.Delete(owlhandler.owlrequest.Key)
		}

	} else {
		owlhandler.Transmit(NOT_FOUND)
	}

}

//PASS命令验证密码
func (owlhandler *OwlHandler) Pass(r *http.Request) {

	if owlconfig.OwlConfigModel.Pass == owlhandler.owlrequest.Pass {
		//token_id=md5(ip+uuid)
		uuid := owltools.GetUUIDString()
		ip := owltools.RemoteAddr2IPAddr(r.RemoteAddr)
		token_id := owltools.GetMd5String(ip + uuid)
		expiration, _ := time.ParseDuration(owlconfig.OwlConfigModel.Tonken_expire_time + "s")
		BaseAuth.Set(token_id, []byte(ip), expiration)
		//创建JWT标准化token
		token, generate_err := owltoken.GenerateToken(token_id)
		if generate_err != nil {
			owlhandler.Transmit(ERROR)
			return
		}
		owlhandler.owlresponse.Data = []byte(token)
		owlhandler.owlresponse.KeyCreateTime = time.Now()
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(ERROR)
	}

}

//验证权限
func (owlhandler *OwlHandler) CheckAuth(r *http.Request) bool {

	//先判断字符串是不是Base64编码如果是就解析
	if owltools.IsBase64(owlhandler.owlrequest.Token) {
		decode_token := string(owltools.Base64Decode(owlhandler.owlrequest.Token, "url"))
		if string(decode_token) == "" {
			return false
		}
		owlhandler.owlrequest.Token = decode_token //重新赋值
	}

	token, parse_err := owltoken.ParseToken(owlhandler.owlrequest.Token)
	//解析出错或者token过期
	if parse_err != nil || token == nil {
		return false
	}
	token_id := token.TokenId
	ip := owltools.RemoteAddr2IPAddr(r.RemoteAddr)
	v, found := BaseAuth.Get(token_id)
	if found == true {
		if string(v) == ip {
			return true
		}
		return false
	}
	return false

}

//Ping 命令
func (owlhandler *OwlHandler) Ping() {

	owlhandler.Transmit(SUCCESS)
	//判断有无输入字符串，有则返回该字符串，无则返回"PONG"
	if owlhandler.owlrequest.Length > 0 {
		owlhandler.owlresponse.Data = []byte(owlhandler.owlrequest.Value)
	} else {
		owlhandler.owlresponse.Data = []byte("PONG")
	}

}
