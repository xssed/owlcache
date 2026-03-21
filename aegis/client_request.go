package aegis

import (
	"fmt"
	"os"
	"strconv"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

// 配置项结构体
type clientConfig struct {
	name         string
	value        string
	minValue     int
	configKey    string
	needMinCheck bool
}

// 检查整数类型配置
func checkIntConfig(config clientConfig) int {
	value, err := strconv.Atoi(config.value)
	if err != nil {
		owllog.OwlLogRun.Println(fmt.Sprintf("The value of %s is not an integer. Set the <%s> option in the configuration file %s.", config.name, config.name, owlconfig.OwlConfigModel.Configfile))
		os.Exit(0)
	}

	if config.needMinCheck && value < config.minValue {
		owllog.OwlLogRun.Println(fmt.Sprintf("The value of %s cannot be less than %d. Set the <%s> option in the configuration file %s.", config.name, config.minValue, config.name, owlconfig.OwlConfigModel.Configfile))
		os.Exit(0)
	}

	return value
}

// 检查客户端配置
func checkClientConfig(name, value, configKey string, needMinCheck bool, minValue int) int {
	config := clientConfig{
		name:         name,
		value:        value,
		minValue:     minValue,
		configKey:    configKey,
		needMinCheck: needMinCheck,
	}
	return checkIntConfig(config)
}

// 检查HTTP客户端请求配置
func checkClientRequestConfig() {
	checkClientConfig("HttpClientRequestLocalCacheLifeTime", owlconfig.OwlConfigModel.HttpClientRequestLocalCacheLifeTime, "", false, 0)
	checkClientConfig("HttpClientRequestTimeout", owlconfig.OwlConfigModel.HttpClientRequestTimeout, "", false, 0)
}

// 检查Memcache客户端配置
func checkMemcacheClientConfig() {
	timeoutSleepTime := checkClientConfig("MemcacheClient_Request_Timeout_Sleeptime", owlconfig.OwlConfigModel.MemcacheClient_Request_Timeout_Sleeptime, "", true, 2)
	maxErrorNumber := checkClientConfig("MemcacheClient_Request_Max_Error_Number", owlconfig.OwlConfigModel.MemcacheClient_Request_Max_Error_Number, "", true, 2)

	_ = timeoutSleepTime
	_ = maxErrorNumber
}

// 检查Redis客户端配置
func checkRedisClientConfig() {
	timeoutSleepTime := checkClientConfig("RedisClient_Request_Timeout_Sleeptime", owlconfig.OwlConfigModel.RedisClient_Request_Timeout_Sleeptime, "", true, 2)
	maxErrorNumber := checkClientConfig("RedisClient_Request_Max_Error_Number", owlconfig.OwlConfigModel.RedisClient_Request_Max_Error_Number, "", true, 2)

	_ = timeoutSleepTime
	_ = maxErrorNumber
}

// 检查客户端请求优化参数
func CheckClientRequest() {
	checkClientRequestConfig()
	checkMemcacheClientConfig()
	checkRedisClientConfig()
}
