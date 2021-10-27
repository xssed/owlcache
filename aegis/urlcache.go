package aegis

import (
	"os"
	"strconv"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

//检查urlcache的配置
func CheckUCConfig() {

	//检查是否启动URL缓存功能
	if owlconfig.OwlConfigModel.Open_Urlcache == "1" {

		for i, site := range owlconfig.OwlUCConfigModel.SiteList {
			index := strconv.Itoa(i)
			//Host
			if len(site.Host) <= 7 {
				owllog.OwlLogRun.Println("Please set a value for the <sites>-<site>-<host> item whose index is [" + index + "]. The url cache configuration file path is " + owlconfig.OwlConfigModel.DBfile + owlconfig.OwlConfigModel.Urlcache_Filename + ".")
				os.Exit(0)
			}
			//Headers
			if len(site.Headers) != 0 {
				for _, rhs := range site.Headers {
					if len(rhs.Value) == 0 {
						owllog.OwlLogRun.Println("Please set a value for the <sites>-<site>-<header>-<value> item whose index is [" + index + "]. The url cache configuration file path is " + owlconfig.OwlConfigModel.DBfile + owlconfig.OwlConfigModel.Urlcache_Filename + ".")
						os.Exit(0)
					} else {
						for _, rh := range rhs.Value {
							if len(rh) == 0 {
								owllog.OwlLogRun.Println("Please set a value for the <sites>-<site>-<header>-<value> item whose index is [" + index + "]. The url cache configuration file path is " + owlconfig.OwlConfigModel.DBfile + owlconfig.OwlConfigModel.Urlcache_Filename + ".")
								os.Exit(0)
							}
						}
					}
				}
			}
			//Proxy
			if len(site.Proxy) != 0 && len(site.Proxy) < 7 {
				owllog.OwlLogRun.Println("Please set a value for the <sites>-<site>-<proxy> item whose index is [" + index + "]. The url cache configuration file path is " + owlconfig.OwlConfigModel.DBfile + owlconfig.OwlConfigModel.Urlcache_Filename + ".")
				os.Exit(0)
			}
		}

	}

}
