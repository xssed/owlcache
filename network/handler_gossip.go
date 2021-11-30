package network

import (
	"errors"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
)

//处理gossip集群发来的更新命令
func gossip_set(key, val, expire string) {

	data, err := gossip_getUrlData(key, val)
	if err != nil {
		owllog.OwlLogGossip.Info(err.Error()) //日志记录
	}
	//请求成功
	owllog.OwlLogGossip.Info(owltools.JoinString("Gossip get url data success: key:", key, " url:", val)) //日志记录

	exptime, _ := time.ParseDuration(owltools.JoinString(expire, "s"))
	ok := BaseCacheDB.Set(key, data, exptime)
	if !ok {
		owllog.OwlLogGossip.Info(owltools.JoinString("gossip:set error key:", key))
	}
	owllog.OwlLogGossip.Info(owltools.JoinString("gossip:set key success:", key)) //日志记录

}

//处理gossip集群发来的删除命令
func gossip_del(key string) {
	ok := BaseCacheDB.Delete(key)
	if !ok {
		owllog.OwlLogGossip.Info(owltools.JoinString("gossip:del error key:", key))
	}
}

//处理gossip集群发来的设置key过期命令
func gossip_expire(key, expire string) {
	exptime, _ := time.ParseDuration(owltools.JoinString(expire, "s"))
	ok := BaseCacheDB.Expire(key, exptime)
	if !ok {
		owllog.OwlLogGossip.Info(owltools.JoinString("gossip:expire error key:", key))
	}
}

//在接收到gossip集群发来的更新命令之后，根据value(存放的目标服务IP地址)发起一个HTTP请求，取出数据，存放到本地数据库
func gossip_getUrlData(key string, val string) ([]byte, error) {

	//创建http client
	var grsa *gorequest.SuperAgent
	grsa = gorequest.New()
	grsa.Get(owltools.JoinString(val, "/data/"))
	grsa.Param("cmd", "get")
	grsa.Param("key", key)
	//设置超时
	ghcrt, _ := strconv.Atoi(owlconfig.OwlConfigModel.GossipHttpClientRequestTimeout)
	//设置毫秒超时
	grsa.Timeout(time.Duration(ghcrt) * time.Millisecond)
	//发送请求获取数据
	resp, _, errs := grsa.EndBytes()
	if errs != nil {
		errstr := owltools.ErrorSliceJoinToString(errs)
		if errstr != "" {
			return []byte(""), errors.New(owltools.JoinString("Gossip get url data error:", errstr, " url:", val))
		}
	}
	if resp.StatusCode != 200 {
		return []byte(""), ErrorGossipGetUrlData
	}
	defer resp.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(resp.Body)
	if ioerr != nil {
		return []byte(""), errors.New(owltools.JoinString("Gossip get url data ioutil.ReadAll error:", ioerr.Error()))
	}

	//清理资源
	grsa.ClearSuperAgent()

	return body, nil

}
