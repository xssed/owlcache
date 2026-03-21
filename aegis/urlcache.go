package aegis

import (
	"fmt"
	"os"
	"strconv"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

// 检查urlcache的配置
func CheckUCConfig() {

	//检查是否启动URL缓存功能
	if owlconfig.OwlConfigModel.Open_Urlcache == "1" {
		checkSitesConfig()
	}

}

// 检查站点配置
func checkSitesConfig() {
	for i, site := range owlconfig.OwlUCConfigModel.SiteList {
		index := strconv.Itoa(i)
		checkSiteConfig(site, index)
	}
}

// 检查单个站点配置
func checkSiteConfig(site interface{}, index string) {
	siteConfig, ok := site.(struct {
		Host    string
		Headers []struct {
			Value []string
		}
		Proxy string
	})
	if !ok {
		return
	}

	checkHostConfig(siteConfig.Host, index)
	checkHeadersConfig(siteConfig.Headers, index)
	checkProxyConfig(siteConfig.Proxy, index)
}

// 检查Host配置
func checkHostConfig(host, index string) {
	if len(host) <= 7 {
		reportConfigError("sites-site-host", index)
	}
}

// 检查Headers配置
func checkHeadersConfig(headers []struct {
	Value []string
}, index string) {
	for _, header := range headers {
		for _, headerValue := range header.Value {
			if len(headerValue) == 0 {
				reportConfigError("sites-site-header-value", index)
			}
		}
	}
}

// 检查Proxy配置
func checkProxyConfig(proxy, index string) {
	if len(proxy) != 0 && len(proxy) < 7 {
		reportConfigError("sites-site-proxy", index)
	}
}

// 报告配置错误
func reportConfigError(configItem, index string) {
	owllog.OwlLogRun.Println(fmt.Sprintf("Please set a value for the <%s> item whose index is [%s]. The url cache configuration file path is %s.", configItem, index, owlconfig.OwlConfigModel.DBfile+owlconfig.OwlConfigModel.Urlcache_Filename))
	os.Exit(0)
}
