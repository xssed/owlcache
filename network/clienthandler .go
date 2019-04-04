package network

import (
	"fmt"
	//"log"
	//"net/http"

	"github.com/xssed/owlcache/group"
	//tools "github.com/xssed/owlcache/tools"
)

//发起请求
func (owlhandler *OwlHandler) GetGroupData() { //r *http.Request

	list := ServerGroupList.Values()
	//fmt.Println(list)
	for k := range list {
		val, ok := list[k].(group.OwlServerGroupRequest)
		if ok {
			fmt.Println(HttpClient.GetValue(val.Address, owlhandler.owlrequest.Key))
			//fmt.Println(val)
		}
	}

	//	HttpClient

	//	token := owlhandler.owlrequest.Token
	//	ip := tools.RemoteAddr2IPAddr(r.RemoteAddr)
	//	v, found := BaseAuth.Get(token)
	//	if found == true {
	//		if v == ip {
	//			return true
	//		}
	//		return false
	//	}
	//	return false

	//	owlservergrouphandler.Transmit(group.NOT_FOUND)

}
