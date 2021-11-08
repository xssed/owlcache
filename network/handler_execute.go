package network

import (
	//"fmt"
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
			owlhandler.getfrommemcache()
		}
		//NOT_FOUND状态下是否从redis中查询数据
		if owlconfig.OwlConfigModel.Get_data_from_redis == "1" {
			owlhandler.getfromredis()
		}
		owlhandler.Transmit(NOT_FOUND)
		return
	}

}

func (owlhandler *OwlHandler) getfrommemcache() {
	//请求优化部分
	mcrts_exptime, _ := time.ParseDuration(owlconfig.OwlConfigModel.MemcacheClient_Request_Timeout_Sleeptime + "s") //睡眠时间
	mrmen_maxnum, _ := strconv.Atoi(owlconfig.OwlConfigModel.MemcacheClient_Request_Max_Error_Number)               //最大请求数
	if MemcacheClientRequestErrorCounter.Add(owlhandler.owlrequest.Key, int64(mrmen_maxnum-1), mcrts_exptime) > 0 {
		//请求数据
		//owllog.OwlLogRun.Println("memcacheclient:get key " + " key:" + owlhandler.owlrequest.Key)
		result, err := memcacheclient.Get(owlhandler.owlrequest.Key)
		if err == nil {
			//找到数据了
			exptime, _ := time.ParseDuration(owlconfig.OwlConfigModel.Get_memcache_data_set_expire_time + "s")
			ok := BaseCacheDB.Set(string(result.Key), []byte(result.Value), exptime)
			//设置数据时出错
			if !ok {
				owllog.OwlLogRun.Println("Get_data_from_memcache:set error " + " key:" + owlhandler.owlrequest.Key)
			} else {
				owlhandler.Transmit(SUCCESS)
				owlhandler.owlresponse.Data = result.Value
				owlhandler.owlresponse.KeyCreateTime = time.Now()
				return
			}
		} else {
			//memcache中也没有找到数据
			owllog.OwlLogRun.Println("Get_data_from_memcache:get error " + " key:" + owlhandler.owlrequest.Key)
		}
	}
}

func (owlhandler *OwlHandler) getfromredis() {
	//请求优化部分
	rcrts_exptime, _ := time.ParseDuration(owlconfig.OwlConfigModel.RedisClient_Request_Timeout_Sleeptime + "s") //睡眠时间
	rcrmen_maxnum, _ := strconv.Atoi(owlconfig.OwlConfigModel.RedisClient_Request_Max_Error_Number)              //最大请求数
	if RedisClientRequestErrorCounter.Add(owlhandler.owlrequest.Key, int64(rcrmen_maxnum-1), rcrts_exptime) > 0 {
		//请求数据
		rcres, err := redisclient.Get(owlhandler.owlrequest.Key)
		if err == nil {
			//找到数据了
			rcexptime, _ := time.ParseDuration(owlconfig.OwlConfigModel.Get_redis_data_set_expire_time + "s")
			ok := BaseCacheDB.Set(owlhandler.owlrequest.Key, []byte(rcres), rcexptime)
			//设置数据时出错
			if !ok {
				owllog.OwlLogRun.Println("Get_data_from_redis:set error " + " key:" + owlhandler.owlrequest.Key)
			} else {
				owlhandler.Transmit(SUCCESS)
				owlhandler.owlresponse.Data = []byte(rcres)
				owlhandler.owlresponse.KeyCreateTime = time.Now()
				return
			}
		} else {
			//redis中也没有找到数据
			owllog.OwlLogRun.Println("Get_data_from_redis:get error " + " key:" + owlhandler.owlrequest.Key)
		}
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
	} else {
		owlhandler.Transmit(ERROR)
	}

	//判断一致性数据同步-设置
	if owlconfig.OwlConfigModel.GroupDataSync == "1" {
		//fmt.Println("gossip.Set()")
		//发送数据到集群
		gossip.Set(owlhandler.owlrequest.Key, string(owlhandler.owlrequest.Value), owlhandler.owlrequest.Expires)
	}

}

func (owlhandler *OwlHandler) Expire() {

	ok := BaseCacheDB.Expire(owlhandler.owlrequest.Key, owlhandler.owlrequest.Expires)
	if ok {
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(ERROR)
	}

	//判断一致性数据同步-设置Key过期
	if owlconfig.OwlConfigModel.GroupDataSync == "1" {
		//fmt.Println("gossip.Expire()")
		//发送数据到集群
		gossip.Expire(owlhandler.owlrequest.Key, owlhandler.owlrequest.Expires)
	}

}

func (owlhandler *OwlHandler) Delete() {

	ok := BaseCacheDB.Delete(owlhandler.owlrequest.Key)
	if ok {
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(ERROR)
	}

	//判断一致性数据同步-删除Key
	if owlconfig.OwlConfigModel.GroupDataSync == "1" {
		//fmt.Println("gossip.Delete()")
		//发送数据到集群
		gossip.Delete(owlhandler.owlrequest.Key)
	}

}

//PASS命令验证密码
func (owlhandler *OwlHandler) Pass(r *http.Request) {

	if owlconfig.OwlConfigModel.Pass == owlhandler.owlrequest.Pass {
		//token=md5(ip+uuid)
		uuid := owltools.GetUUIDString()
		ip := owltools.RemoteAddr2IPAddr(r.RemoteAddr)
		token := owltools.GetMd5String(ip + uuid)
		expiration, _ := time.ParseDuration(owlconfig.OwlConfigModel.Tonken_expire_time + "s")
		BaseAuth.Set(token, []byte(ip), expiration)
		//在返回值中添加UUID返回
		owlhandler.owlresponse.Data = []byte(token)
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(ERROR)
	}

}

//验证权限
func (owlhandler *OwlHandler) CheckAuth(r *http.Request) bool {

	token := owlhandler.owlrequest.Token
	ip := owltools.RemoteAddr2IPAddr(r.RemoteAddr)
	v, found := BaseAuth.Get(token)
	if found == true {
		if string(v) == ip {
			return true
		}
		return false
	}
	return false

}
