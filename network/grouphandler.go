package network

import (
	"net/http"
	//"time"
	//owlconfig "github.com/xssed/owlcache/config"
	//"github.com/xssed/owlcache/network"
	//"fmt"

	tools "github.com/xssed/owlcache/tools"
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

	//验证身份
	if !owlservergrouphandler.CheckAuth(r) {
		owlservergrouphandler.Transmit(NOT_PASS)
		return
	}

	req := owlservergrouphandler.owlservergrouprequest

	command := GroupCommandType(req.Cmd)

	switch command {
	case GroupADD:
		owlservergrouphandler.Add()
	case GroupDELETE:
		owlservergrouphandler.Delete()
	case GroupGetAll:
		owlservergrouphandler.GetAll()
	case GroupGet:
		owlservergrouphandler.Get()
	default:
		owlservergrouphandler.Transmit(UNKNOWN_COMMAND)
	}

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

	token := owlservergrouphandler.owlservergrouprequest.Token
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

//添加一个服务器信息
func (owlservergrouphandler *OwlServerGroupHandler) Add() {

	//数据清理
	owlservergrouphandler.owlservergrouprequest.Cmd = ""

	at, exits := owlservergrouphandler.find(owlservergrouphandler.owlservergrouprequest.Address)
	//存在
	if exits {
		//		res := ServerGroupList.RemoveAt(int32(at))
		//		if res {
		//			owlservergrouphandler.Transmit(SUCCESS)
		//		} else {
		//			owlservergrouphandler.Transmit(ERROR)
		//		}
		ServerGroupList.RemoveAt(int32(at))                                                 //先删除
		ok := ServerGroupList.AddAt(int32(at), owlservergrouphandler.owlservergrouprequest) //后增加
		if ok {
			owlservergrouphandler.Transmit(SUCCESS)
		} else {
			owlservergrouphandler.Transmit(ERROR)
		}
	} else {
		//不存在
		ok := ServerGroupList.Add(owlservergrouphandler.owlservergrouprequest)
		if ok {
			owlservergrouphandler.Transmit(SUCCESS)
		} else {
			owlservergrouphandler.Transmit(ERROR)
		}
	}

}

//内部查找一个服务器信息
func (owlservergrouphandler *OwlServerGroupHandler) find(address string) (int32, bool) {
	var resat int32 = 0
	resbool := false
	list := ServerGroupList.Values()
	for k, _ := range list {
		val, ok := list[k].(*OwlServerGroupRequest)
		if ok {
			if val.Address == address {
				resat = int32(k)
				resbool = true
			}
		}
	}
	return resat, resbool
}

//删除一个服务器信息
func (owlservergrouphandler *OwlServerGroupHandler) Delete() {
	at, exits := owlservergrouphandler.find(owlservergrouphandler.owlservergrouprequest.Address)
	if exits {
		res := ServerGroupList.RemoveAt(int32(at))
		if res {
			owlservergrouphandler.Transmit(SUCCESS)
		} else {
			owlservergrouphandler.Transmit(ERROR)
		}
	} else {
		//不存在
		owlservergrouphandler.Transmit(NOT_FOUND)
	}
}

//获取所有服务器列表信息
func (owlservergrouphandler *OwlServerGroupHandler) GetAll() {
	list := ServerGroupList.Values()
	owlservergrouphandler.owlserveggroupresponse.Data = list
	owlservergrouphandler.Transmit(SUCCESS)
}

//获取一个服务器信息
func (owlservergrouphandler *OwlServerGroupHandler) Get() {

	at, exits := owlservergrouphandler.find(owlservergrouphandler.owlservergrouprequest.Address)
	if exits {
		res, ok := ServerGroupList.GetAt(int32(at))
		if ok {
			owlservergrouphandler.owlserveggroupresponse.Data = res
			owlservergrouphandler.Transmit(SUCCESS)
		} else {
			owlservergrouphandler.Transmit(ERROR)
		}
	} else {
		//不存在
		owlservergrouphandler.Transmit(NOT_FOUND)
	}

}
