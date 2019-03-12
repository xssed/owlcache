package network

import (
	"net/http"
	//"time"
	//owlconfig "github.com/xssed/owlcache/config"
	//"github.com/xssed/owlcache/network"
	//tools "github.com/xssed/owlcache/tools"
)

//一个请求只产生一个 OwlServerGroupHandler
type OwlServerGroupHandler struct {
	owlservergrouprequest  *OwlServerGroupRequest
	owlserveggroupresponse *OwlServerGroupResponse
}

func NewOwlServerGroupHandler() *OwlServerGroupHandler {
	return &OwlServerGroupHandler{&OwlServerGroupRequest{}, &OwlServerGroupResponse{}}
}

//http服务器组执行数据操作
func (owlservergrouphandler *OwlServerGroupHandler) HTTPServerGroupHandle(w http.ResponseWriter, r *http.Request) {

	//	req := owlservergrouphandler.owlservergrouprequest

	//	//验证身份
	//	if !owlservergrouphandler.CheckAuth(r) {
	//		owlservergrouphandler.Transmit(NOT_PASS)
	//		return
	//	}

	//	command := CommandType(req.Cmd)

	//	switch command {
	//	case ADD:
	//		//owlservergrouphandler.Get()
	//	case DELETE:
	//		//owlservergrouphandler.Exists()
	//	case GetAll:

	//	case Get:

	//	default:
	//		owlservergrouphandler.Transmit(UNKNOWN_COMMAND)
	//	}

}

//解析response
func (owlservergrouphandler *OwlServerGroupHandler) Transmit(resstatus ResStatus) {

	switch resstatus {
	case SUCCESS:
		owlservergrouphandler.owlserveggroupresponse.Status = SUCCESS
		owlservergrouphandler.owlserveggroupresponse.Results = ResStatusToString(SUCCESS)
	case ERROR:
		owlservergrouphandler.owlserveggroupresponse.Status = ERROR
		owlservergrouphandler.owlserveggroupresponse.Results = ResStatusToString(ERROR)
	case NOT_FOUND:
		owlservergrouphandler.owlserveggroupresponse.Status = NOT_FOUND
		owlservergrouphandler.owlserveggroupresponse.Results = ResStatusToString(NOT_FOUND)
	case UNKNOWN_COMMAND:
		owlservergrouphandler.owlserveggroupresponse.Status = UNKNOWN_COMMAND
		owlservergrouphandler.owlserveggroupresponse.Results = ResStatusToString(UNKNOWN_COMMAND)
	case NOT_PASS:
		owlservergrouphandler.owlserveggroupresponse.Status = NOT_PASS
		owlservergrouphandler.owlserveggroupresponse.Results = ResStatusToString(NOT_PASS)
	}

	owlservergrouphandler.owlserveggroupresponse.Cmd = owlservergrouphandler.owlservergrouprequest.Cmd
	owlservergrouphandler.owlserveggroupresponse.Address = owlservergrouphandler.owlservergrouprequest.Address

}

//验证权限
func (owlservergrouphandler *OwlServerGroupHandler) CheckAuth(r *http.Request) bool {

	//	token := owlservergrouphandler.owlservergrouprequest.Token
	//	ip := tools.RemoteAddr2IPAddr(r.RemoteAddr)
	//	v, found := owlnetwork.BaseAuth.Get(token)
	//	if found == true {
	//		if v == ip {
	//			return true
	//		}
	//		return false
	//	}
	return false

}

//func (owlservergrouphandler *OwlServerGroupHandler) Set() {
//	ok := network.ServerGroupList.
//	if ok {
//		owlhandler.Transmit(SUCCESS)
//	} else {
//		owlhandler.Transmit(ERROR)
//	}
//}

//func (owlhandler *OwlHandler) Expire() {
//	ok := BaseCacheDB.Expire(owlhandler.owlrequest.Key, owlhandler.owlrequest.Expires)
//	if ok {
//		owlhandler.Transmit(SUCCESS)
//	} else {
//		owlhandler.Transmit(ERROR)
//	}
//}

//func (owlhandler *OwlHandler) Get() {
//	if v, found := BaseCacheDB.Get(owlhandler.owlrequest.Key); found {
//		owlhandler.Transmit(SUCCESS)
//		owlhandler.owlresponse.Data = v
//	} else {
//		owlhandler.Transmit(NOT_FOUND)
//	}
//}

//func (owlhandler *OwlHandler) Delete() {
//	ok := BaseCacheDB.Delete(owlhandler.owlrequest.Key)
//	if ok {
//		owlhandler.Transmit(SUCCESS)
//	} else {
//		owlhandler.Transmit(ERROR)
//	}
//}

//func (owlhandler *OwlHandler) Exists() {
//	ok := BaseCacheDB.Exists(owlhandler.owlrequest.Key)
//	if ok {
//		owlhandler.Transmit(SUCCESS)
//	} else {
//		owlhandler.Transmit(NOT_FOUND)
//	}
//}

////PASS命令验证密码
//func (owlhandler *OwlHandler) Pass(r *http.Request) {

//	if owlconfig.OwlConfigModel.Pass == owlhandler.owlrequest.Pass {
//		//token=md5(ip+uuid)
//		uuid := tools.GetUUIDString()
//		ip := tools.RemoteAddr2IPAddr(r.RemoteAddr)
//		token := tools.GetMd5String(ip + uuid)
//		expiration, _ := time.ParseDuration("1800s")
//		BaseAuth.Set(token, ip, expiration) //30分钟过期
//		//在返回值中添加UUID返回
//		owlhandler.owlresponse.Data = token
//		owlhandler.Transmit(SUCCESS)
//	} else {
//		owlhandler.Transmit(ERROR)
//	}

//}
