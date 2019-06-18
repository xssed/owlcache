package aegis

func AegisInit() {
	//检查密码设置
	CheckConfigPass()
	//检查Redis客户端设置
	CheckRedisConfig()
	//捕获程序正常退出操作 ctrl+c
	OnExit()
}
