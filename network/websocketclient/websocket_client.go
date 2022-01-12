package websocketclient

import (
	"strconv"
	"time"

	"github.com/togettoyou/wsc"
	owlconfig "github.com/xssed/owlcache/config"
	owltools "github.com/xssed/owlcache/tools"
)

//定义websocket_client客户端结构
type OwlWebSocketClient struct {
	*wsc.Wsc
}

//创建websocket_client客户端实体
func NewWebSocketClient(address string) *OwlWebSocketClient {

	var prefix string
	if owlconfig.OwlConfigModel.Open_Https == "1" {
		prefix = "wss://"
	} else {
		prefix = "ws://"
	}

	//自定义WS客户端配置，写超时
	wcww, _ := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_WriteWait)
	//自定义WS客户端配置，支持接受的消息最大长度，默认7000000字节
	wcmms, _ := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_MaxMessageSize)
	//自定义WS客户端配置，最小重连时间间隔
	wcminrt, _ := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_MinRecTime)
	//自定义WS客户端配置，最大重连时间间隔
	wcmaxrt, _ := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_MaxRecTime)
	//自定义WS客户端配置，消息发送缓冲池大小，默认1024
	wcmbs, _ := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_MessageBufferSize)

	ws := wsc.New(owltools.JoinString(prefix, address, "/ws"))
	// 可自定义配置，不使用默认配置
	ws.SetConfig(&wsc.Config{
		// 写超时
		WriteWait: time.Duration(wcww) * time.Second,
		// 支持接受的消息最大长度，默认7000000字节
		MaxMessageSize: int64(wcmms),
		// 最小重连时间间隔
		MinRecTime: time.Duration(wcminrt) * time.Second,
		// 最大重连时间间隔
		MaxRecTime: time.Duration(wcmaxrt) * time.Second,
		// 每次重连失败继续重连的时间间隔递增的乘数因子，递增到最大重连时间间隔为止
		RecFactor: 1.5,
		// 消息发送缓冲池大小，默认1024
		MessageBufferSize: wcmbs,
	})

	return &OwlWebSocketClient{ws}

}
