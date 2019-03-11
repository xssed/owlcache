package main

import (
	"fmt"

	"github.com/xssed/owlcache/owlgroup"
)

type OwlServerGroupRequest struct {
	//请求命令
	Cmd string
	//地址字符串
	Address string
	//链接密码
	Pass string
	//token
	Token string
}

var servergroup *owlgroup.Servergroup

func main() {

	servergroup = owlgroup.NewServergroup()
	servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.1", "1111111111111", ""})
	servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.2", "1111111111111", ""})
	servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.3", "1111111111111", ""})
	servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.4", "1111111111111", ""})
	servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.5", "1111111111111", ""})
	servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.6", "1111111111111", ""})
	servergroup.AddAt(2, OwlServerGroupRequest{"", "http://192.168.0.7", "1111111111111", ""})
	servergroup.AddAt(2, OwlServerGroupRequest{"", "http://192.168.0.8", "1111111111111", ""})
	//	fmt.Println(servergroup.ToSliceString())
	//	servergroup.RemoveFirst()
	//	fmt.Println(servergroup.ToSliceString())
	//	servergroup.RemoveLast()
	//	fmt.Println(servergroup.ToSliceString())
	//	servergroup.RemoveAt(1)
	//	fmt.Println(servergroup.ToSliceString())
	//	servergroup.Clear()
	fmt.Println(servergroup.ToSliceString())
	fmt.Println(servergroup.GetFirst())
	fmt.Println(servergroup.GetLast())
	fmt.Println(servergroup.GetAt(2))
	fmt.Println(servergroup.GetRange(2, 3))
	fmt.Println(servergroup.Count())
	fmt.Println(servergroup.Values())

	//fmt.Println(servergroup.Exists("192.168.1.10"))

}
