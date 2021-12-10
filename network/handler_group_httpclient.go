package network

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/xssed/owlcache/cache"
	owlconfig "github.com/xssed/owlcache/config"
	"github.com/xssed/owlcache/group"
	owllog "github.com/xssed/owlcache/log"
	"github.com/xssed/owlcache/network/httpclient"
	owltools "github.com/xssed/owlcache/tools"
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

	//HttpClient缓存优化
	bhgc_key := owltools.JoinString(address, ":", key)
	bhgc_key_responsehost := owltools.JoinString(address, ":", key, ":", "Responsehost")
	bhgc_key_state := owltools.JoinString(address, ":", key, ":", "State")
	//查询二级缓存中是否有数据
	if v, found := BaseHttpGroupCache.GetKvStore(bhgc_key); found {
		//本地查询到缓存数据
		var resbody OwlResponse
		resbody.Status = ResStatus(200)
		resbody.Key = key
		resbody.Data = v.(*cache.KvStore).Value
		resbody.KeyCreateTime = v.(*cache.KvStore).CreateTime
		//设置响应主机
		if v2, found2 := BaseHttpGroupCache.GetKvStore(bhgc_key_responsehost); found2 {
			resbody.ResponseHost = string(v2.(*cache.KvStore).Value)
		} else {
			resbody.ResponseHost = ""
		}
		kvlist.Add(resbody)
		return

	} else {

		//没有在本地缓存中找到数据
		//请求数据，日志记录
		//owllog.OwlLogHttpG.Info(owltools.JoinString("httpclient:get key", " key:", owlhandler.owlrequest.Key, " address:", address))

		//二级缓存的缓存生命周期
		exptime, _ := time.ParseDuration(owltools.JoinString(owlconfig.OwlConfigModel.HttpClientRequestLocalCacheLifeTime, "ms"))

		//判断上次请求状态 1为能正确获得数据  0为未能正确获取数据 0则跳过http请求
		if v3, found3 := BaseHttpGroupCache.GetKvStore(bhgc_key_state); found3 {
			if string(v3.(*cache.KvStore).Value) == "0" {
				return
			}
		}

		//创建一个的HttpClient客户端
		var HttpClient *httpclient.OwlClient
		//初始化HttpClient客户端
		HttpClient = httpclient.NewOwlClient()
		//请求http数据
		s := HttpClient.GetValue(address, key)
		if s != nil {
			//将获取得数据封装
			var resbody OwlResponse
			resbody.Status = ResStatus(s.StatusCode)
			resbody.Key = s.Header.Get("Key")
			resbody.Data = s.Byte()
			//时间处理部分
			tkt := s.Header.Get("Keycreatetime")
			if len(tkt) > 37 {
				//截取字符串的固定长度格式，并把前后两端的空格过滤
				tkt = strings.TrimSpace(tkt[0:37])
			}
			t, terr := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", tkt)
			if terr != nil {
				owllog.OwlLogHttpG.Info("OwlHandler parseContent Keycreatetime time.Parse failed: " + terr.Error())
			}
			resbody.KeyCreateTime = t
			resbody.ResponseHost = s.Header.Get("Responsehost")
			kvlist.Add(resbody)

			//请求优化，缓存记录,配置参数值为0，则不进行缓存(适合并发量小，数据实时性要求高的场景)。
			if owlconfig.OwlConfigModel.HttpClientRequestLocalCacheLifeTime != "0" {
				BaseHttpGroupCache.Set(bhgc_key_state, []byte("1"), exptime)
				BaseHttpGroupCache.Set(bhgc_key_responsehost, []byte(resbody.ResponseHost), exptime)
				BaseHttpGroupCache.Set(bhgc_key, resbody.Data, exptime)
			}
			return
		}
		//请求优化，缓存记录,配置参数值为0，则不进行缓存(适合并发量小，数据实时性要求高的场景)。
		if owlconfig.OwlConfigModel.HttpClientRequestLocalCacheLifeTime != "0" {
			//请求失败
			BaseHttpGroupCache.Set(bhgc_key_state, []byte("0"), exptime)
		}

		return

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
