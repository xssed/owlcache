package counter

import (
	//"fmt"
	"sync"
	"time"
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
func (c *Counter) Add(name string, max int64, lt time.Duration) int {

	if name == "" {
		return -1
	}

	if v, ok := c.sm.Load(name); ok {
		//存在这个Key
		bc := v.(*BaseCounter)
		if bc.IsUse() {
			bc.AddOne()
			return bc.CurrentNum()
		} else {
			if bc.CurrentStatus() == 0 {
				bc.ReTime()
				bc.ReStatusOff()
			} else {
				if bc.IsExpired() {
					bc.ReStatusOn()
					bc.Reset()
				}
			}
			return -1
		}

	}

	//不存在这个Key
	nb := NewBaseCounter(max, lt)
	nb.AddOne()
	c.sm.Store(name, nb)
	return nb.CurrentNum()

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
