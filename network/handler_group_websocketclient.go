package network

import (
	"strconv"
	//"sync"
	"fmt"
	"time"

	//"github.com/xssed/owlcache/cache"
	//owlconfig "github.com/xssed/owlcache/config"
	"github.com/togettoyou/wsc"
	"github.com/xssed/owlcache/group"
	owllog "github.com/xssed/owlcache/log"
	"github.com/xssed/owlcache/network/websocketclient"
	owltools "github.com/xssed/owlcache/tools"
)

var WSClist *websocketclient.WSCList

//开启Web Socket Client服务，连接到集群
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

//客户端连接到ws服务
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
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString(ws.WebSocket.Url, " OnTextMessageReceived: ", message))
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
