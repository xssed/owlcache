package network

import (
	"encoding/json"
	"fmt"

	"github.com/xssed/owlcache/network/gossip"
	//owlconfig "github.com/xssed/owlcache/config"
)

func startGossip() {

	//list := ServerGroupList.Values()
	fmt.Println("startGossip()")

	//addr := owlconfig.OwlConfigModel.Host + ":" + owlconfig.OwlConfigModel.Tcpport

	if err := gossip.H.StartService(); err != nil {
		fmt.Println(err)
	}

	go func() {

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
					case "add":
						//data.Set(result["key"], result["val"])
						fmt.Println("add")
					case "del":
						//data.Delete(result["key"])
						fmt.Println("del")
					}

				}
			}

		}
	}()

}
