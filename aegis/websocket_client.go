package aegis

import (
	"os"
	"strconv"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

//检查WebSocket Client请求优化参数
func CheckWebSocketClient() {

	// 写超时
	_, wcww_err := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_WriteWait)
	if wcww_err != nil {
		owllog.OwlLogRun.Println("The value of Websocket_Client_WriteWait is not an integer.Set the <Websocket_Client_WriteWait> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	// 支持接受的消息最大长度，默认7000000字节
	_, wcmms_err := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_MaxMessageSize)
	if wcmms_err != nil {
		owllog.OwlLogRun.Println("The value of Websocket_Client_MaxMessageSize is not an integer.Set the <Websocket_Client_MaxMessageSize> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	// 最小重连时间间隔
	_, wcminrt_err := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_MinRecTime)
	if wcminrt_err != nil {
		owllog.OwlLogRun.Println("The value of Websocket_Client_MinRecTime is not an integer.Set the <Websocket_Client_MinRecTime> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	// 最大重连时间间隔
	_, wcmaxrt_err := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_MaxRecTime)
	if wcmaxrt_err != nil {
		owllog.OwlLogRun.Println("The value of Websocket_Client_MaxRecTime is not an integer.Set the <Websocket_Client_MaxRecTime> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	// 消息发送缓冲池大小，默认1024
	_, wcmbs_err := strconv.Atoi(owlconfig.OwlConfigModel.Websocket_Client_MessageBufferSize)
	if wcmbs_err != nil {
		owllog.OwlLogRun.Println("The value of Websocket_Client_MessageBufferSize is not an integer.Set the <Websocket_Client_MessageBufferSize> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

}
