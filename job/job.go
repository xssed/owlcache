package job

import (
	"fmt"
	"os"
	"strconv"
	"time"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owlnetwork "github.com/xssed/owlcache/network"
	owlsystem "github.com/xssed/owlcache/system"
	owltools "github.com/xssed/owlcache/tools"
)

func JobInit() {
	fmt.Println("owlcache  job running...")
	//DB数据备份
	DataBackup()
	//Auth数据备份
	DataAuthBackup()
	//定期清理DB中过期的数据
	ClearExpireData()
	//服务器集群信息数据定期备份
	ServerListBackup()
	//Gossip服务器集群信息数据定期备份
	ServerGossipListBackup()
	//定时自动输出Owl的内存信息
	MemoryInfoToLog()
	//凌晨时分创建新一天的日志
	TimerToCreateLogerInBeforeDawn()
}

// K/V DB数据备份
func DataBackup() {

	task_databackup, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_DataBackup)
	if err != nil {
		owllog.OwlLogTask.Info(owltools.JoinString("Config File Task_DataBackup Parse error:", err.Error())) //日志记录
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_databackup))
	go func() {
		for _ = range ticker.C {
			//fmt.Printf("ticked at %v", time.Now())//调试
			err := owlnetwork.BaseCacheDB.SaveToFile(owlconfig.OwlConfigModel.DBfile, "owlcache.db")
			if err != nil {
				owllog.OwlLogTask.Info(owltools.JoinString("Task: DataBackup() error ", err.Error()))
			}
		}
	}()

}

//Auth数据备份
func DataAuthBackup() {

	task_dataauthbackup, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_DataAuthBackup)
	if err != nil {
		owllog.OwlLogTask.Info("Config File Task_DataAuthBackup Parse error:" + err.Error()) //日志记录
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_dataauthbackup))
	go func() {
		for _ = range ticker.C {
			err := owlnetwork.BaseAuth.SaveToFile(owlconfig.OwlConfigModel.DBfile, "auth.db")
			if err != nil {
				owllog.OwlLogTask.Info(owltools.JoinString("Task: DataAuthBackup() error ", err.Error()))
			}
		}
	}()

}

//清理过期的数据
func ClearExpireData() {

	task_clearexpiredata, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_ClearExpireData)
	if err != nil {
		owllog.OwlLogTask.Info("Config File Task_ClearExpireData Parse error:" + err.Error()) //日志记录
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_clearexpiredata))
	go func() {
		for _ = range ticker.C {
			owlnetwork.BaseCacheDB.ClearExpireData()
		}
	}()

}

//服务器集群信息数据定期备份
func ServerListBackup() {

	task_serverlistbackup, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_ServerListBackup)
	if err != nil {
		owllog.OwlLogTask.Info("Config File Task_ServerListBackup Parse error:" + err.Error()) //日志记录
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_serverlistbackup))
	go func() {
		for _ = range ticker.C {
			err := owlnetwork.ServerGroupList.SaveToFile(owlconfig.OwlConfigModel.DBfile, "server_group_config.json")
			if err != nil {
				owllog.OwlLogTask.Info(owltools.JoinString("Task: ServerListBackup() error ", err.Error()))
			}
		}
	}()

}

//Gossip服务器集群信息数据定期备份
func ServerGossipListBackup() {

	task_servergossiplistbackup, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_ServerGossipListBackup)
	if err != nil {
		owllog.OwlLogTask.Info("Config File Task_ServerGossipListBackup Parse error:" + err.Error()) //日志记录
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_servergossiplistbackup))
	go func() {
		for _ = range ticker.C {
			err := owlnetwork.ServerGroupGossipList.SaveToFile(owlconfig.OwlConfigModel.DBfile, "server_group_gossip_config.json")
			if err != nil {
				owllog.OwlLogTask.Info(owltools.JoinString("Task: ServerGossipListBackup() error ", err.Error()))
			}
		}
	}()

}

//统计内存使用情况
func MemoryInfoToLog() {

	task_memoryinfotolog, err := strconv.Atoi(owlconfig.OwlConfigModel.Task_MemoryInfoToLog)
	if err != nil {
		owllog.OwlLogTask.Info("Config File Task_MemoryInfoToLog Parse error:" + err.Error()) //日志记录
		os.Exit(0)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(task_memoryinfotolog))
	go func() {
		for _ = range ticker.C {
			owllog.OwlLogSystemResource.Info(owlsystem.MemStats())
		}
	}()

}

//凌晨时分创建新一天的日志
func TimerToCreateLogerInBeforeDawn() {
	go func() {
		for {
			//获取当前时间
			now := time.Now()
			//获取24个小时之后的时间
			next := now.Add(time.Hour * 24)
			//获取下一个凌晨零点的日期
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			//计算当前时间到凌晨的时间间隔，设置一个定时器
			sub := next.Sub(now)
			t := time.NewTimer(sub)
			//输出响应
			fmt.Println(owltools.JoinString("owlcache  will create a new log directory at ", next.String(), " after ", sub.String()))
			<-t.C
			//重新初始化创建新的一天的日志
			owllog.LogInit()
			//日志记录
			owllog.OwlLogTask.Info("Create a new log directory.")
		}
	}()
}
