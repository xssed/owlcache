package memcacheclient

import (
	"fmt"
	"os"

	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

var mc *memcache.Client

//客户端初始化
func Start() {

	if len(owlconfig.OwlConfigModel.Memcache_list) < 5 {
		fmt.Println("Configuration file Memcache_list filled in error")
		os.Exit(0)
	}
	list := strings.Split(owlconfig.OwlConfigModel.Memcache_list, "|")
	mc = memcache.New(list...)

}

//Get方法
func Get(key string) (*memcache.Item, error) {
	it, err := mc.Get(key)
	if err != nil {
		owllog.OwlLogRun.Println("Get() error:", err)
		return nil, err
	}
	return it, nil
}
