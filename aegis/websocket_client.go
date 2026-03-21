package aegis

import (
	owlconfig "github.com/xssed/owlcache/config"
)

// 检查WebSocket Client请求优化参数
func CheckWebSocketClient() {
	// 写超时
	checkClientConfig("Websocket_Client_WriteWait", owlconfig.OwlConfigModel.Websocket_Client_WriteWait, "", false, 0)
	// 支持接受的消息最大长度，默认7000000字节
	checkClientConfig("Websocket_Client_MaxMessageSize", owlconfig.OwlConfigModel.Websocket_Client_MaxMessageSize, "", false, 0)
	// 最小重连时间间隔
	checkClientConfig("Websocket_Client_MinRecTime", owlconfig.OwlConfigModel.Websocket_Client_MinRecTime, "", false, 0)
	// 最大重连时间间隔
	checkClientConfig("Websocket_Client_MaxRecTime", owlconfig.OwlConfigModel.Websocket_Client_MaxRecTime, "", false, 0)
	// 消息发送缓冲池大小，默认1024
	checkClientConfig("Websocket_Client_MessageBufferSize", owlconfig.OwlConfigModel.Websocket_Client_MessageBufferSize, "", false, 0)
}
