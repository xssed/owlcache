package network

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
	owltools "github.com/xssed/owlcache/tools"
)

type OwlRequest struct {
	//请求命令
	Cmd CommandType
	//key
	Key string
	//请求内容
	Value []byte
	//请求内容长度
	Length int
	//过期时间
	Expires time.Duration //int64
	//链接密码
	Pass string
	//token
	Token string
	//target
	Target string
}

// request to string
func (req *OwlRequest) String() string {
	return fmt.Sprintf("{OwlRequest cmd=%s , key='%s' , value='%v' , length='%d' , expires='%v' , pass='%s' , token='%s' ,bodylen=%d }",
		req.Cmd, req.Key, req.Value, req.Length, req.Expires, req.Pass, req.Token, int64(len(req.Value)))
}

// 过滤接收数据中的\r\n
func (req *OwlRequest) TrimSpace(str string) string {
	if str != "" {
		return strings.TrimSpace(str)
	}
	return ""
}

// 将字符串切片转换成字符串
func (req *OwlRequest) Slicetostring(slice []string) string {

	return strings.Join(slice, " ")

}

// 将socket请求内容 解析为一个OwlRequest对象
func (req *OwlRequest) TCPReceive(connstr string) {

	params := []string{}
	if strings.TrimSpace(connstr) == "ping" {
		params = append(params, "ping")
	} else {
		params = strings.Split(connstr, " ") //strings.Fields(connstr)
		//判断空字符串请求
		if len(params) <= 1 && strings.TrimSpace(params[0]) == "" {
			return
		}
	}

	command := CommandType(params[0])

	switch command {
	case GET:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
		if len(params) > 2 {
			req.Value = []byte(req.TrimSpace(params[2]))
		}
	case EXIST:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
	case SET:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
		if len(params) > 2 {
			req.Value = []byte(req.TrimSpace(req.Slicetostring(params[2:])))
			req.Length = len(req.Value)
		}
	case EXPIRE:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
		if len(params) > 2 {
			exptime, _ := time.ParseDuration(req.TrimSpace(params[2]) + "s")
			req.Expires = exptime
		}
	case DELETE:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
	case PING:
		req.Cmd = command
		if len(params) > 1 {
			req.Value = []byte(req.TrimSpace(req.Slicetostring(params[1:])))
			req.Length = len(req.Value)
		} else {
			req.Value = []byte("")
			req.Length = 0
		}
	}

}

// 将http请求内容 解析为一个OwlRequest对象
func (req *OwlRequest) HTTPReceive(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //解析参数
	//fmt.Println(r.Form)

	//判断是否开启Urlcache的快捷访问
	if owlconfig.OwlConfigModel.Open_Urlcache == "1" && owlconfig.OwlConfigModel.Urlcache_Request_Easy == "1" && len(r.FormValue("key")) < 1 {
		//开启Urlcache的快捷访问后重新定义key值
		req.Key = r.RequestURI
	} else {
		//判断空字符串请求
		if len(r.Form) <= 1 && strings.TrimSpace(r.FormValue("cmd")) == "" {
			return
		}
		req.Key = r.FormValue("key")
	}

	req.Cmd = CommandType(r.FormValue("cmd"))
	req.Value = []byte(r.FormValue("valuedata"))
	req.Length = len(r.FormValue("valuedata"))
	exptime, _ := time.ParseDuration(owltools.JoinString(req.TrimSpace(r.FormValue("exptime")), "s"))
	req.Expires = exptime
	req.Pass = r.FormValue("pass")
	//避免url cache模式开启时与url的token关键字冲突
	if len(r.FormValue("owl_token")) > 0 {
		req.Token = r.FormValue("owl_token")
	} else {
		req.Token = r.FormValue("token")
	}
	req.Target = r.FormValue("target")

	//fmt.Println(req.String())
}

// 将Websocket请求内容 解析为一个OwlRequest对象
func (req *OwlRequest) WebsocketReceive(w http.ResponseWriter, r *http.Request, connstr string) {

	params := []string{}
	params = strings.Split(connstr, " ") //strings.Fields(connstr)
	//判断空字符串请求
	if len(params) <= 1 && strings.TrimSpace(params[0]) == "" {
		return
	}

	command := CommandType(params[0])

	switch command {
	case GET:
		req.Cmd = command
		req.Key = req.TrimSpace(params[1])
		//存储字符串"info"
		if len(params) > 2 {
			req.Value = []byte(req.TrimSpace(params[2]))
		}
		//存储客户端传过来的字符串,建议是唯一标识，例如UUID
		//客户端发送过来的命令格式 get "key_name" info "Handshake_string"
		if len(params) > 3 && len(req.TrimSpace(params[3])) > 0 {
			req.Pass = owltools.JoinString(req.Key, "@", req.TrimSpace(params[3])) //params[3]是客户端传过来的字符串,建议是唯一标识，例如UUID。 handshake=uuid+“_”+address
		}
	}

}
