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
		//校验Gossip端口
		_, atio_err := strconv.Atoi(owlconfig.OwlConfigModel.Gossipport)
		if atio_err != nil {
			owllog.OwlLogRun.Println("The configuration file <Gossipport> option is not a valid number!")
			os.Exit(0)
		}
		//检测能否正确获取主机名
		_, get_hostname_err := os.Hostname()
		if get_hostname_err != nil {
			owllog.OwlLogRun.Println("When starting the gossip service, getting the Hostname failed! Please check the execution permission of owlcache!")
			os.Exit(0)
		}
	}

}
