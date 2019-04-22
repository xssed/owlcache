package network

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type OwlRequest struct {
	//请求命令
	Cmd CommandType
	//key
	Key string
	//请求内容
	Value interface{}
	//请求内容长度
	Length int
	//过期时间
	Expires time.Duration //int64
	//链接密码
	Pass string
	//token
	Token string
}

//request to string
func (req *OwlRequest) String() string {
	return fmt.Sprintf("{OwlRequest cmd=%s , key='%s' , bodylen=%d }",
		req.Cmd, req.Key, len(req.Value.(string)))
}

//过滤接收数据中的\r\n
func (req *OwlRequest) TrimSpace(str string) string {
	if str != "" {
		return strings.TrimSpace(str)
	}
	return ""
}

//将socket请求内容 解析为一个OwlRequest对象
func (req *OwlRequest) TCPReceive(connstr string) {

	params := strings.Split(connstr, " ") //strings.Fields(connstr)

	//判断空字符串请求
	if len(params) <= 1 && strings.TrimSpace(params[0]) == "" {
		return
	}

	command := CommandType(params[0])

	switch command {
	case SET:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
		req.Value = req.Slicetostring(params[2:])
		req.Length = len(req.Value.(string))
	case EXPIRE:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
		exptime, _ := time.ParseDuration(req.TrimSpace(params[2]) + "s")
		req.Expires = exptime
	case GET:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
	case DELETE:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
	case EXIST:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
		//	case PASS:
		//		req.Cmd = command
		//		req.Pass = req.TrimSpace(params[1])
	}

}

//将http请求内容 解析为一个OwlRequest对象
func (req *OwlRequest) HTTPReceive(w http.ResponseWriter, r *http.Request) {

	//r.ParseForm() //解析参数, 默认是不会解析的
	//fmt.Println(r.Form)
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)

	//判断空字符串请求
	if len(r.Form) <= 1 && strings.TrimSpace(r.FormValue("cmd")) == "" {
		return
	}

	req.Cmd = CommandType(r.FormValue("cmd"))
	req.Key = r.FormValue("key")
	req.Value = r.FormValue("valuedata")
	req.Length = len(req.Value.(string))
	exptime, _ := time.ParseDuration(req.TrimSpace(r.FormValue("exptime")) + "s")
	req.Expires = exptime
	req.Pass = r.FormValue("pass")
	req.Token = r.FormValue("token")

}

//将字符串切片转换成字符串
func (req *OwlRequest) Slicetostring(slice []string) string {
	var returnstr string
	for _, v := range slice {
		returnstr = returnstr + v + " "
	}
	return returnstr
}
