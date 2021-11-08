package network

import (
	"time"
)

type OwlResponse struct {
	//请求命令
	Cmd CommandType
	//返回状态
	Status ResStatus
	//返回结果
	Results string
	//key
	Key string
	//返回内容
	Data []byte
	//程序响应IP
	ResponseHost string
	//内容的创建时间
	KeyCreateTime time.Time
}
