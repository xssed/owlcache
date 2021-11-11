package network

import (
	//"fmt"
	"net/http"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owlsystem "github.com/xssed/owlcache/system"
)

func startHTTP() {

	//设置监听的端口
	addr := owlconfig.OwlConfigModel.Host + ":" + owlconfig.OwlConfigModel.Httpport

	//单机数据执行信息
	http.HandleFunc("/data/", Exe) //设置访问的路由
	//群组数据执行信息
	http.HandleFunc("/group_data/", GroupExe) //设置访问的路由
	//判断是否开启Url Cache
	if owlconfig.OwlConfigModel.Open_Urlcache == "1" {
		//Url Cache数据执行信息
		http.HandleFunc("/uc/", UCExe) //设置访问的路由
	}
	//设置服务器集群
	http.HandleFunc("/server/", Server) //设置服务器集群信息，单机。
	//判断是否开启Urlcache的快捷访问
	if owlconfig.OwlConfigModel.Open_Urlcache == "1" && owlconfig.OwlConfigModel.Urlcache_Request_Easy == "1" {
		//Url Cache数据执行信息
		http.HandleFunc("/", UCExe) //设置访问的路由
	} else {
		//默认信息
		http.HandleFunc("/", IndexPage) //设置访问的路由
	}

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
	var print []byte
	w, print = owlhandler.ToHttp(w)
	//设置响应状态
	w.WriteHeader(int(owlhandler.owlresponse.Status))
	w.Write(print) //输出到客户端的信息

}

//群组数据执行信息
func GroupExe(w http.ResponseWriter, r *http.Request) {

	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.HTTPReceive(w, r)
	owlhandler.HTTPGroupDataHandle(w, r) //执行数据
	var print []byte
	w, print = owlhandler.ToGroupHttp(w, r) //数据转换
	//设置响应状态
	w.WriteHeader(int(owlhandler.owlresponse.Status))
	w.Write(print) //输出到客户端的信息

}

//UrlCache数据执行信息
func UCExe(w http.ResponseWriter, r *http.Request) {

	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.HTTPReceive(w, r)
	var print []byte
	w, print = owlhandler.UCDataHandle(w, r) //执行数据
	w.Write(print)                           //输出到客户端的信息

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
	var print []byte
	w, print = owlservergrouphandler.HttpGroupGetKeyInfoToString(w)
	w.Write(print)

}

//首页
func IndexPage(w http.ResponseWriter, r *http.Request) {
	owlsystem.HttpSayHello(w, r)
}
