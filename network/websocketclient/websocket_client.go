package websocketclient

import (
	//"bufio"
	//"log"
	//"os"
	//"fmt"
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
	ws := wsc.New(owltools.JoinString(prefix, address, "/ws"))
	// 可自定义配置，不使用默认配置
	ws.SetConfig(&wsc.Config{
		// 写超时
		WriteWait: 10 * time.Second,
		// 支持接受的消息最大长度，默认512字节
		MaxMessageSize: 200000,
		// 最小重连时间间隔
		MinRecTime: 2 * time.Second,
		// 最大重连时间间隔
		MaxRecTime: 60 * time.Second,
		// 每次重连失败继续重连的时间间隔递增的乘数因子，递增到最大重连时间间隔为止
		RecFactor: 1.5,
		// 消息发送缓冲池大小，默认256
		MessageBufferSize: 1024,
	})

	return &OwlWebSocketClient{ws}

}
