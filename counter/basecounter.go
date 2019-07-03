package counter

import (
	"sync"
)

//创建
func NewCounter() *BaseCounter {
	return &BaseCounter{}
}

//定义计数器
type BaseCounter struct {
	mut     sync.RWMutex
	currNum int64 //当前数
}

//+1
func (c *BaseCounter) AddOne() int {
	c.mut.Lock()
	c.currNum += 1
	c.mut.Unlock()
	return int(c.currNum)
}

//-1
func (c *BaseCounter) DecOne() int {
	c.mut.Lock()
	c.currNum -= 1
	c.mut.Unlock()
	return int(c.currNum)
}

//获取当前
func (c *BaseCounter) Current() int {
	c.mut.RLock()
	value := c.currNum
	c.mut.RUnlock()
	return int(value)
}

//重置
func (c *BaseCounter) Reset() {
	c.mut.Lock()
	c.currNum = 0
	c.mut.Unlock()
}
