package aegis

import (
	"os"
	"strconv"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

//检查密码设置
func CheckConfigPass() {

	//没有设置密码
	if owlconfig.OwlConfigModel.Pass == "" {
		owllog.OwlLogRun.Println("Please set a password first.Set the <Pass> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	//密码长度过低
	if len(owlconfig.OwlConfigModel.Pass) <= 10 {
		owllog.OwlLogRun.Println("Password must be greater than 10.Set the <Pass> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	//不能是纯数字
	if _, err := strconv.Atoi(owlconfig.OwlConfigModel.Pass); err == nil {
		owllog.OwlLogRun.Println("Password cannot be only numbers.Set the <Pass> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

}
