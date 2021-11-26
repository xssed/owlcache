package network

import (
	"time"

	//"github.com/parnurzeal/gorequest"
	//owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
)

func gossip_set(key, val, expire string) {

	//key_resource := owlconfig.OwlConfigModel.ResponseHost + ":" + owlconfig.OwlConfigModel.Tcpport

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
