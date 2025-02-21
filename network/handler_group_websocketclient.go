package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	//"github.com/xssed/owlcache/cache"
	"github.com/togettoyou/wsc"
	owlconfig "github.com/xssed/owlcache/config"
	"github.com/xssed/owlcache/group"
	owllog "github.com/xssed/owlcache/log"
	"github.com/xssed/owlcache/network/websocketclient"
	owltools "github.com/xssed/owlcache/tools"
	"github.com/xssed/owlcache/tools/timeout"
)

// 存储存活状态的websocketclient列表
var WSClist *websocketclient.WSCList

// 开启Web Socket Client服务，连接到集群
func startWebSocketClient() {

	WSClist = websocketclient.NewWSCList() //WSClist初始化

	list := ServerGroupList.Values()

	for k := range list {
		val, ok := list[k].(group.OwlServerGroupRequest)
		if ok {
			go WebSocketClientConnToServer(val.Address) //客户端连接到ws服务
		}
	}

}

// 客户端连接到ws服务
func WebSocketClientConnToServer(address string) {

	done := make(chan bool)

	ws := websocketclient.NewWebSocketClient(address) //创建资源，配置

	// 连接成功回调
	ws.OnConnected(func() {

		WSClist.AddActiveWSC(ws)                                                                                                      //将活动的WebSocketClient添加到服务器集群信息列表
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("WebSocketClient: Connect to server(OnConnected): ", ws.WebSocket.Url)) //日志记录

	})
	// 连接异常回调，在准备进行连接的过程中发生异常时触发
	ws.OnConnectError(func(err error) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnConnectError: ", err.Error()))
	})
	// 连接断开回调，网络异常，服务端掉线等情况时触发
	ws.OnDisconnected(func(err error) {
		WSClist.RemoveDieWSC(ws) //将失活的WebSocketClient从服务器集群信息列表中删除
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnDisconnected: ", err.Error()))
	})
	// 发送Text消息成功回调
	ws.OnTextMessageSent(func(message string) {
		//owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnTextMessageSent: ", message))
	})
	// 发送Binary消息成功回调
	ws.OnBinaryMessageSent(func(data []byte) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnBinaryMessageSent: ", string(data)))
	})
	// 发送消息异常回调
	ws.OnSentError(func(err error) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnSentError: ", err.Error()))
	})
	// 接受到Ping消息回调
	ws.OnPingReceived(func(appData string) {
		//owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnPingReceived: ", appData))
	})
	// 接受到Pong消息回调
	ws.OnPongReceived(func(appData string) {
		//owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnPongReceived: ", appData))
	})
	// 接受到Text消息回调
	ws.OnTextMessageReceived(func(message string) {
		OwlResponseToWSCGroupCache(message, JsonStrToOwlResponse(message)) //将服务端发送给客户端的数据放进临时存储数据库BaseWSCGroupCache
		//owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnTextMessageReceived: ", message))
	})
	// 接受到Binary消息回调
	ws.OnBinaryMessageReceived(func(data []byte) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnBinaryMessageReceived: ", string(data)))
	})

	// 开始连接
	go ws.Connect()
	// 连接关闭回调，服务端发起关闭信号或客户端主动关闭时触发
	ws.OnClose(func(code int, text string) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnClose: ", strconv.Itoa(code), text))
		WSClist.RemoveDieWSC(ws) //将失活的WebSocketClient从服务器集群信息列表中删除
		done <- true
	})

	for {
		select {
		case <-done:
			return
		}
	}
}

func test() {
	go func() {
		// 连接成功后，测试每5秒发送消息
		t := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-t.C:

				fmt.Println(WSClist.GetList())
				for _, ws := range WSClist.GetList() {
					err := ws.SendTextMessage("get hello data")
					if err == wsc.CloseErr {
						return
					}
				}

			}
		}
	}()
}

// 发起请求获取集合数据
func (owlhandler *OwlHandler) GetWCGroupData(w http.ResponseWriter, r *http.Request) {

	//查询info类型
	resmap := owlhandler.conversionContentInfo(owlhandler.getWSCData())

	//如果valuedata不是字符串info则输出集群第一个
	if string(owlhandler.owlrequest.Value) != "info" {
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

	//如果没有在集群中取到值
	if len(resmap) == 0 {
		owlhandler.Transmit(NOT_FOUND)
		return
	}
	owlhandler.Transmit(SUCCESS)
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")
	data, _ := json.Marshal(&resmap)
	owlhandler.owlresponse.Data = data
	return

}

// 发起WebSocketClient请求获取数据
func (owlhandler *OwlHandler) getWSCData() []OwlResponse {

	list := WSClist.GetList()

	//服务器集群信息存储列表
	groupKVlist := group.NewServergroup()

	var wg sync.WaitGroup

	//判断是否要取指定节点的数据
	if len(owlhandler.owlrequest.Target) > 3 {
		//取指定节点的数据
		var prefix string
		if owlconfig.OwlConfigModel.Open_Https == "1" {
			prefix = "wss://"
		} else {
			prefix = "ws://"
		}
		for _, ws := range list {
			//校验输入的目标服务器是否与配置的服务器一直再进行查询
			//fmt.Println(ws.WebSocket)
			if owltools.JoinString(prefix, owlhandler.owlrequest.Target, "/ws") == ws.WebSocket.Url {
				wg.Add(1)
				wsccontent := NewWSCContent(owlhandler.owlrequest.Key, ws.WebSocket.Url) //定义客户端传输模型
				go owlhandler.parseWSCContent(ws, wsccontent, groupKVlist, &wg)
				wg.Wait()
			}
		}

	} else {
		//遍历服务器集群来取数据
		for _, ws := range list {
			wg.Add(1)
			wsccontent := NewWSCContent(owlhandler.owlrequest.Key, ws.WebSocket.Url) //定义客户端传输模型
			go owlhandler.parseWSCContent(ws, wsccontent, groupKVlist, &wg)
		}
		wg.Wait()
	}

	//fmt.Println(groupKVlist.Values())
	//排序数据
	bubblesortlist := owlhandler.bubbleSortContent(groupKVlist)
	//fmt.Println(bubblesortlist)

	return bubblesortlist

}

// 解析内容
func (owlhandler *OwlHandler) parseWSCContent(ws *websocketclient.OwlWebSocketClient, wsccontent WSCContent, kvlist *group.Servergroup, wg *sync.WaitGroup) {

	//执行完毕自动解除锁定
	defer wg.Done()

	//向WebScoket Server服务器发送查询请求
	//向服务端发送命令格式 get "key_name" info "Handshake_string"    //Handshake_string=uuid+“_”+address
	sql := owltools.JoinString("get ", wsccontent.Key, " info ", wsccontent.Handshake_string)
	err := ws.SendTextMessage(sql)
	if err == wsc.CloseErr {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnTextMessageSent error: ", err.Error()))
		return
	}

	to := timeout.New()
	str := owltools.JoinString(wsccontent.Key, "@", wsccontent.Handshake_string)
	to.SetTimeout(str, time.Second*5)

	//将获取得数据封装
	var resbody OwlResponse

Loop:
	for to.CheckTimeout(str) {
		time.Sleep(time.Millisecond * 2)
		v, b := WSCGroupCacheBurnAfterReading(str) //BaseWSCGroupCache的数据阅后即焚
		if b == true {
			resbody = JsonStrToOwlResponse(string(v)) //将服务之间返回的字符串转换回结构体
			resbody.Key = strings.TrimSpace(strings.Split(resbody.Key, "@")[0])
			//fmt.Println(string(v))
			break Loop
		}
	}

	kvlist.Add(resbody)

	return

}
