package job

import (
	"fmt"
	"os"
	"strconv"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owlnetwork "github.com/xssed/owlcache/network"
)

func JobInit() {
	fmt.Println("owlcache  job running...")
	//DB数据备份(每分钟)
	DataBackup()
	//Auth数据备份(每分钟)
	DataAuthBackup()
	//定期清理DB中过期的数据(每分钟)
	ClearExpireData()
	//服务器集群信息数据定期备份(每15秒)
	ServerListBackup()
}

// K/V DB数据备份
func DataBackup() {

	//因为使用错误，接连写了N个bug   这段注释掉  作为错误示例吧
	//	task_databackup, err := time.ParseDuration(owlconfig.OwlConfigModel.Task_DataBackup + "m")
	//	if err != nil {
	//		//强制异常，退出
	//		owllog.Println("Config File Task_DataBackup Parse error：" + err.Error()) //日志记录
	//		fmt.Println("Config File Task_DataBackup Parse error：" + err.Error())
	//		os.Exit(0)
	//	}
	task_databackup, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_DataBackup)
	if err != nil {
		owllog.Println("Config File Task_DataBackup Parse error：" + err.Error()) //日志记录
		fmt.Println("Config File Task_DataBackup Parse error：" + err.Error())
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_databackup))
	go func() {
		for _ = range ticker.C {
			//fmt.Printf("ticked at %v", time.Now())
			err := owlnetwork.BaseCacheDB.SaveToFile(owlconfig.OwlConfigModel.DBfile, "owlcache.db")
			if err != nil {
				fmt.Println(err)
				owllog.Error(err)
			}
		}
	}()

}

//Auth数据备份
func DataAuthBackup() {

	task_dataauthbackup, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_DataAuthBackup)
	if err != nil {
		owllog.Println("Config File Task_DataAuthBackup Parse error：" + err.Error()) //日志记录
		fmt.Println("Config File Task_DataAuthBackup Parse error：" + err.Error())
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_dataauthbackup))
	go func() {
		for _ = range ticker.C {
			err := owlnetwork.BaseAuth.SaveToFile(owlconfig.OwlConfigModel.DBfile, "auth.db")
			if err != nil {
				fmt.Println(err)
				owllog.Error(err)
			}
		}
	}()

}

//清理过期的数据
func ClearExpireData() {

	task_clearexpiredata, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_ClearExpireData)
	if err != nil {
		owllog.Println("Config File Task_ClearExpireData Parse error：" + err.Error()) //日志记录
		fmt.Println("Config File Task_ClearExpireData Parse error：" + err.Error())
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_clearexpiredata))
	go func() {
		for _ = range ticker.C {
			owlnetwork.BaseCacheDB.ClearExpireData()
			//owllog.Info("exe ClearExpireData()")
		}
	}()

}

//服务器集群信息数据定期备份
func ServerListBackup() {

	task_serverlistbackup, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_ServerListBackup)
	if err != nil {
		owllog.Println("Config File Task_ServerListBackup Parse error：" + err.Error()) //日志记录
		fmt.Println("Config File Task_ServerListBackup Parse error：" + err.Error())
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_serverlistbackup))
	go func() {
		for _ = range ticker.C {
			err := owlnetwork.ServerGroupList.SaveToFile(owlconfig.OwlConfigModel.DBfile, "servergroup.db")
			if err != nil {
				fmt.Println(err)
				owllog.Error(err)
			}
		}
	}()

}
