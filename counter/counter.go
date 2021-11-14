package counter

import (
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

	//去查找计数器之前是否有这个name存进来
	if v, ok := c.sm.Load(name); ok {
		//存在这个Key
		bc := v.(*BaseCounter)
		//是否超过最大值
		if bc.IsBad() {
			//没有超过
			bc.AddOne()            //+1
			return bc.CurrentNum() //返回当前数
		} else {
			//超过最大值，重置
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
	// nb := NewBaseCounter(max, lt)
	// nb.AddOne()
	// c.sm.Store(name, nb)
	// return nb.CurrentNum()

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
		bc.ReTime()
	}

}
