package network

import (
	//"encoding/json"
	//"fmt"
	//"log"

	//"net/http"

	"github.com/xssed/owlcache/group"
	//tools "github.com/xssed/owlcache/tools"
)

//发起请求获取集合数据
func (owlhandler *OwlHandler) GetGroupData() { //r *http.Request

	list := ServerGroupList.Values()
	//fmt.Println(list)

	//ch := make(chan string)
	//count := ServerGroupList.Count() //count 表示活动的协程个数
	//fmt.Println(count)
	for k := range list {
		val, ok := list[k].(group.OwlServerGroupRequest)
		if ok {
			//fmt.Println(val)
			go owlhandler.ParseContent(val.Address, owlhandler.owlrequest.Key)

		}
	}

	owlhandler.Transmit(SUCCESS)
	owlhandler.owlresponse.Data = "123"

}

//解析内容
func (owlhandler *OwlHandler) ParseContent(address, key string) {

	s := HttpClient.GetValue(address, key)
	if s != "" {

		// var resbody *OwlResponse
		// if err := json.Unmarshal([]byte(s), &resbody); err != nil {
		// 	log.Fatalf("OwlHandler ParseContent JSON unmarshling failed: %s", err)
		// }
		// ch <- resbody
		//owlhandler.owlhcp
		//fmt.Println(resbody.KeyCreateTime)
	}

}
