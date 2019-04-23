package network

import (
	//"fmt"
	"net/http"

	owlconfig "github.com/xssed/owlcache/config"
	//owllog "github.com/xssed/owlcache/log"
	owlsystem "github.com/xssed/owlcache/system"
)

type ServerEntity struct {
	handler http.Handler
}

func (se *ServerEntity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//服务器信息，可以在客户端校验版本
	w.Header().Set("Server", "owlcache "+owlsystem.VERSION)
	//Cors值为"1"(开启服务)和"0"(关闭服务)。默认为0关闭服务不允许跨域。
	if owlconfig.OwlConfigModel.Cors == "1" {
		w.Header().Set("Access-Control-Allow-Origin", owlconfig.OwlConfigModel.Access_Control_Allow_Origin)
	}
	// else if owlconfig.OwlConfigModel.Cors == "0" {
	// 	owllog.OwlLogRun.Info("The configuration file Cors value is '0'. 'Access-Control-Allow-Origin' not set.")
	// } else {
	// 	//检测到配置书写异常强制退出
	// 	owllog.OwlLogRun.Fatal(ErrorCors)
	// }
	//继续传递信息
	se.handler.ServeHTTP(w, r)
}
