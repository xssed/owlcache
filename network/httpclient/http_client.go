package httpclient

import (
	"net/http"
	"os"
	"strconv"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

//定义HTTP客户端结构
type OwlClient struct {
	OwlTransport     *http.Transport
	HCRequestTimeout time.Duration
}

//创建一个HTTP客户端
func NewOwlClient() *OwlClient {
	//创建Transport()
	owltransport := NewOwlTransport()
	//从配置中取出集群互相通信时的请求超时时间
	//hcrequesttimeout, err := time.ParseDuration(owlconfig.OwlConfigModel.HttpClientRequestTimeout + "s")//bug
	hcrequesttimeout, err := strconv.Atoi(owlconfig.OwlConfigModel.HttpClientRequestTimeout)
	if err != nil {
		//强制异常，退出
		owllog.OwlLogHttp.Info("Config File HttpClientRequestTimeout Parse error:" + err.Error()) //日志记录
		//fmt.Println("Config File HttpClientRequestTimeout Parse error:" + err.Error())
		os.Exit(0)
	}

	owlhttpclient := &OwlClient{OwlTransport: owltransport, HCRequestTimeout: time.Duration(hcrequesttimeout)}

	return owlhttpclient
}

//获取Key值
func (c *OwlClient) GetValue(address, key string) *Response {

	owlclient := NewOwlHttpClient(c.OwlTransport)
	owlclient.PostForm(address + "/data")
	owlclient.SetTimeout(c.HCRequestTimeout * time.Millisecond)
	owlclient.Query.Add("cmd", "get")
	owlclient.Query.Add("key", key)
	res, err := owlclient.Do()
	if err != nil {
		owllog.OwlLogHttp.Info("owlclient method GetValue error:" + err.Error()) //日志记录
	}
	//owllog.OwlLogHttp.Info("HTTP request OK："+address, key) //日志记录
	owlclient.Claer() //清空查询数据
	if res != nil && res.StatusCode == 200 {
		return res
	} else {
		return nil
	}

}
