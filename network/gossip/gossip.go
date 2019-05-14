package gossip

import (
	"github.com/xssed/owlcache/queue"
)

var (
	H *Handler
	Q *queue.Queue
)

//初始化
func init() {
	H = NewHandler()
	Q = queue.New()
}
