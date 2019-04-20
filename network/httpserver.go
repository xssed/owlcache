package network

import (
	"fmt"
	"net/http"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owlsystem "github.com/xssed/owlcache/system"
)

func stratHTTP() {

	//设置监听的端口
	addr := owlconfig.OwlConfigModel.Host + ":" + owlconfig.OwlConfigModel.Httpport
	//默认信息
	http.HandleFunc("/", IndexPage) //设置访问的路由
	//单机数据执行信息
	http.HandleFunc("/data/", Exe) //设置访问的路由
	//群组数据执行信息
	http.HandleFunc("/group_data/", GroupExe) //设置访问的路由
	//设置服务器集群
	http.HandleFunc("/server/", Server) //设置服务器集群信息，单机。
	//	//设置服务器集群
	//	http.HandleFunc("/server_group/", ServerGroup) //设置服务器集群信息,集群。

	//监听设置
	var err error
	if owlconfig.OwlConfigModel.Open_Https == "1" {
		//支持HTTPS
		err = http.ListenAndServeTLS(addr, owlconfig.OwlConfigModel.Https_CertFile, owlconfig.OwlConfigModel.Https_KeyFile, &ServerEntity{handler: http.DefaultServeMux})
	} else if owlconfig.OwlConfigModel.Open_Https == "0" {
		//普通HTTP
		err = http.ListenAndServe(addr, &ServerEntity{handler: http.DefaultServeMux})
	} else {
		err = ErrorOpenHttpsSelected
	}
	if err != nil {
		owllog.OwlLogRun.Fatal("Error starting HTTP server.", err)
	} else {
		owllog.OwlLogRun.Info("Creating HTTP server with address " + addr)
	}

}

//单机数据执行信息
func Exe(w http.ResponseWriter, r *http.Request) {

	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.HTTPReceive(w, r)
	owlhandler.HTTPHandle(w, r) //执行数据
	resstr := owlhandler.owlresponse.ConvertToString("HTTP")
	fmt.Fprintf(w, resstr) //输出到客户端的信息

}

//群组数据执行信息
func GroupExe(w http.ResponseWriter, r *http.Request) {

	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.HTTPReceive(w, r)
	owlhandler.HTTPGroupDataHandle(w, r) //执行数据
	resstr := owlhandler.owlresponse.ConvertToString("HTTP")
	fmt.Fprintf(w, resstr) //输出到客户端的信息

}

//设置服务器集群信息，单机
func Server(w http.ResponseWriter, r *http.Request) {

	owlservergrouphandler := NewOwlServerGroupHandler()
	owlservergrouphandler.owlservergrouprequest.HTTPReceive(w, r)
	owlservergrouphandler.HTTPServerHandle(w, r) //执行数据
	resstr := owlservergrouphandler.owlserveggroupresponse.ConvertToString()
	fmt.Fprintf(w, resstr) //输出到客户端的信息

}

////设置服务器集群信息，集群
//func ServerGroup(w http.ResponseWriter, r *http.Request) {

//	owlservergrouphandler := NewOwlServerGroupHandler()
//	owlservergrouphandler.owlservergrouprequest.HTTPReceive(w, r)
//	owlservergrouphandler.HTTPServerGroupHandle(w, r) //执行数据
//	resstr := owlservergrouphandler.owlserveggroupresponse.ConvertToString()
//	fmt.Fprintf(w, resstr) //输出到客户端的信息

//}

//首页
func IndexPage(w http.ResponseWriter, r *http.Request) {
	owlsystem.HttpSayHello(w, r)
}
