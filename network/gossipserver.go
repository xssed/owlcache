package network

import (
	"encoding/json"
	"fmt"

	owlconfig "github.com/xssed/owlcache/config"
	"github.com/xssed/owlcache/group"
	"github.com/xssed/owlcache/network/gossip"
	//"github.com/xssed/owlcache/tools"
)

func startGossip() {

	var str_addresslist []string
	list := ServerGroupList.Values()
	for k := range list {
		//fmt.Println(tools.Typeof(list[k]))
		val, ok := list[k].(group.OwlServerGroupRequest)
		if ok {
			str_addresslist = append(str_addresslist, val.Address)
		}
	}

	bindAddress := owlconfig.OwlConfigModel.Host    //host
	bindPort := owlconfig.OwlConfigModel.Gossipport //gossip端口
	passWord := owlconfig.OwlConfigModel.Pass       //交互密码

	if err := gossip.H.StartService(str_addresslist, passWord, bindAddress, bindPort); err != nil {
		fmt.Println(err)
	}

	go listenGossipQueue()

}

func listenGossipQueue() {

	for {

		size := gossip.Q.Size()
		if size >= 1 {
			e := gossip.Q.Pop()
			//fmt.Println("结果:", e)
			if e != nil {

				var result gossip.Execute
				v, convert_ok := e.(string)
				if convert_ok {
					//fmt.Println("string:", v)
					if err := json.Unmarshal([]byte(v), &result); err != nil {
						fmt.Println(err)
					}
					//fmt.Println("json to map ", result)
				}

				switch result["cmd"] {
				case "set":
					fmt.Println("set")
					fmt.Println(result)
				case "expire":
					fmt.Println("expire")
					fmt.Println(result)
				case "del":
					fmt.Println("del")
					fmt.Println(result)
				}

			}
		}

	}

}
