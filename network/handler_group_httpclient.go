package network

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/xssed/owlcache/group"
	owllog "github.com/xssed/owlcache/log"
)

//发起请求获取集合数据
func (owlhandler *OwlHandler) GetGroupData(w http.ResponseWriter, r *http.Request) {

	//如果valuedata不是字符串info则输出集群第一个
	if string(owlhandler.owlrequest.Value) != "info" {
		//默认是同时去请求集群数据，所有数据请求完毕(超时则中止请求)，再返回数据
		resmap := owlhandler.conversionContent(owlhandler.getHttpData())
		//如果没有在集群中取到值
		if len(resmap) == 0 {
			owlhandler.Transmit(NOT_FOUND)
			owlhandler.owlresponse.Data = []byte("")
			return
		}
		//取出集群第一个最新的数据
		owlhandler.Transmit(SUCCESS)
		owlhandler.owlresponse.Data = resmap[0]["Data"].([]byte)
		return

	}
	//查询info类型
	resmap := owlhandler.conversionContentInfo(owlhandler.getHttpData())
	//如果没有在集群中取到值
	if len(resmap) == 0 {
		owlhandler.Transmit(NOT_FOUND)
		owlhandler.owlresponse.Data = []byte("")
		return
	}
	owlhandler.Transmit(SUCCESS)
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	data, _ := json.Marshal(&resmap)
	owlhandler.owlresponse.Data = data
	return

}

//发起请求获取数据
func (owlhandler *OwlHandler) getHttpData() []OwlResponse {

	list := ServerGroupList.Values()
	//fmt.Println(list)

	//count := ServerGroupList.Count()
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
	if s != nil {

		var resbody OwlResponse

		resbody.Status = ResStatus(s.StatusCode)
		resbody.Key = s.Header.Get("Key")
		resbody.Data = s.Byte()
		//时间处理部分
		tkt := s.Header.Get("Keycreatetime")
		if len(tkt) > 37 {
			tkt = tkt[0:37]
		}
		t, terr := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", tkt)
		if terr != nil {
			owllog.OwlLogHttp.Info("OwlHandler parseContent Keycreatetime time.Parse failed: " + terr.Error())
		}
		resbody.KeyCreateTime = t
		resbody.ResponseHost = s.Header.Get("Responsehost")
		kvlist.Add(resbody)

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
	}

	return array

}

//封装返回key集合信息
func (owlhandler *OwlHandler) conversionContentInfo(res_slice []OwlResponse) []map[string]interface{} {

	var response_list []map[string]interface{}

	for index := range res_slice {

		oldresponse := res_slice[index]
		//fmt.Println(oldresponse)
		//只接受存在的数据（响应内容状态为200）
		if oldresponse.Status == 200 {
			temp_map := make(map[string]interface{})
			temp_map["Address"] = oldresponse.ResponseHost
			temp_map["Key"] = oldresponse.Key
			temp_map["Status"] = oldresponse.Status
			temp_map["KeyCreateTime"] = oldresponse.KeyCreateTime
			response_list = append(response_list, temp_map)
		}

	}

	return response_list

}

//封装返回key集合数据
func (owlhandler *OwlHandler) conversionContent(res_slice []OwlResponse) []map[string]interface{} {

	var response_list []map[string]interface{}

	for index := range res_slice {

		oldresponse := res_slice[index]
		//只接受存在的数据（响应内容状态为200）
		if oldresponse.Status == 200 {
			temp_map := make(map[string]interface{})
			temp_map["Data"] = oldresponse.Data
			response_list = append(response_list, temp_map)
		}

	}

	return response_list

}
