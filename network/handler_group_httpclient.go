package network

import (
	"encoding/json"
	"sync"

	"github.com/xssed/owlcache/group"
	owllog "github.com/xssed/owlcache/log"
)

//发起请求获取集合数据
func (owlhandler *OwlHandler) GetGroupData() {

	owlhandler.Transmit(SUCCESS)
	owlhandler.owlresponse.Data = owlhandler.conversionContent(owlhandler.getHttpData())

}

//发起请求获取数据
func (owlhandler *OwlHandler) getHttpData() []OwlResponse {

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
			go owlhandler.parseContent(val.Address, owlhandler.owlrequest.Key, groupKVlist, &wg)
		}
	}
	wg.Wait()

	//fmt.Println(groupKVlist.Values())
	//排序数据
	bubblesortlist := owlhandler.bubbleSortContent(groupKVlist)
	//fmt.Println(bubblesortlist)

	return bubblesortlist

}

//解析内容
func (owlhandler *OwlHandler) parseContent(address, key string, kvlist *group.Servergroup, wg *sync.WaitGroup) {

	defer wg.Done()

	s := HttpClient.GetValue(address, key)
	if s != "" {
		var resbody OwlResponse
		if err := json.Unmarshal([]byte(s), &resbody); err != nil {
			owllog.OwlLogHttp.Fatalf("OwlHandler parseContent JSON unmarshling failed: %s", err)
		}
		kvlist.Add(resbody)
		//kvlist.Add(s)
		//fmt.Println(resbody)
	}

}

//排序
func (owlhandler *OwlHandler) bubbleSortContent(kvlist *group.Servergroup) []OwlResponse {

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

//封装返回数据
func (owlhandler *OwlHandler) conversionContent(res_slice []OwlResponse) []byte {

	var response_list []map[string]interface{}

	for index := range res_slice {

		oldresponse := res_slice[index]
		//只接受存在的数据（响应内容状态为200）
		if oldresponse.Status == 200 {
			temp_map := make(map[string]interface{})
			temp_map["Address"] = oldresponse.ResponseHost
			temp_map["Status"] = oldresponse.Status
			temp_map["Data"] = oldresponse.Data
			temp_map["KeyCreateTime"] = oldresponse.KeyCreateTime
			response_list = append(response_list, temp_map)
		}

	}

	data, err := json.Marshal(response_list)
	if err != nil {
		data = []byte("")
	}
	return data

}
