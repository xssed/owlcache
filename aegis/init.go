package aegis

func AegisInit() {
	//检查密码设置
	CheckConfigPass()
	//检查Redis客户端设置
	CheckRedisConfig()
	//检查Gossip密码设置
	CheckGossipConfig()
	//检查客户端请求优化参数
	CheckClientRequest()
	//检查WebSocket Client请求优化参数
	CheckWebSocketClient()
	//检查urlcache的配置
	CheckUCConfig()
	//捕获程序正常退出操作 ctrl+c
	OnExit()
}
