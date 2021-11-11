package network

import (
	"encoding/json"
	"net/http"

	owlconfig "github.com/xssed/owlcache/config"
	owltools "github.com/xssed/owlcache/tools"
)

//一个请求只产生一个 OwlHandler
type OwlHandler struct {
	owlrequest  *OwlRequest
	owlresponse *OwlResponse
}

func NewOwlHandler() *OwlHandler {
	return &OwlHandler{&OwlRequest{}, &OwlResponse{}}
}

//TCP执行数据操作
func (owlhandler *OwlHandler) TCPHandle() {

	req := owlhandler.owlrequest

	command := CommandType(req.Cmd)

	switch command {
	case GET:
		owlhandler.Get()
	case EXIST:
		owlhandler.Exists()
	case SET:
		owlhandler.Set()
	case EXPIRE:
		owlhandler.Expire()
	case DELETE:
		owlhandler.Delete()
	default:
		owlhandler.Transmit(UNKNOWN_COMMAND)
	}

}

//http单机执行数据操作
func (owlhandler *OwlHandler) HTTPHandle(w http.ResponseWriter, r *http.Request) {

	req := owlhandler.owlrequest

	command := CommandType(req.Cmd)

	switch command {
	case GET:
		owlhandler.Get()
	case EXIST:
		owlhandler.Exists()
	case SET:
		if !owlhandler.CheckAuth(r) {
			owlhandler.Transmit(NOT_PASS)
			break
		}
		owlhandler.Set()
	case EXPIRE:
		if !owlhandler.CheckAuth(r) {
			owlhandler.Transmit(NOT_PASS)
			break
		}
		owlhandler.Expire()
	case DELETE:
		if !owlhandler.CheckAuth(r) {
			owlhandler.Transmit(NOT_PASS)
			break
		}
		owlhandler.Delete()
	case PASS:
		owlhandler.Pass(r)
	default:
		owlhandler.Transmit(UNKNOWN_COMMAND)
	}

}

//http群组执行数据操作
func (owlhandler *OwlHandler) HTTPGroupDataHandle(w http.ResponseWriter, r *http.Request) {

	req := owlhandler.owlrequest

	command := CommandType(req.Cmd)

	switch command {
	case GET:
		//HttpClient
		owlhandler.GetGroupData(w, r)
	default:
		owlhandler.Transmit(UNKNOWN_COMMAND)
	}

}

//UrlCache数据执行信息
func (owlhandler *OwlHandler) UCDataHandle(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, []byte) {

	w, print := owlhandler.GeUrlCacheData(w, r)
	return w, print

}

//解析response
func (owlhandler *OwlHandler) Transmit(resstatus ResStatus) {

	switch resstatus {
	case SUCCESS:
		owlhandler.owlresponse.Status = SUCCESS
		owlhandler.owlresponse.Results = ResStatusToString(SUCCESS)
	case ERROR:
		owlhandler.owlresponse.Status = ERROR
		owlhandler.owlresponse.Results = ResStatusToString(ERROR)
	case NOT_FOUND:
		owlhandler.owlresponse.Status = NOT_FOUND
		owlhandler.owlresponse.Results = ResStatusToString(NOT_FOUND)
	case UNKNOWN_COMMAND:
		owlhandler.owlresponse.Status = UNKNOWN_COMMAND
		owlhandler.owlresponse.Results = ResStatusToString(UNKNOWN_COMMAND)
	case NOT_PASS:
		owlhandler.owlresponse.Status = NOT_PASS
		owlhandler.owlresponse.Results = ResStatusToString(NOT_PASS)
	}

	owlhandler.owlresponse.Cmd = owlhandler.owlrequest.Cmd
	owlhandler.owlresponse.Key = owlhandler.owlrequest.Key

}

//将数据转换成json(单节点)
func (owlhandler *OwlHandler) ToHttp(w http.ResponseWriter) (http.ResponseWriter, []byte) {

	owlhandler.owlresponse.ResponseHost = owltools.JoinString(owlconfig.OwlConfigModel.ResponseHost, ":", owlconfig.OwlConfigModel.Httpport) //设置响应的主机信息
	//设置Ke的响应信息
	w.Header().Set("ResponseHost", owlhandler.owlresponse.ResponseHost)
	w.Header().Set("Key", owlhandler.owlresponse.Key)
	w.Header().Set("KeyCreateTime", owlhandler.owlresponse.KeyCreateTime.String())
	//GET请求优先处理
	if owlhandler.owlrequest.Cmd == GET {
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
	owlhandler.owlresponse.ResponseHost = owlconfig.OwlConfigModel.ResponseHost + ":" + owlconfig.OwlConfigModel.Tcpport
	data, _ := json.Marshal(owlhandler.owlresponse)
	return data

}
