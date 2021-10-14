package aegis

import (
	"os"
	"strconv"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

//检查Gossip密码设置
func CheckGossipConfig() {

	//检查密码
	if owlconfig.OwlConfigModel.GroupDataSync == "1" {
		//没有设置密码
		if owlconfig.OwlConfigModel.GossipDataSyncAuthKey == "" {
			owllog.OwlLogRun.Println("Please set a password first.Set the <GossipDataSyncAuthKey> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
			os.Exit(0)
		}
		//密码长度过低
		if len(owlconfig.OwlConfigModel.GossipDataSyncAuthKey) <= 10 {
			owllog.OwlLogRun.Println("Password must be greater than 10.Set the <GossipDataSyncAuthKey> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
			os.Exit(0)
		}
		//不能是纯数字
		if _, err := strconv.Atoi(owlconfig.OwlConfigModel.GossipDataSyncAuthKey); err == nil {
			owllog.OwlLogRun.Println("Password cannot be only numbers.Set the <GossipDataSyncAuthKey> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
			os.Exit(0)
		}
	}

}
