package network

import (
	"encoding/json"
	"fmt"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
	"github.com/xssed/owlcache/group"
	"github.com/xssed/owlcache/network/gossip"

	//"github.com/xssed/owlcache/tools"
	owllog "github.com/xssed/owlcache/log"
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
				str_addresslist = append(str_addresslist, val.Address)
			}
		}
	}

	bindAddress := owlconfig.OwlConfigModel.Host               //host
	bindPort := owlconfig.OwlConfigModel.Gossipport            //gossip端口
	passWord := owlconfig.OwlConfigModel.GossipDataSyncAuthKey //交互密码

	if err := gossip.H.StartService(str_addresslist, passWord, bindAddress, bindPort); err != nil {
		fmt.Println(err)
	}

	go listenGossipQueue()

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
						fmt.Println(err)
					}
					//fmt.Println("json to map ", result)
				}

				switch result["cmd"] {
				case "set":
					exptime, _ := time.ParseDuration(result["expire"] + "s")
					ok := BaseCacheDB.Set(result["key"], []byte(result["val"]), exptime)
					if !ok {
						owllog.OwlLogHttp.Println("gossip:set error " + " key:" + result["key"])
					}
				case "expire":
					exptime, _ := time.ParseDuration(result["expire"] + "s")
					ok := BaseCacheDB.Expire(result["key"], exptime)
					if !ok {
						owllog.OwlLogHttp.Println("gossip:expire error " + " key:" + result["key"])
					}
				case "del":
					ok := BaseCacheDB.Delete(result["key"])
					if !ok {
						owllog.OwlLogHttp.Println("gossip:del error " + " key:" + result["key"])
					}
				}

			}
		}

	}

}
