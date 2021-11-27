package network

import (
	"errors"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
)

func gossip_set(key, val, expire string) {

	//key_resource := owlconfig.OwlConfigModel.ResponseHost + ":" + owlconfig.OwlConfigModel.Tcpport
	data, err := gossip_getUrlData(key, val)

	exptime, _ := time.ParseDuration(owltools.JoinString(expire, "s"))
	ok := BaseCacheDB.Set(key, []byte(val), exptime)
	if !ok {
		owllog.OwlLogHttp.Println(owltools.JoinString("gossip:set error key:", key))
	}
}

func gossip_del(key string) {
	ok := BaseCacheDB.Delete(key)
	if !ok {
		owllog.OwlLogHttp.Println(owltools.JoinString("gossip:del error key:", key))
	}
}

func gossip_expire(key, expire string) {
	exptime, _ := time.ParseDuration(owltools.JoinString(expire, "s"))
	ok := BaseCacheDB.Expire(key, exptime)
	if !ok {
		owllog.OwlLogHttp.Println(owltools.JoinString("gossip:expire error key:", key))
	}
}

func gossip_getUrlData(key string, val string) ([]byte, error) {

	//创建http client
	var grsa *gorequest.SuperAgent
	grsa = gorequest.New()
	grsa.Get(owltools.JoinString(val, "/data"))
	grsa.Param("cmd", "get")
	grsa.Param("key", key)
	//设置超时
	ghcrt, _ := strconv.Atoi(owlconfig.OwlConfigModel.GossipHttpClientRequestTimeout)
	grsa.Timeout(time.Duration(ghcrt) * time.Millisecond)
	//发送请求获取数据
	resp, _, errs := grsa.EndBytes()
	if errs != nil {
		errstr := owltools.ErrorSliceJoinToString(errs)
		if errstr != "" {
			er := owltools.JoinString("Gossip get url data error:", errstr, " url:", val)
			owllog.OwlLogHttp.Info(er) //日志记录
			return []byte(""), errors.New(er)
		}
	}
	if resp.StatusCode != 200 {
		return []byte(""), ErrorGossipGetUrlData
	}
	defer resp.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(resp.Body)
	if ioerr != nil {
		owllog.OwlLogHttp.Info(owltools.JoinString("Gossip get url data ioutil.ReadAll error:", ioerr.Error())) //日志记录
		return []byte(""), ioerr.Error()
	}

	//清理资源
	grsa.ClearSuperAgent()

	return body, nil

}
