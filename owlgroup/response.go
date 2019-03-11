package owlgroup

import (
	"encoding/json"
)

type OwlServerGroupResponse struct {
	//请求命令
	Cmd CommandType
	//返回状态
	Status ResStatus
	//返回结果
	Results string
	//address
	Address string
}

//将数据转换成json
func (p *OwlServerGroupResponse) ConvertToString() string {
	data, _ := json.Marshal(p)
	s := string(data)
	return s
}
