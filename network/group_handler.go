package network

import (
	"net/http"

	"github.com/xssed/owlcache/group"
	owllog "github.com/xssed/owlcache/log"
	tools "github.com/xssed/owlcache/tools"
)

//一个请求只产生一个 OwlServerGroupHandler
type OwlServerGroupHandler struct {
	owlservergrouprequest  *group.OwlServerGroupRequest
	owlserveggroupresponse *group.OwlServerGroupResponse
}

func NewOwlServerGroupHandler() *OwlServerGroupHandler {
	return &OwlServerGroupHandler{&group.OwlServerGroupRequest{}, &group.OwlServerGroupResponse{}}
}

//http服务器组执行数据操作
func (owlservergrouphandler *OwlServerGroupHandler) HTTPServerHandle(w http.ResponseWriter, r *http.Request) {

	//验证身份
	if !owlservergrouphandler.CheckAuth(r) {
		owlservergrouphandler.Transmit(group.NOT_PASS)
		return
	}

	req := owlservergrouphandler.owlservergrouprequest

	command := group.GroupCommandType(req.Cmd)

	switch command {
	case group.GroupADD:
		owlservergrouphandler.Add()
	case group.GroupDELETE:
		owlservergrouphandler.Delete()
	case group.GroupGetAll:
		owlservergrouphandler.GetAll()
	case group.GroupGet:
		owlservergrouphandler.Get()
	default:
		owlservergrouphandler.Transmit(group.UNKNOWN_COMMAND)
	}

}

//http服务器组执行数据操作,集群
func (owlservergrouphandler *OwlServerGroupHandler) HTTPServerGroupHandle(w http.ResponseWriter, r *http.Request) {

	//验证身份
	if !owlservergrouphandler.CheckAuth(r) {
		owlservergrouphandler.Transmit(group.NOT_PASS)
		return
	}

	req := owlservergrouphandler.owlservergrouprequest

	command := group.GroupCommandType(req.Cmd)

	switch command {
	case group.GroupADD:
		owlservergrouphandler.Add()
	case group.GroupDELETE:
		owlservergrouphandler.Delete()
	case group.GroupGetAll:
		owlservergrouphandler.GetAll()
	case group.GroupGet:
		owlservergrouphandler.Get()
	default:
		owlservergrouphandler.Transmit(group.UNKNOWN_COMMAND)
	}

}

//解析response
func (owlservergrouphandler *OwlServerGroupHandler) Transmit(resstatus group.ResStatus) {

	switch resstatus {
	case group.SUCCESS:
		owlservergrouphandler.owlserveggroupresponse.Status = group.SUCCESS
		owlservergrouphandler.owlserveggroupresponse.Results = group.ResStatusToString(group.SUCCESS)
	case group.ERROR:
		owlservergrouphandler.owlserveggroupresponse.Status = group.ERROR
		owlservergrouphandler.owlserveggroupresponse.Results = group.ResStatusToString(group.ERROR)
	case group.NOT_FOUND:
		owlservergrouphandler.owlserveggroupresponse.Status = group.NOT_FOUND
		owlservergrouphandler.owlserveggroupresponse.Results = group.ResStatusToString(group.NOT_FOUND)
	case group.UNKNOWN_COMMAND:
		owlservergrouphandler.owlserveggroupresponse.Status = group.UNKNOWN_COMMAND
		owlservergrouphandler.owlserveggroupresponse.Results = group.ResStatusToString(group.UNKNOWN_COMMAND)
	case group.NOT_PASS:
		owlservergrouphandler.owlserveggroupresponse.Status = group.NOT_PASS
		owlservergrouphandler.owlserveggroupresponse.Results = group.ResStatusToString(group.NOT_PASS)
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
	//创建数据
	var reqs group.OwlServerGroupRequest
	reqs.Cmd = owlservergrouphandler.owlservergrouprequest.Cmd
	reqs.Address = owlservergrouphandler.owlservergrouprequest.Address
	reqs.Pass = owlservergrouphandler.owlservergrouprequest.Pass
	reqs.Token = owlservergrouphandler.owlservergrouprequest.Token

	at, exits := owlservergrouphandler.find(owlservergrouphandler.owlservergrouprequest.Address)
	//存在
	if exits {
		ServerGroupList.RemoveAt(int32(at))          //先删除
		ok := ServerGroupList.AddAt(int32(at), reqs) //后增加
		if ok {
			owlservergrouphandler.Transmit(group.SUCCESS)
		} else {
			owlservergrouphandler.Transmit(group.ERROR)
		}
	} else {
		//不存在
		ok := ServerGroupList.Add(reqs)
		if ok {
			owlservergrouphandler.Transmit(group.SUCCESS)
		} else {
			owlservergrouphandler.Transmit(group.ERROR)
		}
	}

}

//内部查找一个服务器信息
func (owlservergrouphandler *OwlServerGroupHandler) find(address string) (int32, bool) {
	var resat int32 = 0
	resbool := false
	list := ServerGroupList.Values()
	for k, _ := range list {

		defer func() {
			if err := recover(); err != nil {
				owllog.OwlLogHttp.Panicln(err)
			}
		}()

		//fmt.Println(fmt.Sprintf("%T", list[k]))
		val, ok := list[k].(group.OwlServerGroupRequest)
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
			owlservergrouphandler.Transmit(group.SUCCESS)
		} else {
			owlservergrouphandler.Transmit(group.ERROR)
		}
	} else {
		//不存在
		owlservergrouphandler.Transmit(group.NOT_FOUND)
	}
}

//获取所有服务器列表信息
func (owlservergrouphandler *OwlServerGroupHandler) GetAll() {
	list := ServerGroupList.Values()
	owlservergrouphandler.owlserveggroupresponse.Data = list
	owlservergrouphandler.Transmit(group.SUCCESS)
}

//获取一个服务器信息
func (owlservergrouphandler *OwlServerGroupHandler) Get() {

	at, exits := owlservergrouphandler.find(owlservergrouphandler.owlservergrouprequest.Address)
	if exits {
		res, ok := ServerGroupList.GetAt(int32(at))
		if ok {
			owlservergrouphandler.owlserveggroupresponse.Data = res
			owlservergrouphandler.Transmit(group.SUCCESS)
		} else {
			owlservergrouphandler.Transmit(group.ERROR)
		}
	} else {
		//不存在
		owlservergrouphandler.Transmit(group.NOT_FOUND)
	}

}
