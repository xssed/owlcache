package network

import (
	"encoding/json"
	//"fmt"
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
	Data interface{}
}

//将数据转换成json
func (p *OwlResponse) ConvertToString() string {
	data, _ := json.Marshal(p)
	s := string(data)
	return s
}
