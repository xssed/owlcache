package gossip

import (
	"encoding/json"
	"fmt"
)

//代表
type delegate struct{}

//节点元数据
func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

//通知消息
func (d *delegate) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}
	//将通讯数据的头字节取出(自定义数据协议)
	str := string(b)
	pass := H.Password
	protocol_head_len := len(pass)

	switch str[:protocol_head_len] {
	//验证密码
	case pass:
		//定义数据包列表
		var executes []*Execute
		//绑定数据到数据包列表
		if err := json.Unmarshal([]byte(str[protocol_head_len:]), &executes); err != nil {
			return
		}
		//遍历取单个数据包操作
		for _, u := range executes {
			//把map[string]string的Data数据取出
			jsons, errs := json.Marshal(&u) //转换成JSON返回的是byte[]
			if errs != nil {
				fmt.Println(errs.Error())
			}
			Q.Push(string(jsons)) //将数据发送到队列
		}
	}
}

//获取广播
func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return H.broadcasts.GetBroadcasts(overhead, limit)
}

//本地状态，将数据转化成JSON数据返回
func (d *delegate) LocalState(join bool) []byte {
	return []byte{}
}

//合并远程状态
func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}
}
