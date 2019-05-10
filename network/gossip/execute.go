package gossip

import (
	"encoding/json"
)

//数据包结构(数据交换)
type Execute map[string]string

func PostQueueBroadcast(exedata *Execute) {
	var exedata_list []*Execute
	exedata_list = append(exedata_list, exedata)
	b, err := json.Marshal(exedata_list)
	if err != nil {
		return
	}
	//发送数据到集群
	H.QueueBroadcast(b)
}

//设置Key数据
func Set(key, val string) {
	exedata := make(Execute)
	exedata["cmd"] = "add"
	exedata["key"] = key
	exedata["val"] = val
	PostQueueBroadcast(&exedata)
}

//为Key设置过期时间
func Delete(key string) {
	exedata := make(Execute)
	exedata["cmd"] = "del"
	exedata["key"] = key
	exedata["val"] = ""
	PostQueueBroadcast(&exedata)
}

//为Key设置过期时间
func Expire(key, expires string) {
	exedata := make(Execute)
	exedata["cmd"] = "expire"
	exedata["key"] = key
	exedata["val"] = expires
	PostQueueBroadcast(&exedata)
}
