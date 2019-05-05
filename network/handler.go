package network

import (
	"net/http"
	"time"

	"github.com/xssed/owlcache/cache"
	owlconfig "github.com/xssed/owlcache/config"
	tools "github.com/xssed/owlcache/tools"
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
	case SET:
		owlhandler.Set()
	case EXPIRE:
		owlhandler.Expire()
	case GET:
		owlhandler.Get()
	case DELETE:
		owlhandler.Delete()
	case EXIST:
		owlhandler.Exists()
		//	case PASS:
		//		owlhandler.Pass()
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
		owlhandler.GetGroupData()
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

func (owlhandler *OwlHandler) Set() {
	ok := BaseCacheDB.Set(owlhandler.owlrequest.Key, owlhandler.owlrequest.Value, owlhandler.owlrequest.Expires)
	if ok {
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(ERROR)
	}
}

func (owlhandler *OwlHandler) Expire() {
	ok := BaseCacheDB.Expire(owlhandler.owlrequest.Key, owlhandler.owlrequest.Expires)
	if ok {
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(ERROR)
	}
}

func (owlhandler *OwlHandler) Get() {
	if v, found := BaseCacheDB.GetKvStore(owlhandler.owlrequest.Key); found {
		owlhandler.Transmit(SUCCESS)
		owlhandler.owlresponse.Data = v.(*cache.KvStore).Value
		owlhandler.owlresponse.KeyCreateTime = v.(*cache.KvStore).CreateTime
	} else {
		owlhandler.Transmit(NOT_FOUND)
	}
}

func (owlhandler *OwlHandler) Delete() {
	ok := BaseCacheDB.Delete(owlhandler.owlrequest.Key)
	if ok {
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(ERROR)
	}
}

func (owlhandler *OwlHandler) Exists() {
	ok := BaseCacheDB.Exists(owlhandler.owlrequest.Key)
	if ok {
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(NOT_FOUND)
	}
}

//PASS命令验证密码
func (owlhandler *OwlHandler) Pass(r *http.Request) {

	if owlconfig.OwlConfigModel.Pass == owlhandler.owlrequest.Pass {
		//token=md5(ip+uuid)
		uuid := tools.GetUUIDString()
		ip := tools.RemoteAddr2IPAddr(r.RemoteAddr)
		token := tools.GetMd5String(ip + uuid)
		expiration, _ := time.ParseDuration("3600s")
		BaseAuth.Set(token, ip, expiration) //60分钟过期
		//在返回值中添加UUID返回
		owlhandler.owlresponse.Data = token
		owlhandler.Transmit(SUCCESS)
	} else {
		owlhandler.Transmit(ERROR)
	}

}

//验证权限
func (owlhandler *OwlHandler) CheckAuth(r *http.Request) bool {

	//uuid := owlhandler.owlrequest.Pass
	token := owlhandler.owlrequest.Token
	ip := tools.RemoteAddr2IPAddr(r.RemoteAddr)
	v, found := BaseAuth.Get(token)
	if found == true {
		if v == ip {
			return true
		}
		return false
	}
	return false

}
