package gossip

import (
	"github.com/hashicorp/memberlist"
)

//广播
type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

//无效信号
func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

//消息
func (b *broadcast) Message() []byte {
	return b.msg
}

//完成操作
func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}
