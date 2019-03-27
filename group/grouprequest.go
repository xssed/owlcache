package group

import (
	//"bufio"
	"fmt"
	//"io"
	//"strconv"
	"net/http"
	"strings"
	//"time"
)

type OwlServerGroupRequest struct {
	//请求命令
	Cmd GroupCommandType
	//地址字符串
	Address string
	//链接密码
	Pass string
	//token
	Token string
}

//request to string
func (req *OwlServerGroupRequest) String() string {
	return fmt.Sprintf("{OwlServerGroupRequest cmd=%s , address='%s' , pass=%d }",
		req.Cmd, req.Address, req.Pass)
}

//过滤接收数据中的\r\n
func (req *OwlServerGroupRequest) TrimSpace(str string) string {
	if str != "" {
		return strings.TrimSpace(str)
	}
	return ""
}

//将http请求内容 解析为一个OwlServerGroupRequest对象
func (req *OwlServerGroupRequest) HTTPReceive(w http.ResponseWriter, r *http.Request) {

	//判断空字符串请求
	if len(r.Form) <= 1 && req.TrimSpace(r.FormValue("cmd")) == "" {
		return
	}

	req.Cmd = GroupCommandType(r.FormValue("cmd"))
	req.Address = r.FormValue("address")
	req.Pass = r.FormValue("pass")
	req.Token = r.FormValue("token")

}
