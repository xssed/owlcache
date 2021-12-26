package websocketclient

import (
	"sync"
)

//创建一个客户端文本传输模型
type WSCList struct {
	lock sync.RWMutex          //活动的服务器集群信息列表的读写锁
	list []*OwlWebSocketClient //创建活动的服务器集群信息列表
}

//WSCList初始化
func NewWSCList() *WSCList {
	return &WSCList{}
}

//返回活动的WebSocketClient服务器集群信息列表
func (wsclist *WSCList) GetList() []*OwlWebSocketClient {

	wsclist.lock.RLock()
	defer wsclist.lock.RUnlock()
	return wsclist.list

}

//将活动的WebSocketClient添加到服务器集群信息列表
func (wsclist *WSCList) AddActiveWSC(temp_wsc *OwlWebSocketClient) {

	wsclist.lock.Lock()
	defer wsclist.lock.Unlock()

	wsclist.list = append(wsclist.list, temp_wsc)

}

//将失活的WebSocketClient从服务器集群信息列表中删除
func (wsclist *WSCList) RemoveDieWSC(temp_wsc *OwlWebSocketClient) {

	wsclist.lock.Lock()
	defer wsclist.lock.Unlock()

	for index, t_ws := range wsclist.list {
		if temp_wsc == t_ws {
			wsclist.list = append(wsclist.list[:index], wsclist.list[index+1:]...)
		}
	}

}
