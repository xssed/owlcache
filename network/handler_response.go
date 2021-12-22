package network

import (
	"encoding/json"
	"net/http"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
	owltools "github.com/xssed/owlcache/tools"
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
	Data []byte
	//程序响应IP
	ResponseHost string
	//内容的创建时间
	KeyCreateTime time.Time
}

//将数据转换成json(单节点)
func (owlhandler *OwlHandler) ToHttp(w http.ResponseWriter) (http.ResponseWriter, []byte) {

	owlhandler.owlresponse.ResponseHost = owltools.JoinString(owlconfig.OwlConfigModel.ResponseHost, ":", owlconfig.OwlConfigModel.Httpport) //设置响应的主机信息
	//设置Ke的响应信息
	w.Header().Set("ResponseHost", owlhandler.owlresponse.ResponseHost)
	w.Header().Set("Key", owlhandler.owlresponse.Key)
	w.Header().Set("KeyCreateTime", owlhandler.owlresponse.KeyCreateTime.String())
	//GET、PING命令请求优先处理
	if owlhandler.owlrequest.Cmd == GET || owlhandler.owlrequest.Cmd == PING {
		return w, owlhandler.owlresponse.Data
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	data, _ := json.Marshal(owlhandler.owlresponse)
	return w, data

}

//将数据转换成json(集群)
func (owlhandler *OwlHandler) ToGroupHttp(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, []byte) {

	//如果valuedata不是字符串info则输出集群第一个
	if string(owlhandler.owlrequest.Value) != "info" {
		return w, owlhandler.owlresponse.Data
	}
	//查询info类型
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	return w, owlhandler.owlresponse.Data

}

//TCP服务将数据进行转换输出
func (owlhandler *OwlHandler) ToTcp() []byte {

	if owlhandler.owlrequest.Cmd == GET {
		if string(owlhandler.owlrequest.Value) != "info" {
			return owlhandler.owlresponse.Data
		}
		owlhandler.owlresponse.Data = []byte("")
	}
	//PING命令
	if owlhandler.owlrequest.Cmd == PING {
		return owlhandler.owlresponse.Data
	}
	owlhandler.owlresponse.ResponseHost = owlconfig.OwlConfigModel.ResponseHost + ":" + owlconfig.OwlConfigModel.Tcpport
	data, _ := json.Marshal(owlhandler.owlresponse)
	return data

}

//Websocket服务将数据进行转换输出
func (owlhandler *OwlHandler) ToWebsocket() []byte {

	if owlhandler.owlrequest.Cmd == GET && string(owlhandler.owlrequest.Value) != "" {
		return owlhandler.owlresponse.Data
	}
	owlhandler.owlresponse.Data = []byte("")
	data, _ := json.Marshal(owlhandler.owlresponse)
	return data

}
