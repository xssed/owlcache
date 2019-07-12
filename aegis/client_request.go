package aegis

import (
	"os"
	"strconv"

	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
)

//检查客户端请求优化参数
func CheckClientRequest() {

	hrts, hrts_err := strconv.Atoi(owlconfig.OwlConfigModel.HttpClient_Request_Timeout_Sleeptime)
	if hrts_err != nil {
		owllog.OwlLogRun.Println("The value of HttpClient_Request_Timeout_Sleeptime is not an integer.Set the <HttpClient_Request_Timeout_Sleeptime> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	if hrts < 2 {
		owllog.OwlLogRun.Println("The value of HttpClient_Request_Timeout_Sleeptime cannot be less than 2.Set the <HttpClient_Request_Timeout_Sleeptime> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

	hrmen, hrmen_err := strconv.Atoi(owlconfig.OwlConfigModel.HttpClient_Request_Max_Error_Number)
	if hrmen_err != nil {
		owllog.OwlLogRun.Println("The value of HttpClient_Request_Max_Error_Number is not an integer.Set the <HttpClient_Request_Max_Error_Number> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	if hrmen < 2 {
		owllog.OwlLogRun.Println("The value of HttpClient_Request_Max_Error_Number cannot be less than 2.Set the <HttpClient_Request_Max_Error_Number> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

	mrts, mrts_err := strconv.Atoi(owlconfig.OwlConfigModel.MemcacheClient_Request_Timeout_Sleeptime)
	if mrts_err != nil {
		owllog.OwlLogRun.Println("The value of MemcacheClient_Request_Timeout_Sleeptime is not an integer.Set the <MemcacheClient_Request_Timeout_Sleeptime> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	if mrts < 2 {
		owllog.OwlLogRun.Println("The value of MemcacheClient_Request_Timeout_Sleeptime cannot be less than 2.Set the <MemcacheClient_Request_Timeout_Sleeptime> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

	mrmen, mrmen_err := strconv.Atoi(owlconfig.OwlConfigModel.MemcacheClient_Request_Max_Error_Number)
	if mrmen_err != nil {
		owllog.OwlLogRun.Println("The value of MemcacheClient_Request_Max_Error_Number is not an integer.Set the <MemcacheClient_Request_Max_Error_Number> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	if mrmen < 2 {
		owllog.OwlLogRun.Println("The value of MemcacheClient_Request_Max_Error_Number cannot be less than 2.Set the <MemcacheClient_Request_Max_Error_Number> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

	rrts, rrts_err := strconv.Atoi(owlconfig.OwlConfigModel.RedisClient_Request_Timeout_Sleeptime)
	if rrts_err != nil {
		owllog.OwlLogRun.Println("The value of RedisClient_Request_Timeout_Sleeptime is not an integer.Set the <RedisClient_Request_Timeout_Sleeptime> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	if rrts < 2 {
		owllog.OwlLogRun.Println("The value of RedisClient_Request_Timeout_Sleeptime cannot be less than 2.Set the <RedisClient_Request_Timeout_Sleeptime> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

	rrmen, rrmen_err := strconv.Atoi(owlconfig.OwlConfigModel.RedisClient_Request_Max_Error_Number)
	if rrmen_err != nil {
		owllog.OwlLogRun.Println("The value of RedisClient_Request_Max_Error_Number is not an integer.Set the <RedisClient_Request_Max_Error_Number> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	if rrmen < 2 {
		owllog.OwlLogRun.Println("The value of RedisClient_Request_Max_Error_Number cannot be less than 2.Set the <RedisClient_Request_Max_Error_Number> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

}
