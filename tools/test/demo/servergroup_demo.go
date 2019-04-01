package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xssed/owlcache/group"
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

var servergroup *group.Servergroup

func main() {

	servergroup = group.NewServergroup()
	go func() {
		servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.1", "1111111111111", "111"})
	}()
	go func() {
		servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.2", "1111111111111", "222"})
	}()
	go func() {
		servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.3", "1111111111111", "333"})
	}()
	go func() {
		servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.4", "1111111111111", "444"})
	}()
	go func() {
		servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.5", "1111111111111", "555"})
	}()
	go func() {
		servergroup.Add(OwlServerGroupRequest{"", "http://192.168.0.6", "1111111111111", "666"})
	}()
	go func() {
		servergroup.AddAt(2, OwlServerGroupRequest{"", "http://192.168.0.7", "1111111111111", ""})
	}()

	for i := 0; i < 100; i++ {
		go servergroup.AddAt(2, OwlServerGroupRequest{"", "http://192.168.10." + strconv.Itoa(i), strconv.Itoa(i), ""})
	}

	time.Sleep(time.Second * 10)

	//	fmt.Println(servergroup.AddAt(2, OwlServerGroupRequest{"", "http://192.168.0.8", "1111111111111", ""}))
	//	fmt.Println(servergroup.AddAt(222, OwlServerGroupRequest{"", "http://192.168.0.8", "1111111111111", ""}))
	//	fmt.Println(servergroup.ToSliceString())
	//	servergroup.RemoveFirst()
	//	fmt.Println(servergroup.ToSliceString())
	//	servergroup.RemoveLast()
	//	fmt.Println(servergroup.ToSliceString())
	//	servergroup.RemoveAt(1)
	//	fmt.Println(servergroup.ToSliceString())
	//	servergroup.Clear()
	//	fmt.Println(servergroup.ToSliceString())
	//	fmt.Println(servergroup.GetFirst())
	//	fmt.Println(servergroup.GetLast())
	//	fmt.Println(servergroup.GetAt(2))
	//	fmt.Println(servergroup.GetAt(2222))
	//	fmt.Println(servergroup.GetRange(2, 3))
	//	fmt.Println(servergroup.Count())
	resat := 0
	resbool := false
	fmt.Println(servergroup.Count())

	list := servergroup.Values()
	fmt.Println(list)
	for k := range list {
		val, ok := list[k].(OwlServerGroupRequest)
		if ok {
			if val.Address == "http://192.168.0.2" {
				resat = k
				resbool = true
			}
		}
	}
	fmt.Println(resat)   //位置
	fmt.Println(resbool) //存在与否

	fmt.Println(servergroup.Values())

	servergroup.SaveToFile("./db_file/servergroup.db")
	servergroup.Clear() //Clear()清空
	fmt.Println(servergroup.Values())
	servergroup.LoadFromFile("./db_file/servergroup.db")
	fmt.Println(servergroup.Values())
	//fmt.Println(servergroup.Exists("192.168.1.10"))
	fmt.Println(servergroup.Count())

}
