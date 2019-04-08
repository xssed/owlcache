package network

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/xssed/owlcache/group"
	//tools "github.com/xssed/owlcache/tools"
)

//发起请求获取集合数据
func (owlhandler *OwlHandler) GetGroupData() {

	list := ServerGroupList.Values()
	//fmt.Println(list)

	//count := ServerGroupList.Count() //count 表示活动的协程个数
	//fmt.Println("count:", count)

	//服务器集群信息存储列表
	groupKVlist := group.NewServergroup()

	var wg sync.WaitGroup

	for k := range list {
		val, ok := list[k].(group.OwlServerGroupRequest)
		if ok {
			//fmt.Println(val)
			wg.Add(1)
			go owlhandler.ParseContent(val.Address, owlhandler.owlrequest.Key, groupKVlist, &wg)
		}
	}
	wg.Wait()

	//fmt.Println(groupKVlist.Values())

	fmt.Println(owlhandler.BubbleSortContent(groupKVlist))

	owlhandler.Transmit(SUCCESS)
	owlhandler.owlresponse.Data = "123"

}

//解析内容
func (owlhandler *OwlHandler) ParseContent(address, key string, kvlist *group.Servergroup, wg *sync.WaitGroup) {

	defer wg.Done()

	s := HttpClient.GetValue(address, key)
	if s != "" {
		var resbody OwlResponse
		if err := json.Unmarshal([]byte(s), &resbody); err != nil {
			log.Fatalf("OwlHandler ParseContent JSON unmarshling failed: %s", err)
		}
		kvlist.Add(resbody)
		//kvlist.Add(s)
		//fmt.Println(resbody)
	}

}

//排序
func (owlhandler *OwlHandler) BubbleSortContent(kvlist *group.Servergroup) []OwlResponse {

	var array []OwlResponse

	list := kvlist.Values()
	for k := range list {
		val, ok := list[k].(OwlResponse)
		if ok {
			array = append(array, val)
		}
	}

	var sorted = false
	for !sorted {
		sorted = true
		for i := 0; i < len(array)-1; i++ {
			if array[i].KeyCreateTime.Unix() < array[i+1].KeyCreateTime.Unix() {
				sorted = false
				array[i], array[i+1] = array[i+1], array[i]
			}
		}
		//fmt.Println(array)
	}

	return array

}
