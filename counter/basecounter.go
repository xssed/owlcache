package counter

import (
	"sync"
	"time"
)

//创建
func NewBaseCounter(max int64, lt time.Duration) *BaseCounter {
	return &BaseCounter{maxNum: max, lifeTime: lt, status: 0}
}

//定义计数器
type BaseCounter struct {
	mut        sync.RWMutex
	currNum    int64 //当前数
	maxNum     int64 //最大数
	status     int
	lifeTime   time.Duration
	createTime time.Time
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

//获取当前状态
func (c *BaseCounter) CurrentStatus() int {
	c.mut.RLock()
	value := c.status
	c.mut.RUnlock()
	return value
}

//获取当前统计了多少数
func (c *BaseCounter) CurrentNum() int {
	c.mut.RLock()
	value := c.currNum
	c.mut.RUnlock()
	return int(value)
}

//获取允许最大值
func (c *BaseCounter) CurrentMaxNum() int {
	c.mut.RLock()
	value := c.maxNum
	c.mut.RUnlock()
	return int(value)
}

//重置
func (c *BaseCounter) Reset() {
	c.mut.Lock()
	c.currNum = 0
	c.mut.Unlock()
}

//重置时间
func (c *BaseCounter) ReTime() {
	c.mut.Lock()
	c.createTime = time.Now()
	c.mut.Unlock()
}

//重置
func (c *BaseCounter) ReStatusOn() {
	c.mut.Lock()
	c.status = 0
	c.mut.Unlock()
}

//重置
func (c *BaseCounter) ReStatusOff() {
	c.mut.Lock()
	c.status = -1
	c.mut.Unlock()
}

//检验是否过期
func (c *BaseCounter) IsExpired() bool {
	//判断永不过期
	if c.lifeTime == 0 {
		return false
	}
	return time.Since(c.createTime) > c.lifeTime
}

//超过最大值
func (c *BaseCounter) IsBad() bool {
	//判断不设最大值
	if c.CurrentMaxNum() == 0 {
		return true
	}
	return c.CurrentMaxNum() < c.CurrentNum()
}

//判断可用性
func (c *BaseCounter) IsUse() bool {
	if !c.IsBad() { //&& !c.IsExpired()
		return true
	} else {
		return false
	}
}
