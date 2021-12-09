package httpclient

import (
	"time"

	"github.com/parnurzeal/gorequest"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
)

//定义HTTP客户端结构
type OwlClient struct {
	*OwlHttp
}

//创建HttpClient实体
func NewOwlClient() *OwlClient {

	return &OwlClient{NewOwlHttpClient()}

}

//获取Key值
func (c *OwlClient) GetValue(address, key string) *Response {

	//判断集群之间是否开启HTTPS安全通道
	if owlconfig.OwlConfigModel.HttpsClient_InsecureSkipVerify == "0" {
		return c.primitive_get(address, key)
	}
	return c.modern_get(address, key)

}

//获取Key值-自己的HttpClient封装包
func (c *OwlClient) primitive_get(address, key string) *Response {

	owlclient := c
	owlclient.PostForm(owltools.JoinString(address, "/data/"))
	owlclient.SetTimeout(c.HCRequestTimeout * time.Millisecond)
	owlclient.Query.Add("cmd", "get")
	owlclient.Query.Add("key", key)
	res, err := owlclient.Do()
	if err != nil {
		owllog.OwlLogHttpG.Info("owlclient method GetValue error:" + err.Error()) //日志记录
	}
	//owllog.OwlLogHttpG.Info("HTTP request OK："+address, key) //日志记录
	owlclient.Claer() //清空数据
	if res != nil && res.StatusCode == 200 {
		return res
	} else {
		return nil
	}

}

//获取Key值-开源的HttpClient封装包gorequest
func (c *OwlClient) modern_get(address, key string) *Response {

	var grsa *gorequest.SuperAgent
	grsa = gorequest.New()

	grsa.Timeout(c.HCRequestTimeout * time.Millisecond)
	grsa.Get(owltools.JoinString(address, "/data/"))
	grsa.Param("cmd", "get")
	grsa.Param("key", key)

	//发送请求获取数据
	r_res, _, r_err_slices := grsa.EndBytes()
	//判断请求的响应数据是否超过本地允许的最大值
	if r_err_slices != nil {
		owllog.OwlLogHttpG.Info(owltools.JoinString("owlclient method GetValue error:", owltools.ErrorSliceJoinToString(r_err_slices))) //日志记录
	}
	//清理资源
	grsa.ClearSuperAgent()
	if r_res != nil && r_res.StatusCode == 200 {
		return &Response{r_res}
	} else {
		return nil
	}

}
