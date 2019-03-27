package network

import (
	"fmt"
	"log"
	"net/http"

	owlconfig "github.com/xssed/owlcache/config"
	owlsystem "github.com/xssed/owlcache/system"
)

func stratHTTP() {

	addr := owlconfig.OwlConfigModel.Host + ":" + owlconfig.OwlConfigModel.Httpport
	//默认信息
	http.HandleFunc("/", IndexPage) //设置访问的路由
	//单机数据执行信息
	http.HandleFunc("/data/", Exe) //设置访问的路由
	//群组数据执行信息
	http.HandleFunc("/group_data/", GroupExe) //设置访问的路由
	//设置服务器集群
	http.HandleFunc("/server/", Server) //设置服务器集群信息，单机。
	//设置服务器集群
	http.HandleFunc("/server_group/", ServerGroup) //设置服务器集群信息，单机。

	//监听设置
	err := http.ListenAndServe(addr, nil) //设置监听的端口
	if err != nil {
		fmt.Println("Error starting HTTP server.")
		log.Fatal("ListenAndServe: ", err)
	}
}

//单机数据执行信息
func Exe(w http.ResponseWriter, r *http.Request) {

	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.HTTPReceive(w, r)
	owlhandler.HTTPHandle(w, r) //执行数据
	resstr := owlhandler.owlresponse.ConvertToString()
	fmt.Fprintf(w, resstr) //输出到客户端的信息

}

//群组数据执行信息
func GroupExe(w http.ResponseWriter, r *http.Request) {

	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.HTTPReceive(w, r)
	owlhandler.HTTPGroupDataHandle(w, r) //执行数据
	resstr := owlhandler.owlresponse.ConvertToString()
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

//设置服务器集群信息，集群
func ServerGroup(w http.ResponseWriter, r *http.Request) {

	owlservergrouphandler := NewOwlServerGroupHandler()
	owlservergrouphandler.owlservergrouprequest.HTTPReceive(w, r)
	owlservergrouphandler.HTTPServerGroupHandle(w, r) //执行数据
	resstr := owlservergrouphandler.owlserveggroupresponse.ConvertToString()
	fmt.Fprintf(w, resstr) //输出到客户端的信息

}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<style type='text/css'>*{ padding: 0; margin: 0; } div{ padding: 4px 48px;} a{color:#2E5CD5;cursor: pointer;text-decoration: none} a:hover{text-decoration:underline; } body{ background: #fff; font-family: 'Century Gothic','Microsoft yahei'; color: #333;font-size:18px;} h1{ font-size: 100px; font-weight: normal; margin-bottom: 12px; } p{ line-height: 1.6em; font-size: 42px }</style><div style='padding: 24px 48px;'><h1>:)</h1><p>Welcome to use Owlcache. Version:"+owlsystem.VERSION+"<br/><span style='font-size:25px'>If you have any questions,Please contact us: <a href=\"mailto:xsser@xsser.cc\">xsser@xsser.cc</a><br>Project Home : <a href=\"https://github.com/xssed/owlcache\" target=\"_blank\">https://github.com/xssed/owlcache</a></span></p><div>")
}
