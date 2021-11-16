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

//执行计数器-1是不允许执行，>0是计数器开始工作，可执行(错误请求时+1)
func (c *Counter) Exe(name string, max int64, lt time.Duration) int {
	//name没有写
	if name == "" {
		return -1
	}
	//判读Key是否存在
	if !c.Exists(name) {
		return c.Add(name, max, lt)
	} else {
		v, _ := c.sm.Load(name)
		bc := v.(*BaseCounter)

		//是否超过最大值
		if !bc.IsBad() {
			bc.AddOne()            //+1
			return bc.CurrentNum() //返回当前数
		} else {
			//超过最大值
			//判断当前状态
			if bc.CurrentStatus() == 0 {
				bc.ReTime()
				bc.ReStatusOff()
			} else {
				if bc.IsExpired() {
					bc.ReStatusOn()
					bc.Reset()
				}
			}
			return bc.CurrentStatus()
		}
	}
	return -1
}

//+1
func (c *Counter) Add(name string, max int64, lt time.Duration) int {
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
		bc.ReTime()
	}
}

//-1
func (c *Counter) Dec(name string) {
	if name == "" {
		return
	}
	if v, ok := c.sm.Load(name); ok {
		//存在这个Key
		bc := v.(*BaseCounter)
		bc.DecOne()
	}
}

//判断一个key是否存在
func (c *Counter) Exists(key string) bool {
	_, ok := c.sm.Load(key)
	return ok
}

//删除一个key
func (c *Counter) Delete(key string) bool {
	c.sm.Delete(key)
	return true
}
