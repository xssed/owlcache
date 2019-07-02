package counter

import (
	"sync"
	"sync/atomic"
)

//创建
func NewCounter() *BaseCounter {
	return &BaseCounter{}
}

//定义计数器
type BaseCounter struct {
	mut     sync.Mutex
	currNum int64 //当前数
}

//+1
func (c *BaseCounter) AddOne() int {
	return int(atomic.AddInt64(&c.currNum, 1))
}

//-1
func (c *BaseCounter) DecOne() int {
	return int(atomic.AddInt64(&c.currNum, -1))

}

//获取当前
func (c *BaseCounter) Current() int {
	return int(atomic.LoadInt64(&c.currNum))
}

//重置
func (c *BaseCounter) Reset() {
	c.mut.Lock()
	c.currNum = 0
	c.mut.Unlock()
}
