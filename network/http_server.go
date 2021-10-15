package network

import (
	"fmt"
	"net/http"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owlsystem "github.com/xssed/owlcache/system"
)

func startHTTP() {

	//设置监听的端口
	addr := owlconfig.OwlConfigModel.Host + ":" + owlconfig.OwlConfigModel.Httpport
	//默认信息
	http.HandleFunc("/", IndexPage) //设置访问的路由
	//单机数据执行信息
	http.HandleFunc("/data/", Exe) //设置访问的路由
	//设置服务器集群
	http.HandleFunc("/server/", Server) //设置服务器集群信息，单机。
	//群组数据执行信息
	http.HandleFunc("/group_data/", GroupExe) //设置访问的路由
	//启动gossip数据最终一致服务。检查是否开启gossip服务。默认为关闭。
	// if owlconfig.OwlConfigModel.GroupWorkMode == "gossip" {
	// 	//什么也不做
	// } else if owlconfig.OwlConfigModel.GroupWorkMode == "owlcache" {
	// 	//群组数据执行信息
	// 	http.HandleFunc("/group_data/", GroupExe) //设置访问的路由
	// } else {
	// 	//检测到配置书写异常强制退出
	// 	owllog.OwlLogRun.Fatal(ErrorGroupWorkMode)
	// }

	//监听设置
	var err error
	if owlconfig.OwlConfigModel.Open_Https == "1" {
		//支持HTTPS
		owllog.OwlLogRun.Info("Creating HTTPS server with address:" + addr)
		err = http.ListenAndServeTLS(addr, owlconfig.OwlConfigModel.Https_CertFile, owlconfig.OwlConfigModel.Https_KeyFile, &ServerEntity{handler: http.DefaultServeMux})
	} else if owlconfig.OwlConfigModel.Open_Https == "0" {
		//普通HTTP
		owllog.OwlLogRun.Info("Creating HTTP server with address:" + addr)
		err = http.ListenAndServe(addr, &ServerEntity{handler: http.DefaultServeMux})
	} else {
		err = ErrorOpenHttpsSelected
	}
	if err != nil {
		owllog.OwlLogRun.Fatal("Error starting HTTP server.", err)
	}

}

//单机数据执行信息
func Exe(w http.ResponseWriter, r *http.Request) {

	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.HTTPReceive(w, r)
	owlhandler.HTTPHandle(w, r) //执行数据
	resstr := owlhandler.owlresponse.ConvertToString("HTTP")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, resstr) //输出到客户端的信息

}

//群组数据执行信息
func GroupExe(w http.ResponseWriter, r *http.Request) {

	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.HTTPReceive(w, r)
	owlhandler.HTTPGroupDataHandle(w, r) //执行数据
	resstr := owlhandler.owlresponse.ConvertToString("HTTP")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, resstr) //输出到客户端的信息

}

//设置服务器集群信息，默认集群与gossip
func Server(w http.ResponseWriter, r *http.Request) {

	owlservergrouphandler := NewOwlServerGroupHandler()
	owlservergrouphandler.owlservergrouprequest.HTTPReceive(w, r)

	if r.FormValue("group_type") == "gossip" {
		owlservergrouphandler.HTTPServerGroupHandle(w, r) //执行数据
	} else {
		owlservergrouphandler.HTTPServerHandle(w, r) //执行数据
	}

	resstr := owlservergrouphandler.owlserveggroupresponse.ConvertToString()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, resstr) //输出到客户端的信息

}

//首页
func IndexPage(w http.ResponseWriter, r *http.Request) {
	owlsystem.HttpSayHello(w, r)
}
