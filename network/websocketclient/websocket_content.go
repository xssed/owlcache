package websocketclient

import (
	"sync"
)

//创建一个客户端文本传输模型
type WSCContent struct {
	lock     sync.RWMutex
	KeyItems sync.Map
}
