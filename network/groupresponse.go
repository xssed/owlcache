package network

import (
	"encoding/json"
)

type OwlServerGroupResponse struct {
	//请求命令
	Cmd GroupCommandType
	//返回状态
	Status ResStatus
	//返回结果
	Results string
	//address
	Address string
	//返回内容
	Data interface{}
}

//将数据转换成json
func (p *OwlServerGroupResponse) ConvertToString() string {
	data, _ := json.Marshal(p)
	s := string(data)
	return s
}
