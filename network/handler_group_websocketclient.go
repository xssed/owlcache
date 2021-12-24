package network

import (
	"bufio"
	//"fmt"
	"os"
	"strconv"
	"time"

	//"github.com/xssed/owlcache/cache"
	//owlconfig "github.com/xssed/owlcache/config"
	"github.com/togettoyou/wsc"
	"github.com/xssed/owlcache/group"
	owllog "github.com/xssed/owlcache/log"
	"github.com/xssed/owlcache/network/websocketclient"
	owltools "github.com/xssed/owlcache/tools"
)

//开启Web Socket Client服务，连接到集群
func startWebSocketClient() {

	list := ServerGroupList.Values()

	for k := range list {
		val, ok := list[k].(group.OwlServerGroupRequest)
		if ok {
			go WebSocketClientConnToServer(val.Address)
		}
	}

}

//客户端连接到ws服务
func WebSocketClientConnToServer(address string) {

	done := make(chan bool)
	ws := websocketclient.NewWebSocketClient(address)

	// 设置回调处理
	ws.OnConnected(func() {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("WebSocketClient: Connect to server(OnConnected): ", ws.WebSocket.Url))

		// 连接成功后，测试每5秒发送消息
		go func() {
			t := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-t.C:
					err := ws.SendTextMessage("get /77.jpg data")
					if err == wsc.CloseErr {
						return
					}
				}
			}
		}()

	})
	ws.OnConnectError(func(err error) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnConnectError: ", err.Error()))
	})
	ws.OnDisconnected(func(err error) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnDisconnected: ", err.Error()))
	})
	ws.OnTextMessageSent(func(message string) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnTextMessageSent: ", message))
	})
	ws.OnBinaryMessageSent(func(data []byte) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnBinaryMessageSent: ", string(data)))
	})
	ws.OnSentError(func(err error) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnSentError: ", err.Error()))
	})
	ws.OnPingReceived(func(appData string) {
		//owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnPingReceived: ", appData))
	})
	ws.OnPongReceived(func(appData string) {
		//owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnPongReceived: ", appData))
	})
	ws.OnTextMessageReceived(func(message string) {
		//owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnTextMessageReceived: ", message))
		writeResult(message, "2.jpg")
	})
	ws.OnBinaryMessageReceived(func(data []byte) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnBinaryMessageReceived: ", string(data)))
	})

	// 开始连接
	go ws.Connect()

	ws.OnClose(func(code int, text string) {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("OnClose: ", strconv.Itoa(code), text))
		done <- true
	})

	for {
		select {
		case <-done:
			return
		}
	}
}

func writeResult(text string, outfile string) error {

	file, err := os.Create(outfile)
	if err != nil {
		owllog.OwlLogWebsocketClient.Info(owltools.JoinString("writer: ", err.Error()))
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(text)
	writer.Flush()

	return err
}
