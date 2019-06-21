package redisclient

import (
	"os"
	"strconv"

	"github.com/go-redis/redis"
	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

var rc *redis.Client

//客户端初始化
func Start() {
	//检查数据库设置是否正确
	select_db, _ := strconv.Atoi(owlconfig.OwlConfigModel.Redis_DB)
	//检查是否能够连接到Redis服务
	rc = redis.NewClient(&redis.Options{
		Addr:     owlconfig.OwlConfigModel.Redis_Addr,
		Password: owlconfig.OwlConfigModel.Redis_Password,
		DB:       select_db,
	})
	pong, err := rc.Ping().Result()
	if err != nil && pong != "PONG" {
		owllog.OwlLogRun.Println("owlcache failed to connect to redis:", err)
		owllog.OwlLogRun.Println("Please alter the redis password & address.Set the <Redis_Addr> <Redis_Password> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
}

//Get方法
func Get(key string) (string, error) {

	val, err := rc.Get(key).Result()
	if err != nil {
		owllog.OwlLogRun.Println("Redis Client Get() error:", err)
		return "", err
	}
	return val, nil

}
