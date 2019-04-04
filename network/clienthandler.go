package network

import (
	//"encoding/json"
	"fmt"
	//"log"

	//"net/http"

	"github.com/xssed/owlcache/group"
	//tools "github.com/xssed/owlcache/tools"
)

//发起请求获取集合数据
func (owlhandler *OwlHandler) GetGroupData() { //r *http.Request

	list := ServerGroupList.Values()
	//fmt.Println(list)

	ch := make(chan string)
	count := ServerGroupList.Count() //count 表示活动的协程个数
	fmt.Println(count)
	for k := range list {
		val, ok := list[k].(group.OwlServerGroupRequest)
		if ok {
			//fmt.Println(val)
			go owlhandler.ParseContent(val.Address, owlhandler.owlrequest.Key, ch)

		}
	}

	for range ch {
		// 每次从ch中接收数据，表明一个活动的协程结束
		count--

		fmt.Println(ch)

		// 当所有活动的协程都结束时，关闭管道
		if count == 0 {
			close(ch)
			continue
		}
	}

	//	owlservergrouphandler.Transmit(group.NOT_FOUND)

}

//解析内容
func (owlhandler *OwlHandler) ParseContent(address, key string, ch chan string) {

	s := HttpClient.GetValue(address, key)
	if s != "" {

		// var resbody *OwlResponse
		// if err := json.Unmarshal([]byte(s), &resbody); err != nil {
		// 	log.Fatalf("OwlHandler ParseContent JSON unmarshling failed: %s", err)
		// }
		// ch <- resbody
		ch <- s
		//owlhandler.owlhcp
		//fmt.Println(resbody.KeyCreateTime)
	}

}
