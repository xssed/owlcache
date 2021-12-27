package network

import (
	"net/http"
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
	case PING:
		owlhandler.Ping()
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
	case PING:
		owlhandler.Ping()
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

//Websocket数据执行操作
func (owlhandler *OwlHandler) WebsocketHandle(w http.ResponseWriter, r *http.Request) {

	req := owlhandler.owlrequest

	command := CommandType(req.Cmd)

	switch command {
	case GET:
		// if string(req.Value) == "data" {
		owlhandler.Get()
		// } else {
		// 	owlhandler.Transmit(UNKNOWN_COMMAND)
		// }
	default:
		owlhandler.Transmit(UNKNOWN_COMMAND)
	}

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
