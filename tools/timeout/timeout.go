package timeout

import (
	"sync"
	"time"
)

//定义Timeout数据结构
type Timeout struct {
	timeouts map[string]time.Time
	mutex    sync.RWMutex
}

//创建一个Timeout
func New() Timeout {
	return Timeout{
		timeouts: make(map[string]time.Time),
	}
}

//为一个Timeout的key设置生命周期
func (to *Timeout) SetTimeout(key string, duration time.Duration) {
	now := time.Now()
	to.mutex.Lock()
	to.timeouts[key] = now.Add(duration) //(time.Time)
	to.mutex.Unlock()
}

//检查key是否超时，删除过期的Key
func (to *Timeout) CheckTimeout(key string) bool {
	now := time.Now()

	to.mutex.RLock()
	timeout, ok := to.timeouts[key]
	to.mutex.RUnlock()
	val := ok && (now.Before(timeout) || now.Equal(timeout))

	if !val {
		to.RemoveTimeout(key)
	}
	return val
}

//删除超时的Key
func (to *Timeout) RemoveTimeout(key string) {
	to.mutex.Lock()
	delete(to.timeouts, key)
	to.mutex.Unlock()
}

//获取所有的key列表
func (to *Timeout) Timeouts() map[string]time.Time {
	to.PruneTimeouts()

	copy := make(map[string]time.Time, len(to.timeouts))
	for key, val := range to.timeouts {
		copy[key] = val
	}
	return copy
}

//遍历key列表删除过期的Key，返回所有没过期的key集合
func (to *Timeout) PruneTimeouts() []string {
	now := time.Now()

	var handles []string
	for key, timeout := range to.timeouts {
		if now.After(timeout) {
			handles = append(handles, key)
			to.RemoveTimeout(key)
		}
	}
	return handles
}
