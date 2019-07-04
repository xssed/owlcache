package counter

import (
	"sync"
)

//创建
func NewCounter() *Counter {
	return &Counter{}
}

//定义
type Counter struct {
	sm sync.Map
}

//+1
func (c *Counter) Add(name string) int {

	if name == "" {
		return 0
	}

	if v, ok := c.sm.Load(name); ok {
		//存在这个Key
		bc := v.(*BaseCounter)
		bc.AddOne()
		return bc.Current()
	}
	//不存在这个Key
	nb := NewBaseCounter()
	nb.AddOne()
	c.sm.Store(name, nb)
	return nb.Current()

}

//重置
func (c *Counter) Reset(name string) {

	if name == "" {
		return
	}

	if v, ok := c.sm.Load(name); ok {
		//存在这个Key
		bc := v.(*BaseCounter)
		bc.Reset()
	}

}
