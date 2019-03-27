package job

import (
	"fmt"
	"time"

	owllog "github.com/xssed/owlcache/log"
	owlnetwork "github.com/xssed/owlcache/network"
)

func JobInit() {
	fmt.Println("owlcache  job running...")
	//数据备份
	DataBackup()
	//清理过期的数据
	ClearExpireData()
	//服务器集群信息数据定期备份
	ServerListBackup()
}

// K/V数据备份
func DataBackup() {

	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			//fmt.Printf("ticked at %v", time.Now())
			err := owlnetwork.BaseCacheDB.SaveToFile("./owlcache.db")
			if err != nil {
				fmt.Println(err)
				owllog.Error(err)
			}
		}
	}()

}

//清理过期的数据
func ClearExpireData() {

	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			//fmt.Printf("ticked at %v", time.Now())
			owlnetwork.BaseCacheDB.ClearExpireData()
			//owllog.Info("exe ClearExpireData()")
		}
	}()

}

//服务器集群信息数据定期备份
func ServerListBackup() {

	ticker := time.NewTicker(time.Minute * 1)
	go func() {
		for _ = range ticker.C {
			//fmt.Printf("ticked at %v", time.Now())
			err := owlnetwork.ServerGroupList.SaveToFile("./servergroup.db")
			if err != nil {
				fmt.Println(err)
				owllog.Error(err)
			}
		}
	}()

}
