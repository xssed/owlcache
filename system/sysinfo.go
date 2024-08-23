package system

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/xssed/owlcache/tools"
)

const (
	VERSION      string = "0.5.0-beta"
	VERSION_DATE string = "2024-08-23"
)

//程序启动欢迎信息
func DosSayHello() {

	fmt.Println("Welcome to use owlcache. Version:" + VERSION + "\nIf you have any questions,Please contact us: xsser@xsser.cc \nProject Home:https://github.com/xssed/owlcache")
	fmt.Println(`                _                _          `)
	fmt.Println(`   _____      _| | ___ __ _  ___| |__   ___ `)
	fmt.Println(`  / _ \ \ /\ / / |/ __/ _' |/ __| '_ \ / _ \`)
	fmt.Println(` | (_) \ V  V /| | (_| (_| | (__| | | |  __/`)
	fmt.Println(`  \___/ \_/\_/ |_|\___\__,_|\___|_| |_|\___|`)
	fmt.Println(`                                            `)

}

//http服务欢迎页
func HttpSayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<style type='text/css'>*{ padding: 0; margin: 0; } div{ padding: 4px 48px;} a{color:#2E5CD5;cursor: pointer;text-decoration: none} a:hover{text-decoration:underline; } body{ background: #fff; font-family: 'Century Gothic','Microsoft yahei'; color: #333;font-size:18px;} h1{ font-size: 100px; font-weight: normal; margin-bottom: 12px; } p{ line-height: 1.6em; font-size: 42px }</style><div style='padding: 24px 48px;'><h1>:)</h1><p>Welcome to use owlcache. Version:"+VERSION+"<br/><span style='font-size:25px'>If you have any questions,Please contact us: <a href=\"mailto:xsser@xsser.cc\">xsser@xsser.cc</a><br>Project Home : <a href=\"https://github.com/xssed/owlcache\" target=\"_blank\">https://github.com/xssed/owlcache</a></span></p><div>")
}

//输出Owl的内存信息
func MemStats() string {

	var m runtime.MemStats

	// Sys 服务现在系统使用的内存
	// NumGC 垃圾回收调用次数
	// Alloc 堆空间分配的字节数
	// TotalAlloc 从服务开始运行至今分配器为分配的堆空间总和，只有增加，释放的时候不减少。

	runtime.ReadMemStats(&m)

	unit_gb := 1024 * 1024 * 1024.0
	alloc := float64(m.Alloc) / unit_gb
	totalalloc := float64(m.TotalAlloc) / unit_gb
	sys := float64(m.Sys) / unit_gb
	numgc := m.NumGC

	logstr := fmt.Sprintf("Sys = %vGB  NumGC = %v  Alloc = %vGB  TotalAlloc = %vGB  ", tools.RoundedFixed(sys, 7), numgc, tools.RoundedFixed(alloc, 7), tools.RoundedFixed(totalalloc, 7))
	return logstr

}
