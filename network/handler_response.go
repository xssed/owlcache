package network

import (
	"encoding/json"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
)

type OwlResponse struct {
	//请求命令
	Cmd CommandType
	//返回状态
	Status ResStatus
	//返回结果
	Results string
	//key
	Key string
	//返回内容
	Data interface{}
	//程序响应IP
	ResponseHost string
	//内容的创建时间
	KeyCreateTime time.Time
	//响应内容格式
	ContentType string
}

//将数据转换成json
func (p *OwlResponse) ConvertToString(mode string) string {

	if mode == "TCP" {
		p.ResponseHost = owlconfig.OwlConfigModel.ResponseHost + ":" + owlconfig.OwlConfigModel.Tcpport
	} else if mode == "HTTP" {
		p.ResponseHost = owlconfig.OwlConfigModel.ResponseHost + ":" + owlconfig.OwlConfigModel.Httpport
	} else {
		p.ResponseHost = ""
	}

	data, _ := json.Marshal(p)
	s := string(data)
	return s
}
