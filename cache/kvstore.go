package cache

import (
	"time"
)

//创建key/value模型
type KvStore struct {
	Key        interface{}
	Value      interface{}
	LifeTime   time.Duration
	CreateTime time.Time
}

//创建一条key/value内容
func newKvStore(key interface{}, value interface{}, lifeTime time.Duration) *KvStore {
	return &KvStore{
		key,
		value,
		lifeTime,
		time.Now()}
}

//检验是否过期
func (kvstore *KvStore) IsExpired() bool {
	//判断永不过期
	if kvstore.LifeTime == 0 {
		return false
	}
	return time.Since(kvstore.CreateTime) > kvstore.LifeTime
}
