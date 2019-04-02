package httpclient

import (
	"fmt"
	owllog "log"
	"net/http"
	"os"
	"strconv"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
	//owlnetwork "github.com/xssed/owlcache/network"
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
		owllog.Println("Config File HttpClientRequestTimeout Parse error：" + err.Error()) //日志记录
		fmt.Println("Config File HttpClientRequestTimeout Parse error：" + err.Error())
		os.Exit(0)
	}

	owlhttpclient := &OwlClient{OwlTransport: owltransport, HCRequestTimeout: time.Duration(hcrequesttimeout)}

	return owlhttpclient
}

//登录获取Token
func (c *OwlClient) GetToken(address, cmd, pass string) {
	owlclient := NewOwlHttpClient(c.OwlTransport)
	fmt.Println(owlclient)
	owlclient.Get(address)
	owlclient.SetTimeout(c.HCRequestTimeout * time.Second)
	owlclient.Query.Add("cmd", cmd)
	owlclient.Query.Add("pass", pass)
	res, err := owlclient.Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.StatusCode)
	fmt.Println(res.String())

	//owlnetwork.PASS

	owlclient.Claer()
}
