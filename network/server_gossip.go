package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
	"github.com/xssed/owlcache/group"
	"github.com/xssed/owlcache/network/gossip"

	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
)

func startGossip() {

	var str_addresslist []string
	//list := ServerGroupList.Values()
	list := ServerGroupGossipList.Values() //Gossip集群信息与默认集群方式分离
	for k := range list {
		//fmt.Println(tools.Typeof(list[k]))
		val, ok := list[k].(group.OwlServerGroupRequest)
		if ok {
			if val.Address != "" {
				//加载初始化时校验配置文件中gossip服务地址能否连接
				conn, err := net.DialTimeout("tcp", val.Address, 3*time.Second)
				if err != nil {
					owllog.OwlLogGossip.Info(owltools.JoinString("gossip:start error ", err.Error()))
					os.Exit(0)
				} else {
					if conn != nil {
						_ = conn.Close()
					}
				}
				//将节点添加进入集群信息
				str_addresslist = append(str_addresslist, val.Address)
			}
		}
	}

	bindAddress := owlconfig.OwlConfigModel.Host               //host
	bindPort := owlconfig.OwlConfigModel.Gossipport            //gossip端口
	passWord := owlconfig.OwlConfigModel.GossipDataSyncAuthKey //交互密码

	if err := gossip.H.StartService(str_addresslist, passWord, bindAddress, bindPort); err != nil {
		fmt.Println(err.Error())
	}

	//监听队列
	listenGossipQueue()

}

func listenGossipQueue() {

	for {

		time.Sleep(time.Microsecond * 7) //微秒级阻塞

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
						fmt.Println(err.Error())
					}
					//fmt.Println("json to map ", result)
				}

				switch result["cmd"] {
				case "set":
					go gossip_set(result["key"], result["val"], result["expire"])
				case "expire":
					go gossip_expire(result["key"], result["expire"])
				case "del":
					go gossip_del(result["key"])
				}

			}
		}

	}

}
