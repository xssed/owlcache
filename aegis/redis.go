package aegis

import (
	"os"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

//检查Redis客户端设置
func CheckRedisConfig() {

	//如果开启从Redis中获取数据
	if owlconfig.OwlConfigModel.Get_data_from_redis == "1" {
		//地址必填
		if len(owlconfig.OwlConfigModel.Redis_Addr) <= 6 {
			owllog.OwlLogRun.Println("Redis_Addr length must be greater than 6.Set the <Redis_Addr> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
			os.Exit(0)
		}
		//没有设置密码
		if owlconfig.OwlConfigModel.Redis_Password == "" {
			owllog.OwlLogRun.Println("Please set a redis password.Set the <Redis_Password> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
			os.Exit(0)
		}

	}

}
